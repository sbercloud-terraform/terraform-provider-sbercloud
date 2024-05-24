// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ER
// ---------------------------------------------------------------

package er

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceCreate,
		UpdateContext: resourceInstanceUpdate,
		ReadContext:   resourceInstanceRead,
		DeleteContext: resourceInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The router name.`,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5\\w.-]*$"), "The name only english and "+
						"chinese letters, digits, underscore (_) and hyphens (-) are allowed."),
				),
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `The availability zone list where the Enterprise router is located.`,
			},
			"asn": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `The BGP AS number of the Enterprise router.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the Enterprise router.`,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 255),
					validation.StringMatch(regexp.MustCompile(`^[^<>]*$`),
						"The angle brackets (< and >) are not allowed."),
				),
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The enterprise project ID to which the Enterprise router belongs.`,
			},
			"enable_default_propagation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to enable the propagation of the default route table.`,
			},
			"enable_default_association": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to enable the association of the default route table.`,
			},
			"auto_accept_shared_attachments": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to automatically accept the creation of shared attachment.`,
			},
			"default_propagation_route_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the default propagation route table.`,
			},
			"default_association_route_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the default association route table.`,
			},
			"tags": common.TagsSchema(),
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Current status of the router.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time.`,
			},
		},
	}
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createInstance: Create an Enterprise router instance.
	createInstanceHttpUrl := "enterprise-router/instances"
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 Client: %s", err)
	}

	createInstancePath := client.ResourceBaseURL() + createInstanceHttpUrl

	createInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	createInstanceOpt.JSONBody = utils.RemoveNil(buildCreateInstanceBodyParams(d, cfg))
	createInstanceResp, err := client.Request("POST", createInstancePath, &createInstanceOpt)
	if err != nil {
		return diag.Errorf("error creating Instance: %s", err)
	}

	createInstanceRespBody, err := utils.FlattenResponse(createInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("instance.id", createInstanceRespBody)
	if err != nil {
		return diag.Errorf("error creating Instance: ID is not found in API response")
	}
	d.SetId(id.(string))

	err = instanceWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Create of Instance (%s) to complete: %s", d.Id(), err)
	}
	// updateInstanceDefaultRouteTables: Update default route tables
	var (
		updateInstanceDefaultRouteTablesHttpUrl = "v3/{project_id}/enterprise-router/instances/{id}"
		updateInstanceDefaultRouteTablesProduct = "er"
	)
	updateInstanceDefaultRouteTablesClient, err := cfg.NewServiceClient(updateInstanceDefaultRouteTablesProduct, region)
	if err != nil {
		return diag.Errorf("error creating Instance Client: %s", err)
	}

	updateInstanceDefaultRouteTablesPath := updateInstanceDefaultRouteTablesClient.Endpoint + updateInstanceDefaultRouteTablesHttpUrl
	updateInstanceDefaultRouteTablesPath = strings.ReplaceAll(updateInstanceDefaultRouteTablesPath, "{project_id}",
		updateInstanceDefaultRouteTablesClient.ProjectID)
	updateInstanceDefaultRouteTablesPath = strings.ReplaceAll(updateInstanceDefaultRouteTablesPath, "{id}", d.Id())

	updateInstanceDefaultRouteTablesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	updateInstanceDefaultRouteTablesOpt.JSONBody = utils.RemoveNil(buildUpdateInstanceDefaultRouteTablesBodyParams(d))
	updateInstanceDefaultRouteTablesResp, err := updateInstanceDefaultRouteTablesClient.Request("PUT", updateInstanceDefaultRouteTablesPath,
		&updateInstanceDefaultRouteTablesOpt)
	if err != nil {
		return diag.Errorf("error creating Instance: %s", err)
	}

	_, err = utils.FlattenResponse(updateInstanceDefaultRouteTablesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	err = utils.UpdateResourceTags(client, d, "instance", d.Id())
	if err != nil {
		return diag.Errorf("error creating instance tags: %s", err)
	}

	return resourceInstanceRead(ctx, d, meta)
}

func buildCreateInstanceBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance": buildCreateInstanceInstanceChildBody(d, cfg),
	}
	return bodyParams
}

func buildCreateInstanceInstanceChildBody(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	params := map[string]interface{}{
		"name":                           utils.ValueIngoreEmpty(d.Get("name")),
		"availability_zone_ids":          utils.ValueIngoreEmpty(d.Get("availability_zones")),
		"asn":                            utils.ValueIngoreEmpty(d.Get("asn")),
		"description":                    utils.ValueIngoreEmpty(d.Get("description")),
		"enterprise_project_id":          utils.ValueIngoreEmpty(common.GetEnterpriseProjectID(d, cfg)),
		"enable_default_propagation":     utils.ValueIngoreEmpty(d.Get("enable_default_propagation")),
		"enable_default_association":     utils.ValueIngoreEmpty(d.Get("enable_default_association")),
		"auto_accept_shared_attachments": utils.ValueIngoreEmpty(d.Get("auto_accept_shared_attachments")),
	}
	return params
}

func buildUpdateInstanceDefaultRouteTablesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance": buildUpdateInstanceDefaultRouteTablesInstanceChildBody(d),
	}
	return bodyParams
}

func buildUpdateInstanceDefaultRouteTablesInstanceChildBody(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"default_propagation_route_table_id": utils.ValueIngoreEmpty(d.Get("default_propagation_route_table_id")),
		"default_association_route_table_id": utils.ValueIngoreEmpty(d.Get("default_association_route_table_id")),
	}
	return params
}

func instanceWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			// createInstanceWaiting: Query the Enterprise router instance status
			var (
				createInstanceWaitingHttpUrl = "v3/{project_id}/enterprise-router/instances/{id}"
				createInstanceWaitingProduct = "er"
			)
			createInstanceWaitingClient, err := cfg.NewServiceClient(createInstanceWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating Instance Client: %s", err)
			}

			createInstanceWaitingPath := createInstanceWaitingClient.Endpoint + createInstanceWaitingHttpUrl
			createInstanceWaitingPath = strings.ReplaceAll(createInstanceWaitingPath, "{project_id}", createInstanceWaitingClient.ProjectID)
			createInstanceWaitingPath = strings.ReplaceAll(createInstanceWaitingPath, "{id}", d.Id())

			createInstanceWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			createInstanceWaitingResp, err := createInstanceWaitingClient.Request("GET", createInstanceWaitingPath, &createInstanceWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createInstanceWaitingRespBody, err := utils.FlattenResponse(createInstanceWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			status, err := jmespath.Search(`instance.state`, createInstanceWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `instance.state`)
			}

			if utils.StrSliceContains([]string{"fail"}, status.(string)) {
				return createInstanceWaitingRespBody, "", fmt.Errorf("unexpected status '%s'", status.(string))
			}

			if utils.StrSliceContains([]string{"available"}, status.(string)) {
				return createInstanceWaitingRespBody, "COMPLETED", nil
			}

			return createInstanceWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getInstance: Query the Enterprise router instance detail
	var (
		getInstanceHttpUrl = "v3/{project_id}/enterprise-router/instances/{id}"
		getInstanceProduct = "er"
	)
	getInstanceClient, err := cfg.NewServiceClient(getInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating Instance Client: %s", err)
	}

	getInstancePath := getInstanceClient.Endpoint + getInstanceHttpUrl
	getInstancePath = strings.ReplaceAll(getInstancePath, "{project_id}", getInstanceClient.ProjectID)
	getInstancePath = strings.ReplaceAll(getInstancePath, "{id}", d.Id())

	getInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getInstanceResp, err := getInstanceClient.Request("GET", getInstancePath, &getInstanceOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Instance")
	}

	getInstanceRespBody, err := utils.FlattenResponse(getInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("instance.name", getInstanceRespBody, nil)),
		d.Set("description", utils.PathSearch("instance.description", getInstanceRespBody, nil)),
		d.Set("status", utils.PathSearch("instance.state", getInstanceRespBody, nil)),
		d.Set("created_at", utils.PathSearch("instance.created_at", getInstanceRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("instance.updated_at", getInstanceRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("instance.enterprise_project_id", getInstanceRespBody, nil)),
		d.Set("asn", utils.PathSearch("instance.asn", getInstanceRespBody, nil)),
		d.Set("enable_default_propagation", utils.PathSearch("instance.enable_default_propagation", getInstanceRespBody, nil)),
		d.Set("enable_default_association", utils.PathSearch("instance.enable_default_association", getInstanceRespBody, nil)),
		d.Set("default_propagation_route_table_id", utils.PathSearch("instance.default_propagation_route_table_id", getInstanceRespBody, nil)),
		d.Set("default_association_route_table_id", utils.PathSearch("instance.default_association_route_table_id", getInstanceRespBody, nil)),
		d.Set("availability_zones", utils.PathSearch("instance.availability_zone_ids", getInstanceRespBody, nil)),
		d.Set("auto_accept_shared_attachments", utils.PathSearch("instance.auto_accept_shared_attachments", getInstanceRespBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("instance.tags", getInstanceRespBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 Client: %s", err)
	}

	updateInstancehasChanges := []string{
		"name",
		"description",
		"enable_default_propagation",
		"enable_default_association",
		"default_propagation_route_table_id",
		"default_association_route_table_id",
		"auto_accept_shared_attachments",
	}

	if d.HasChanges(updateInstancehasChanges...) {
		// updateInstance: Update the configuration of Enterprise router instance
		updateInstanceHttpUrl := "enterprise-router/instances/{id}"
		updateInstancePath := client.ResourceBaseURL() + updateInstanceHttpUrl
		updateInstancePath = strings.ReplaceAll(updateInstancePath, "{project_id}", client.ProjectID)
		updateInstancePath = strings.ReplaceAll(updateInstancePath, "{id}", d.Id())

		updateInstanceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateInstanceOpt.JSONBody = utils.RemoveNil(buildUpdateInstanceBodyParams(d))
		_, err = client.Request("PUT", updateInstancePath, &updateInstanceOpt)
		if err != nil {
			return diag.Errorf("error updating Instance: %s", err)
		}
		err = instanceWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the Update of Instance (%s) to complete: %s", d.Id(), err)
		}
	}

	updateInstanceAvailabilityZoneshasChanges := []string{
		"availability_zones",
	}

	if d.HasChanges(updateInstanceAvailabilityZoneshasChanges...) {
		// updateInstanceAvailabilityZones: Update the availability zone list where the Enterprise router instance is located
		updateInstanceAvailabilityZonesHttpUrl := "enterprise-router/instances/{id}/change-availability-zone-ids"
		updateInstanceAvailabilityZonesPath := client.ResourceBaseURL() + updateInstanceAvailabilityZonesHttpUrl
		updateInstanceAvailabilityZonesPath = strings.ReplaceAll(updateInstanceAvailabilityZonesPath, "{project_id}", client.ProjectID)
		updateInstanceAvailabilityZonesPath = strings.ReplaceAll(updateInstanceAvailabilityZonesPath, "{id}", d.Id())

		updateInstanceAvailabilityZonesOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				202,
			},
		}
		updateInstanceAvailabilityZonesOpt.JSONBody = utils.RemoveNil(buildUpdateInstanceAvailabilityZonesBodyParams(d))
		_, err = client.Request("POST", updateInstanceAvailabilityZonesPath, &updateInstanceAvailabilityZonesOpt)
		if err != nil {
			return diag.Errorf("error updating Instance: %s", err)
		}
		err = instanceWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the Update of Instance (%s) to complete: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		err = utils.UpdateResourceTags(client, d, "instance", d.Id())
		if err != nil {
			return diag.Errorf("error updating instance tags: %s", err)
		}
	}
	return resourceInstanceRead(ctx, d, meta)
}

func buildUpdateInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance": buildUpdateInstanceInstanceChildBody(d),
	}
	return bodyParams
}

func buildUpdateInstanceInstanceChildBody(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"name":                               utils.ValueIngoreEmpty(d.Get("name")),
		"description":                        utils.ValueIngoreEmpty(d.Get("description")),
		"enable_default_propagation":         utils.ValueIngoreEmpty(d.Get("enable_default_propagation")),
		"enable_default_association":         utils.ValueIngoreEmpty(d.Get("enable_default_association")),
		"default_propagation_route_table_id": utils.ValueIngoreEmpty(d.Get("default_propagation_route_table_id")),
		"default_association_route_table_id": utils.ValueIngoreEmpty(d.Get("default_association_route_table_id")),
		"auto_accept_shared_attachments":     utils.ValueIngoreEmpty(d.Get("auto_accept_shared_attachments")),
	}
	return params
}

func buildUpdateInstanceAvailabilityZonesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"availability_zone_ids": utils.ValueIngoreEmpty(d.Get("availability_zones")),
	}
	return bodyParams
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteInstance: Deleter an existing router instance
	var (
		deleteInstanceHttpUrl = "v3/{project_id}/enterprise-router/instances/{id}"
		deleteInstanceProduct = "er"
	)
	deleteInstanceClient, err := cfg.NewServiceClient(deleteInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating Instance Client: %s", err)
	}

	deleteInstancePath := deleteInstanceClient.Endpoint + deleteInstanceHttpUrl
	deleteInstancePath = strings.ReplaceAll(deleteInstancePath, "{project_id}", deleteInstanceClient.ProjectID)
	deleteInstancePath = strings.ReplaceAll(deleteInstancePath, "{id}", d.Id())

	deleteInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	_, err = deleteInstanceClient.Request("DELETE", deleteInstancePath, &deleteInstanceOpt)
	if err != nil {
		return diag.Errorf("error deleting Instance: %s", err)
	}

	err = deleteInstanceWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the Delete of Instance (%s) to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteInstanceWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// deleteInstanceWaiting: missing operation notes
			var (
				deleteInstanceWaitingHttpUrl = "v3/{project_id}/enterprise-router/instances/{id}"
				deleteInstanceWaitingProduct = "er"
			)
			deleteInstanceWaitingClient, err := config.NewServiceClient(deleteInstanceWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating Instance Client: %s", err)
			}

			deleteInstanceWaitingPath := deleteInstanceWaitingClient.Endpoint + deleteInstanceWaitingHttpUrl
			deleteInstanceWaitingPath = strings.ReplaceAll(deleteInstanceWaitingPath, "{project_id}", deleteInstanceWaitingClient.ProjectID)
			deleteInstanceWaitingPath = strings.ReplaceAll(deleteInstanceWaitingPath, "{id}", d.Id())

			deleteInstanceWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			deleteInstanceWaitingResp, err := deleteInstanceWaitingClient.Request("GET", deleteInstanceWaitingPath, &deleteInstanceWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return deleteInstanceWaitingResp, "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			deleteInstanceWaitingRespBody, err := utils.FlattenResponse(deleteInstanceWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			status, err := jmespath.Search(`instance.state`, deleteInstanceWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `instance.state`)
			}

			if utils.StrSliceContains([]string{"fail"}, status.(string)) {
				return deleteInstanceWaitingRespBody, "", fmt.Errorf("unexpected status '%s'", status.(string))
			}

			return deleteInstanceWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
