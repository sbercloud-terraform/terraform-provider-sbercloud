package ecs

import (
	"context"
	"fmt"
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

// @API ECS POST /v1/{project_id}/cloudservers/{server_id}/nics
// @API ECS POST /v1/{project_id}/cloudservers/{server_id}/nics/delete
// @API ECS GET /v1/{project_id}/cloudservers/{server_id}/os-interface
// @API VPC GET /v1/{project_id}/ports/{port_id}
// @API VPC PUT /v1/{project_id}/ports/{port_id}
// @API ECS GET /v1/{project_id}/jobs/{job_id}
func ResourceComputeInterfaceAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeInterfaceAttachCreate,
		ReadContext:   resourceComputeInterfaceAttachRead,
		UpdateContext: resourceComputeInterfaceAttachUpdate,
		DeleteContext: resourceComputeInterfaceAttachDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceComputeInterfaceAttachImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"port_id", "network_id"},
			},
			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"fixed_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source_dest_check": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"ipv6_bandwidth_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ipv6_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"fixed_ipv6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mac": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeInterfaceAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	computeClient, err := cfg.NewServiceClient("ecs", region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	// Create NIC and get `job_id`
	nic := make([]map[string]interface{}, 1)
	nic[0] = utils.RemoveNil(buildCreateNicBodyParams(d))

	createNicHttpUrl := "v1/{project_id}/cloudservers/{server_id}/nics"
	createNicPath := computeClient.Endpoint + createNicHttpUrl
	createNicPath = strings.ReplaceAll(createNicPath, "{project_id}", computeClient.ProjectID)
	createNicPath = strings.ReplaceAll(createNicPath, "{server_id}", d.Get("instance_id").(string))
	createNicOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"nics": nic,
		},
	}
	createNicResp, err := computeClient.Request("POST", createNicPath, &createNicOpt)
	if err != nil {
		return diag.Errorf("error creating ECS NIC: %s", err)
	}
	createNicRespBody, err := utils.FlattenResponse(createNicResp)
	if err != nil {
		return diag.FromErr(err)
	}
	jobID, err := jmespath.Search("job_id", createNicRespBody)
	if err != nil {
		return diag.Errorf("error creating ECS NIC: `job_id` not found in API response")
	}

	// Wait for job status become `SUCCESS`.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      getJobRefreshFunc(computeClient, jobID.(string)),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for ECS NIC created: %s", err)
	}

	// Get `nic_id`
	id, err := jmespath.Search("entities.sub_jobs|[0].entities.nic_id", result)
	if err != nil {
		return diag.Errorf("error creating ECS NIC: `nic_id` not found in API response")
	}
	d.SetId(id.(string))

	// Update port if `source_dest_check` is false, else skip update,
	// because `security_group_ids` is set when NIC is created.
	vpcClient, err := cfg.NewServiceClient("vpc", region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}
	if !d.Get("source_dest_check").(bool) {
		err = updatePort(vpcClient, d.Id(), d.Get("security_group_ids"), d.Get("source_dest_check").(bool))
		if err != nil {
			return diag.Errorf("error creating ECS NIC: %s", err)
		}
	}

	return resourceComputeInterfaceAttachRead(ctx, d, meta)
}

func buildCreateNicBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"subnet_id":       utils.ValueIngoreEmpty(d.Get("network_id")),
		"security_groups": buildNicsRequestBodySecurityGroups(d.Get("security_group_ids")),
		"ip_address":      utils.ValueIngoreEmpty(d.Get("fixed_ip")),
		"port_id":         utils.ValueIngoreEmpty(d.Get("port_id")),
		"ipv6_enable":     utils.ValueIngoreEmpty(d.Get("ipv6_enable")),
		"ipv6_bandwidth": map[string]interface{}{
			"id": utils.ValueIngoreEmpty(d.Get("ipv6_bandwidth_id")),
		},
	}
	return bodyParams
}

func buildNicsRequestBodySecurityGroups(rawParams interface{}) []map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	ids := make([]map[string]interface{}, len(rawArray))
	for i, val := range rawArray {
		id := val.(string)
		params := map[string]interface{}{
			"id": id,
		}
		ids[i] = params
	}
	return ids
}

func getJob(computeClient *golangsdk.ServiceClient, id string) (interface{}, error) {
	getJobHttpUrl := "v1/{project_id}/jobs/{job_id}"
	getJobPath := computeClient.Endpoint + getJobHttpUrl
	getJobPath = strings.ReplaceAll(getJobPath, "{project_id}", computeClient.ProjectID)
	getJobPath = strings.ReplaceAll(getJobPath, "{job_id}", id)
	getJobOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getJobResp, err := computeClient.Request("GET", getJobPath, &getJobOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getJobResp)
}

func getJobRefreshFunc(computeClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result, err := getJob(computeClient, id)
		if err != nil {
			return nil, "ERROR", err
		}
		status, err := jmespath.Search("status", result)
		if err != nil {
			return nil, "ERROR", err
		}
		if status.(string) == "FAIL" {
			err = fmt.Errorf("job failed with code %s: %s",
				utils.PathSearch("error_code", result, ""), utils.PathSearch("fail_reason", result, ""))
			return nil, "FAIL", err
		}
		if status.(string) == "SUCCESS" {
			return result, "SUCCESS", nil
		}
		return result, "PENDING", nil
	}
}

func updatePort(vpcClient *golangsdk.ServiceClient, portID string, securityGroupIds interface{}, sourceDestCheck bool) error {
	updatePortHttpUrl := "v1/{project_id}/ports/{port_id}"
	updatePortPath := vpcClient.Endpoint + updatePortHttpUrl
	updatePortPath = strings.ReplaceAll(updatePortPath, "{project_id}", vpcClient.ProjectID)
	updatePortPath = strings.ReplaceAll(updatePortPath, "{port_id}", portID)

	// Update `allowedAddressPairs` of the port to `1.1.1.1/0` to disable the source/destination check.
	allowedAddressPairs := make([]map[string]interface{}, 0)
	if !sourceDestCheck {
		allowedAddressPairs = append(allowedAddressPairs, map[string]interface{}{
			"ip_address": "1.1.1.1/0",
		})
	}

	port := make(map[string]interface{})
	port["allowed_address_pairs"] = allowedAddressPairs
	if len(securityGroupIds.([]interface{})) > 0 {
		port["security_groups"] = securityGroupIds
	}

	updatePortOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"port": port,
		},
	}

	_, err := vpcClient.Request("PUT", updatePortPath, &updatePortOpt)
	if err != nil {
		return err
	}

	return nil
}

func resourceComputeInterfaceAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	computeClient, err := cfg.NewServiceClient("ecs", region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	id := d.Id()

	// Get NIC.
	listNicsHttpUrl := "v1/{project_id}/cloudservers/{server_id}/os-interface"
	listNicsPath := computeClient.Endpoint + listNicsHttpUrl
	listNicsPath = strings.ReplaceAll(listNicsPath, "{project_id}", computeClient.ProjectID)
	listNicsPath = strings.ReplaceAll(listNicsPath, "{server_id}", d.Get("instance_id").(string))
	listNicsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listNicsResp, err := computeClient.Request("GET", listNicsPath, &listNicsOpt)
	if err != nil {
		return diag.Errorf("error getting NIC list: %s", err)
	}
	listNicsRespBody, err := utils.FlattenResponse(listNicsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPaths := fmt.Sprintf("interfaceAttachments[?port_id=='%s']|[0]", id)
	nic, err := jmespath.Search(jsonPaths, listNicsRespBody)
	if err != nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	// Getting VPC port.
	vpcClient, err := cfg.NewServiceClient("vpc", region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}
	port, err := readVPCPort(vpcClient, id)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("port_id", id),
		d.Set("network_id", utils.PathSearch("net_id", nic, nil)),
		d.Set("fixed_ip", utils.PathSearch("fixed_ips|[0].ip_address", nic, nil)),
		d.Set("fixed_ipv6", utils.PathSearch("fixed_ips|[1].ip_address", nic, nil)),
		d.Set("ipv6_enable", len(utils.PathSearch("fixed_ips", nic, make([]interface{}, 0)).([]interface{})) == 2),
		d.Set("mac", utils.PathSearch("mac_addr", nic, nil)),
		d.Set("security_group_ids", utils.PathSearch("port.security_groups", port, make([]interface{}, 0))),
		d.Set("source_dest_check", flattenSourceDestCheck(utils.PathSearch("port.allowed_address_pairs", port, make([]interface{}, 0)))),
		d.Set("ipv6_bandwidth_id", utils.PathSearch("port.ipv6_bandwidth_id", port, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func readVPCPort(vpcClient *golangsdk.ServiceClient, portID string) (interface{}, error) {
	getPortHttpUrl := "v1/{project_id}/ports/{port_id}"
	getPortPath := vpcClient.Endpoint + getPortHttpUrl
	getPortPath = strings.ReplaceAll(getPortPath, "{project_id}", vpcClient.ProjectID)
	getPortPath = strings.ReplaceAll(getPortPath, "{port_id}", portID)
	getPortOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPortResp, err := vpcClient.Request("GET", getPortPath, &getPortOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving VPC port: %s", err)
	}
	return utils.FlattenResponse(getPortResp)
}

func flattenSourceDestCheck(allowedAddressPairs interface{}) bool {
	pairs := allowedAddressPairs.([]interface{})
	return len(pairs) == 0
}

func resourceComputeInterfaceAttachUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	vpcClient, err := cfg.NewServiceClient("vpc", region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	err = updatePort(vpcClient, d.Id(), d.Get("security_group_ids"), d.Get("source_dest_check").(bool))
	if err != nil {
		return diag.Errorf("error updating ECS NIC port (%s): %s", d.Id(), err)
	}

	return resourceComputeInterfaceAttachRead(ctx, d, meta)
}

func resourceComputeInterfaceAttachDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	computeClient, err := cfg.NewServiceClient("ecs", region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	serverID := d.Get("instance_id").(string)

	deleteNicsHttpUrl := "v1/{project_id}/cloudservers/{server_id}/nics/delete"
	deleteNicsPath := computeClient.Endpoint + deleteNicsHttpUrl
	deleteNicsPath = strings.ReplaceAll(deleteNicsPath, "{project_id}", computeClient.ProjectID)
	deleteNicsPath = strings.ReplaceAll(deleteNicsPath, "{server_id}", serverID)

	nic := make([]map[string]interface{}, 0)
	nic = append(nic, map[string]interface{}{
		"id": d.Get("port_id"),
	})
	deleteNicOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"nics": nic,
		},
	}
	// The `DELETE` API use `POST` method actually.
	deleteNicResp, err := computeClient.Request("POST", deleteNicsPath, &deleteNicOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	deleteNicRespBody, err := utils.FlattenResponse(deleteNicResp)
	if err != nil {
		return diag.FromErr(err)
	}
	jobID, err := jmespath.Search("job_id", deleteNicRespBody)
	if err != nil {
		return diag.Errorf("error deleting ECS NIC: `job_id` not found in API response")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      getJobRefreshFunc(computeClient, jobID.(string)),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for the ECS NIC to be deleted: %s", err)
	}

	return nil
}

func resourceComputeInterfaceAttachImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	partLength := len(parts)

	if partLength == 2 {
		d.SetId(parts[1])
		return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
	}
	return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<port_id>")
}
