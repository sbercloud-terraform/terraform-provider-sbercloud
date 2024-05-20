// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DLI
// ---------------------------------------------------------------

package dli

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceDatasourceConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatasourceConnectionCreate,
		UpdateContext: resourceDatasourceConnectionUpdate,
		ReadContext:   resourceDatasourceConnectionRead,
		DeleteContext: resourceDatasourceConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				ForceNew:    true,
				Description: `The name of a datasource connection.`,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z-_0-9]*$`),
					"the input is invalid"),
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The VPC ID of the service to be connected.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The subnet ID of the service to be connected.`,
			},
			"route_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The route table ID associated with the subnet of the service to be connected.`,
			},
			"queues": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `List of queue names that are available for datasource connections.`,
			},
			"hosts": {
				Type:        schema.TypeList,
				Elem:        datasourceConnectionHostSchema(),
				Optional:    true,
				Computed:    true,
				Description: `The user-defined host information. A maximum of 20,000 records are supported.`,
			},
			"routes": {
				Type:        schema.TypeSet,
				Elem:        datasourceConnectionRouteSchema(),
				Optional:    true,
				Computed:    true,
				Description: `List of routes.`,
			},
			"tags": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				ForceNew:    true,
				Description: `The key/value pairs to associate with the datasource connection.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The connection status.`,
			},
		},
	}
}

func datasourceConnectionHostSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The user-defined host name.`,
			},
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `IPv4 address of the host.`,
			},
		},
	}
	return &sc
}

func datasourceConnectionRouteSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  `The route Name`,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"cidr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The CIDR of the route.`,
			},
		},
	}
	return &sc
}

func resourceDatasourceConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createDatasourceConnection: create a DLI enhanced connection.
	var (
		createDatasourceConnectionHttpUrl = "v2.0/{project_id}/datasource/enhanced-connections"
		createDatasourceConnectionProduct = "dli"
	)
	createDatasourceConnectionClient, err := cfg.NewServiceClient(createDatasourceConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	createDatasourceConnectionPath := createDatasourceConnectionClient.Endpoint + createDatasourceConnectionHttpUrl
	createDatasourceConnectionPath = strings.ReplaceAll(createDatasourceConnectionPath, "{project_id}",
		createDatasourceConnectionClient.ProjectID)

	createDatasourceConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createDatasourceConnectionOpt.JSONBody = utils.RemoveNil(buildCreateDatasourceConnectionBodyParams(d))
	createDatasourceConnectionResp, err := createDatasourceConnectionClient.Request("POST",
		createDatasourceConnectionPath, &createDatasourceConnectionOpt)

	if err != nil {
		return diag.Errorf("error creating DatasourceConnection: %s", err)
	}

	createDatasourceConnectionRespBody, err := utils.FlattenResponse(createDatasourceConnectionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("connection_id", createDatasourceConnectionRespBody)
	if err != nil {
		return diag.Errorf("error creating DatasourceConnection: ID is not found in API response")
	}
	d.SetId(id.(string))

	// add routes
	if v, ok := d.GetOk("routes"); ok {
		err = addRoutes(createDatasourceConnectionClient, d.Id(), v.(*schema.Set))
		if err != nil {
			return diag.Errorf("error adding routes to DatasourceConnection: %s", d.Id())
		}
	}

	return resourceDatasourceConnectionRead(ctx, d, meta)
}

func buildCreateDatasourceConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":            utils.ValueIngoreEmpty(d.Get("name")),
		"dest_vpc_id":     utils.ValueIngoreEmpty(d.Get("vpc_id")),
		"dest_network_id": utils.ValueIngoreEmpty(d.Get("subnet_id")),
		"routetable_id":   utils.ValueIngoreEmpty(d.Get("route_table_id")),
		"queues":          d.Get("queues").(*schema.Set).List(),
		"hosts":           buildCreateDatasourceConnectionRequestBodyHost(d.Get("hosts")),
		"tags":            utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	return bodyParams
}

func buildCreateDatasourceConnectionRequestBodyHost(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name": utils.ValueIngoreEmpty(raw["name"]),
				"ip":   utils.ValueIngoreEmpty(raw["ip"]),
			}
		}
		return rst
	}
	return nil
}

func resourceDatasourceConnectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDatasourceConnection: Query the DLI instance
	var (
		getDatasourceConnectionHttpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{id}"
		getDatasourceConnectionProduct = "dli"
	)
	getDatasourceConnectionClient, err := cfg.NewServiceClient(getDatasourceConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	getDatasourceConnectionPath := getDatasourceConnectionClient.Endpoint + getDatasourceConnectionHttpUrl
	getDatasourceConnectionPath = strings.ReplaceAll(getDatasourceConnectionPath, "{project_id}", getDatasourceConnectionClient.ProjectID)
	getDatasourceConnectionPath = strings.ReplaceAll(getDatasourceConnectionPath, "{id}", d.Id())

	getDatasourceConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getDatasourceConnectionResp, err := getDatasourceConnectionClient.Request("GET", getDatasourceConnectionPath, &getDatasourceConnectionOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DatasourceConnection")
	}

	getDatasourceConnectionRespBody, err := utils.FlattenResponse(getDatasourceConnectionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if utils.PathSearch("status", getDatasourceConnectionRespBody, "") == "DELETED" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "the datasource connection has been deleted")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getDatasourceConnectionRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("dest_vpc_id", getDatasourceConnectionRespBody, nil)),
		d.Set("subnet_id", utils.PathSearch("dest_network_id", getDatasourceConnectionRespBody, nil)),
		d.Set("queues", utils.PathSearch("available_queue_info[*].name", getDatasourceConnectionRespBody, nil)),
		d.Set("hosts", flattenGetDatasourceConnectionResponseBodyHost(getDatasourceConnectionRespBody)),
		d.Set("routes", flattenGetDatasourceConnectionResponseBodyRoute(getDatasourceConnectionRespBody)),
		d.Set("status", utils.PathSearch("status", getDatasourceConnectionRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetDatasourceConnectionResponseBodyHost(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("hosts", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
			"ip":   utils.PathSearch("ip", v, nil),
		})
	}
	return rst
}

func flattenGetDatasourceConnectionResponseBodyRoute(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("routes", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
			"cidr": utils.PathSearch("cidr", v, nil),
		})
	}
	return rst
}

func resourceDatasourceConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateDatasourceConnectionHostsChanges := []string{
		"hosts",
	}

	if d.HasChanges(updateDatasourceConnectionHostsChanges...) {
		// updateDatasourceConnectionHosts: update hosts
		var (
			updateDatasourceConnectionHostsHttpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{id}"
			updateDatasourceConnectionHostsProduct = "dli"
		)
		updateDatasourceConnectionHostsClient, err := cfg.NewServiceClient(updateDatasourceConnectionHostsProduct, region)
		if err != nil {
			return diag.Errorf("error creating DLI Client: %s", err)
		}

		updateDatasourceConnectionHostsPath := updateDatasourceConnectionHostsClient.Endpoint + updateDatasourceConnectionHostsHttpUrl
		updateDatasourceConnectionHostsPath = strings.ReplaceAll(updateDatasourceConnectionHostsPath, "{project_id}",
			updateDatasourceConnectionHostsClient.ProjectID)
		updateDatasourceConnectionHostsPath = strings.ReplaceAll(updateDatasourceConnectionHostsPath, "{id}", d.Id())

		updateDatasourceConnectionHostsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateDatasourceConnectionHostsOpt.JSONBody = utils.RemoveNil(buildUpdateDatasourceConnectionHostsBodyParams(d, cfg))
		_, err = updateDatasourceConnectionHostsClient.Request("PUT", updateDatasourceConnectionHostsPath, &updateDatasourceConnectionHostsOpt)
		if err != nil {
			return diag.Errorf("error updating DatasourceConnection: %s", err)
		}
	}

	// updateDatasourceConnectionQueues: update queues
	updateDatasourceConnectionQueuesChanges := []string{
		"queues",
	}

	if d.HasChanges(updateDatasourceConnectionQueuesChanges...) {
		o, n := d.GetChange("queues")

		addRaws := n.(*schema.Set).Difference(o.(*schema.Set))
		delRaws := o.(*schema.Set).Difference(n.(*schema.Set))
		if addRaws.Len() > 0 {
			var (
				updateDatasourceConnectionQueuesHttpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{id}/associate-queue"
				updateDatasourceConnectionQueuesProduct = "dli"
			)
			updateDatasourceConnectionQueuesClient, err := cfg.NewServiceClient(updateDatasourceConnectionQueuesProduct, region)
			if err != nil {
				return diag.Errorf("error creating DLI Client: %s", err)
			}

			updateDatasourceConnectionQueuesPath := updateDatasourceConnectionQueuesClient.Endpoint + updateDatasourceConnectionQueuesHttpUrl
			updateDatasourceConnectionQueuesPath = strings.ReplaceAll(updateDatasourceConnectionQueuesPath, "{project_id}",
				updateDatasourceConnectionQueuesClient.ProjectID)
			updateDatasourceConnectionQueuesPath = strings.ReplaceAll(updateDatasourceConnectionQueuesPath, "{id}", d.Id())

			updateDatasourceConnectionQueuesOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			updateDatasourceConnectionQueuesOpt.JSONBody = buildDatasourceConnectionQueuesBodyParams(addRaws)
			_, err = updateDatasourceConnectionQueuesClient.Request("POST", updateDatasourceConnectionQueuesPath,
				&updateDatasourceConnectionQueuesOpt)
			if err != nil {
				return diag.Errorf("error updating DatasourceConnection: %s", err)
			}
		}

		if delRaws.Len() > 0 {
			var (
				deleteDatasourceConnectionQueuesHttpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{id}/disassociate-queue"
				deleteDatasourceConnectionQueuesProduct = "dli"
			)
			deleteDatasourceConnectionQueuesClient, err := cfg.NewServiceClient(deleteDatasourceConnectionQueuesProduct, region)
			if err != nil {
				return diag.Errorf("error creating DLI Client: %s", err)
			}

			deleteDatasourceConnectionQueuesPath := deleteDatasourceConnectionQueuesClient.Endpoint + deleteDatasourceConnectionQueuesHttpUrl
			deleteDatasourceConnectionQueuesPath = strings.ReplaceAll(deleteDatasourceConnectionQueuesPath, "{project_id}",
				deleteDatasourceConnectionQueuesClient.ProjectID)
			deleteDatasourceConnectionQueuesPath = strings.ReplaceAll(deleteDatasourceConnectionQueuesPath, "{id}", d.Id())

			deleteDatasourceConnectionQueuesOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			deleteDatasourceConnectionQueuesOpt.JSONBody = buildDatasourceConnectionQueuesBodyParams(delRaws)
			_, err = deleteDatasourceConnectionQueuesClient.Request("POST", deleteDatasourceConnectionQueuesPath,
				&deleteDatasourceConnectionQueuesOpt)
			if err != nil {
				return diag.Errorf("error updating DatasourceConnection: %s", err)
			}
		}
	}

	// updateDatasourceConnectionRoutes: update routes
	updateDatasourceConnectionRoutesChanges := []string{
		"routes",
	}

	if d.HasChanges(updateDatasourceConnectionRoutesChanges...) {
		connectionRouteClient, err := cfg.NewServiceClient("dli", region)
		if err != nil {
			return diag.Errorf("error creating DLI Client: %s", err)
		}

		o, n := d.GetChange("routes")
		addRaws := n.(*schema.Set).Difference(o.(*schema.Set))
		delRaws := o.(*schema.Set).Difference(n.(*schema.Set))

		if addRaws.Len() > 0 {
			err := addRoutes(connectionRouteClient, d.Id(), addRaws)
			if err != nil {
				return diag.Errorf("error updating DatasourceConnection: %s", err)
			}
		}

		if delRaws.Len() > 0 {
			err := removeRoutes(connectionRouteClient, d.Id(), delRaws)
			if err != nil {
				return diag.Errorf("error updating DatasourceConnection: %s", err)
			}
		}
	}
	return resourceDatasourceConnectionRead(ctx, d, meta)
}

func addRoutes(connectionRouteClient *golangsdk.ServiceClient, id string, addRaws *schema.Set) error {
	var (
		addConnectionRouteHttpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{id}/routes"
	)

	addConnectionRoutePath := connectionRouteClient.Endpoint + addConnectionRouteHttpUrl
	addConnectionRoutePath = strings.ReplaceAll(addConnectionRoutePath, "{project_id}", connectionRouteClient.ProjectID)
	addConnectionRoutePath = strings.ReplaceAll(addConnectionRoutePath, "{id}", id)

	addConnectionRouteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	for _, params := range addRaws.List() {
		addConnectionRouteOpt.JSONBody = params
		_, err := connectionRouteClient.Request("POST", addConnectionRoutePath, &addConnectionRouteOpt)
		if err != nil {
			return err
		}
	}
	return nil
}

func removeRoutes(connectionRouteClient *golangsdk.ServiceClient, id string, raws *schema.Set) error {
	for _, params := range raws.List() {
		var (
			removeDatasourceConnectionRoutesHttpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{id}/routes/{name}"
		)
		removeDatasourceConnectionRoutesPath := connectionRouteClient.Endpoint + removeDatasourceConnectionRoutesHttpUrl
		removeDatasourceConnectionRoutesPath = strings.ReplaceAll(removeDatasourceConnectionRoutesPath, "{project_id}",
			connectionRouteClient.ProjectID)
		removeDatasourceConnectionRoutesPath = strings.ReplaceAll(removeDatasourceConnectionRoutesPath, "{id}", id)
		removeDatasourceConnectionRoutesPath = strings.ReplaceAll(removeDatasourceConnectionRoutesPath, "{name}",
			fmt.Sprintf("%v", utils.PathSearch("name", params, nil)))

		removeDatasourceConnectionRoutesOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		_, err := connectionRouteClient.Request("DELETE", removeDatasourceConnectionRoutesPath,
			&removeDatasourceConnectionRoutesOpt)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildDatasourceConnectionQueuesBodyParams(v *schema.Set) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"queues": v.List(),
	}
	return bodyParams
}

func buildUpdateDatasourceConnectionHostsBodyParams(d *schema.ResourceData, _ *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"hosts": buildUpdateDatasourceConnectionHostsRequestBodyHost(d.Get("hosts")),
	}
	return bodyParams
}

func buildUpdateDatasourceConnectionHostsRequestBodyHost(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name": utils.ValueIngoreEmpty(raw["name"]),
				"ip":   utils.ValueIngoreEmpty(raw["ip"]),
			}
		}
		return rst
	}
	return nil
}

func resourceDatasourceConnectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteDatasourceConnectionHttpUrl = "v2.0/{project_id}/datasource/enhanced-connections/{id}"
		deleteDatasourceConnectionProduct = "dli"
	)
	deleteDatasourceConnectionClient, err := cfg.NewServiceClient(deleteDatasourceConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	deleteDatasourceConnectionPath := deleteDatasourceConnectionClient.Endpoint + deleteDatasourceConnectionHttpUrl
	deleteDatasourceConnectionPath = strings.ReplaceAll(deleteDatasourceConnectionPath, "{project_id}", deleteDatasourceConnectionClient.ProjectID)
	deleteDatasourceConnectionPath = strings.ReplaceAll(deleteDatasourceConnectionPath, "{id}", d.Id())

	deleteDatasourceConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deleteDatasourceConnectionClient.Request("DELETE", deleteDatasourceConnectionPath, &deleteDatasourceConnectionOpt)
	if err != nil {
		return diag.Errorf("error deleting DatasourceConnection: %s", err)
	}

	return nil
}
