package dms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Kafka POST /v2/kafka/{project_id}/instances/{instance_id}/kafka-user-client-quota
// @API Kafka PUT /v2/kafka/{project_id}/instances/{instance_id}/kafka-user-client-quota
// @API Kafka GET /v2/kafka/{project_id}/instances/{instance_id}/kafka-user-client-quota
// @API Kafka DELETE /v2/kafka/{project_id}/instances/{instance_id}/kafka-user-client-quota
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/tasks
// @API Kafka GET /v2/{project_id}/instances/{instance_id}
func ResourceDmsKafkaUserClientQuota() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaUserClientQuotaCreate,
		ReadContext:   resourceDmsKafkaUserClientQuotaRead,
		UpdateContext: resourceDmsKafkaUserClientQuotaUpdate,
		DeleteContext: resourceDmsKafkaUserClientQuotaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(50 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the Kafka instance.`,
			},
			"user": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  `Specifies the user name to apply the quota.`,
				AtLeastOneOf: []string{"user", "user_default", "client", "client_default"},
			},
			"user_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the user default configuration of the quota.`,
			},
			"client": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the client to which the quota applies.`,
			},
			"client_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the client default configuration of the quota.`,
			},
			"producer_byte_rate": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  `Specifies an upper limit on the prodction rate. The unit is B/s.`,
				AtLeastOneOf: []string{"producer_byte_rate", "consumer_byte_rate"},
			},
			"consumer_byte_rate": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies an upper limit on the consumption rate. The unit is B/s.`,
			},
		},
	}
}

func resourceDmsKafkaUserClientQuotaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	// createKafkaUserClientQuota: create DMS kafka user client quota
	var (
		createKafkaUserClientQuotaHttpUrl = "v2/kafka/{project_id}/instances/{instance_id}/kafka-user-client-quota"
		createKafkaUserClientQuotaProduct = "dmsv2"
	)
	createKafkaUserClientQuotaClient, err := cfg.NewServiceClient(createKafkaUserClientQuotaProduct, region)

	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createKafkaUserClientQuotaPath := createKafkaUserClientQuotaClient.Endpoint + createKafkaUserClientQuotaHttpUrl
	createKafkaUserClientQuotaPath = strings.ReplaceAll(createKafkaUserClientQuotaPath, "{project_id}",
		createKafkaUserClientQuotaClient.ProjectID)
	createKafkaUserClientQuotaPath = strings.ReplaceAll(createKafkaUserClientQuotaPath, "{instance_id}", instanceID)

	createKafkaUserClientQuotaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createKafkaUserClientQuotaOpt.JSONBody = utils.RemoveNil(buildKafkaUserClientQuotaBodyParams(d))

	// The quota is allowd to create only when the instance status is RUNNING.
	retryFunc := func() (interface{}, bool, error) {
		createKafkaUserClientQuotaResp, createErr := createKafkaUserClientQuotaClient.Request("POST",
			createKafkaUserClientQuotaPath, &createKafkaUserClientQuotaOpt)
		retry, err := handleOperationConflictError(createErr)
		return createKafkaUserClientQuotaResp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     KafkaInstanceStateRefreshFunc(createKafkaUserClientQuotaClient, instanceID),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return diag.Errorf("error creating DMS kafka user client quota: %v", err)
	}

	_, err = utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	user := d.Get("user").(string)
	userDefault := d.Get("user_default").(bool)
	client := d.Get("client").(string)
	clientDefault := d.Get("client_default").(bool)
	d.SetId(instanceID + "/" + user + "/" + strconv.FormatBool(userDefault) + "/" + client + "/" + strconv.FormatBool(clientDefault))

	// The quota creation triggers a related task, if the task status is SUCCESS, the quota has been created.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATED"},
		Target:       []string{"SUCCESS"},
		Refresh:      userClientQuotaTaskRefreshFunc(createKafkaUserClientQuotaClient, instanceID, d, "kafkaClientQuotaCreate"),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the quota (%s) to be created: %s", d.Id(), err)
	}

	return resourceDmsKafkaUserClientQuotaRead(ctx, d, meta)
}

func buildKafkaUserClientQuotaBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user":               utils.ValueIngoreEmpty(d.Get("user")),
		"user-default":       utils.ValueIngoreEmpty(d.Get("user_default")),
		"client":             utils.ValueIngoreEmpty(d.Get("client")),
		"client-default":     utils.ValueIngoreEmpty(d.Get("client_default")),
		"producer-byte-rate": utils.ValueIngoreEmpty(d.Get("producer_byte_rate")),
		"consumer-byte-rate": utils.ValueIngoreEmpty(d.Get("consumer_byte_rate")),
	}
	return bodyParams
}

func userClientQuotaTaskRefreshFunc(client *golangsdk.ServiceClient, instanceID string,
	d *schema.ResourceData, taskName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// getUserClientQuotaTask: query user client quota task
		getUserClientQuotaTaskHttpUrl := "v2/{project_id}/instances/{instance_id}/tasks"
		getUserClientQuotaTaskPath := client.Endpoint + getUserClientQuotaTaskHttpUrl
		getUserClientQuotaTaskPath = strings.ReplaceAll(getUserClientQuotaTaskPath, "{project_id}",
			client.ProjectID)
		getUserClientQuotaTaskPath = strings.ReplaceAll(getUserClientQuotaTaskPath, "{instance_id}", instanceID)

		getUserClientQuotaTaskPathOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getUserClientQuotaTaskPathResp, err := client.Request("GET", getUserClientQuotaTaskPath,
			&getUserClientQuotaTaskPathOpt)

		if err != nil {
			return nil, "QUERY ERROR", err
		}

		getUserClientQuotaTaskRespBody, err := utils.FlattenResponse(getUserClientQuotaTaskPathResp)
		if err != nil {
			return nil, "PARSE ERROR", err
		}

		task := filterUserClientQuotaTask(taskName, d, getUserClientQuotaTaskRespBody)
		if task == nil {
			return nil, "NIL ERROR", fmt.Errorf("can not find the task of the quota")
		}
		status := utils.PathSearch("status", task, "").(string)
		return task, status, nil
	}
}

func filterUserClientQuotaTask(taskName string, d *schema.ResourceData, resp interface{}) interface{} {
	taskJson := utils.PathSearch("tasks", resp, make([]interface{}, 0))
	taskArray := taskJson.([]interface{})
	if len(taskArray) < 1 {
		return nil
	}

	rawUser, rawUserOK := d.GetOk("user")
	rawUserDefault := d.Get("user_default").(bool)
	rawClient, rawClientOK := d.GetOk("client")
	rawClientDefault := d.Get("client_default").(bool)

	for _, task := range taskArray {
		name := utils.PathSearch("name", task, nil)
		params := utils.PathSearch("params", task, nil).(string)
		paramsData := []byte(params)
		var paramsJons interface{}
		err := json.Unmarshal(paramsData, &paramsJons)
		if err != nil {
			fmt.Println(err)
		}
		userClientQuota := utils.PathSearch("new_kafka_user_client_quota", paramsJons, nil)
		user := utils.PathSearch("user", userClientQuota, nil)
		userDefault := utils.PathSearch(`"user-default"`, userClientQuota, false).(bool)
		client := utils.PathSearch("client", userClientQuota, nil)
		clientDefault := utils.PathSearch(`"client-default"`, userClientQuota, false).(bool)
		if taskName != name {
			continue
		}
		if rawUserOK && rawUser != user {
			continue
		}
		if rawUserDefault != userDefault {
			continue
		}
		if rawClientOK && rawClient != client {
			continue
		}
		if rawClientDefault != clientDefault {
			continue
		}

		return task
	}

	return nil
}

func resourceDmsKafkaUserClientQuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getKafkaUserClientQuota: query DMS kafka user client quota
	var (
		getKafkaUserClientQuotaHttpUrl = "v2/kafka/{project_id}/instances/{instance_id}/kafka-user-client-quota"
		getKafkaUserClientQuotaProduct = "dms"
	)
	getKafkaUserClientQuotaClient, err := cfg.NewServiceClient(getKafkaUserClientQuotaProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 5 {
		return diag.Errorf("invalid id format, must be <instance_id>/<user>/<user_default>/<client>/<client_default>")
	}
	instanceID := parts[0]
	getKafkaUserClientQuotaPath := getKafkaUserClientQuotaClient.Endpoint + getKafkaUserClientQuotaHttpUrl
	getKafkaUserClientQuotaPath = strings.ReplaceAll(getKafkaUserClientQuotaPath, "{project_id}",
		getKafkaUserClientQuotaClient.ProjectID)
	getKafkaUserClientQuotaPath = strings.ReplaceAll(getKafkaUserClientQuotaPath, "{instance_id}", instanceID)

	getKafkaUserClientQuotaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getKafkaUserClientQuotaResp, err := getKafkaUserClientQuotaClient.Request("GET", getKafkaUserClientQuotaPath,
		&getKafkaUserClientQuotaOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving the quota")
	}

	getKafkaUserClientQuotaRespBody, respBodyerr := utils.FlattenResponse(getKafkaUserClientQuotaResp)
	if respBodyerr != nil {
		return diag.FromErr(respBodyerr)
	}

	quota := filterUserClientQuota(parts, getKafkaUserClientQuotaRespBody)
	if quota == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("user", utils.PathSearch("user", quota, nil)),
		d.Set("user_default", utils.PathSearch(`"user-default"`, quota, nil)),
		d.Set("client", utils.PathSearch("client", quota, nil)),
		d.Set("client_default", utils.PathSearch(`"client-default"`, quota, nil)),
		d.Set("producer_byte_rate", utils.PathSearch(`"producer-byte-rate"`, quota, nil)),
		d.Set("consumer_byte_rate", utils.PathSearch(`"consumer-byte-rate"`, quota, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func filterUserClientQuota(parts []string, resp interface{}) interface{} {
	quotaJson := utils.PathSearch("quotas", resp, make([]interface{}, 0))
	quotaArray := quotaJson.([]interface{})
	if len(quotaArray) < 1 || len(parts) != 5 {
		return nil
	}
	rawUserDefault, _ := strconv.ParseBool(parts[2])
	rawClientDefault, _ := strconv.ParseBool(parts[4])

	for _, quota := range quotaArray {
		user := utils.PathSearch("user", quota, nil)
		userDefault := utils.PathSearch(`"user-default"`, quota, false).(bool)
		client := utils.PathSearch("client", quota, nil)
		clientDefault := utils.PathSearch(`"client-default"`, quota, false).(bool)
		if parts[1] != user {
			continue
		}
		if rawUserDefault != userDefault {
			continue
		}
		if parts[3] != client {
			continue
		}
		if rawClientDefault != clientDefault {
			continue
		}
		return quota
	}
	return nil
}

func resourceDmsKafkaUserClientQuotaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	// updateKafkaUserClientQuota: update DMS kafka user client quota
	var (
		updateKafkaUserClientQuotaHttpUrl = "v2/kafka/{project_id}/instances/{instance_id}/kafka-user-client-quota"
		updateKafkaUserClientQuotaProduct = "dmsv2"
	)
	updateKafkaUserClientQuotaClient, err := cfg.NewServiceClient(updateKafkaUserClientQuotaProduct, region)

	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	updateKafkaUserClientQuotaPath := updateKafkaUserClientQuotaClient.Endpoint + updateKafkaUserClientQuotaHttpUrl
	updateKafkaUserClientQuotaPath = strings.ReplaceAll(updateKafkaUserClientQuotaPath, "{project_id}",
		updateKafkaUserClientQuotaClient.ProjectID)
	updateKafkaUserClientQuotaPath = strings.ReplaceAll(updateKafkaUserClientQuotaPath, "{instance_id}", instanceID)

	updateKafkaUserClientQuotaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateKafkaUserClientQuotaOpt.JSONBody = utils.RemoveNil(buildKafkaUserClientQuotaBodyParams(d))

	// The quota is allowd to update only when the instance status is RUNNING.
	retryFunc := func() (interface{}, bool, error) {
		updateKafkaUserClientQuotaResp, createErr := updateKafkaUserClientQuotaClient.Request("PUT",
			updateKafkaUserClientQuotaPath, &updateKafkaUserClientQuotaOpt)
		retry, err := handleMultiOperationsError(createErr)
		return updateKafkaUserClientQuotaResp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     KafkaInstanceStateRefreshFunc(updateKafkaUserClientQuotaClient, instanceID),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return diag.Errorf("error updating the quota: %v", err)
	}

	_, err = utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	// The quota modification triggers a related task, if the task status is SUCCESS, the quota has been modified.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATED"},
		Target:       []string{"SUCCESS"},
		Refresh:      userClientQuotaTaskRefreshFunc(updateKafkaUserClientQuotaClient, instanceID, d, "kafkaClientQuotaModify"),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        1 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the quota (%s) to be updated: %s", d.Id(), err)
	}

	return resourceDmsKafkaUserClientQuotaRead(ctx, d, meta)
}

func resourceDmsKafkaUserClientQuotaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteKafkaUserClientQuota: delete DMS kafka user client quota
	var (
		deleteKafkaUserClientQuotaHttpUrl = "v2/kafka/{project_id}/instances/{instance_id}/kafka-user-client-quota"
		deleteKafkaUserClientQuotaProduct = "dmsv2"
	)
	deleteKafkaUserClientQuotaClient, err := cfg.NewServiceClient(deleteKafkaUserClientQuotaProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	deleteKafkaUserClientQuotaPath := deleteKafkaUserClientQuotaClient.Endpoint + deleteKafkaUserClientQuotaHttpUrl
	deleteKafkaUserClientQuotaPath = strings.ReplaceAll(deleteKafkaUserClientQuotaPath, "{project_id}",
		deleteKafkaUserClientQuotaClient.ProjectID)
	deleteKafkaUserClientQuotaPath = strings.ReplaceAll(deleteKafkaUserClientQuotaPath, "{instance_id}", instanceID)

	deleteKafkaUserClientQuotaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteKafkaUserClientQuotaOpt.JSONBody = utils.RemoveNil(buildKafkaUserClientQuotaBodyParams(d))

	// The quota is allowd to delete only when the instance status is RUNNING.
	retryFunc := func() (interface{}, bool, error) {
		deleteKafkaUserClientQuotaResp, deleteErr := deleteKafkaUserClientQuotaClient.Request("DELETE",
			deleteKafkaUserClientQuotaPath, &deleteKafkaUserClientQuotaOpt)
		retry, err := handleOperationConflictError(deleteErr)
		return deleteKafkaUserClientQuotaResp, retry, err
	}
	_, retryErr := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     KafkaInstanceStateRefreshFunc(deleteKafkaUserClientQuotaClient, instanceID),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if retryErr != nil {
		return diag.Errorf("error deleting the quota: %v", err)
	}

	// The quota deletion triggers a related task, if the task status is SUCCESS, the quota has been deleted.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATED"},
		Target:       []string{"SUCCESS"},
		Refresh:      userClientQuotaTaskRefreshFunc(deleteKafkaUserClientQuotaClient, instanceID, d, "kafkaClientQuotaDelete"),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        1 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the quota (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}

func handleOperationConflictError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault404); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		// DMS.00404022 This instance does not exist.
		if errorCode.(string) == "DMS.00404022" {
			return true, err
		}
	}
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		// DMS.00400026 This operation is not allowed due to the instance status.
		if errorCode.(string) == "DMS.00400026" {
			return true, err
		}
	}
	return false, err
}
