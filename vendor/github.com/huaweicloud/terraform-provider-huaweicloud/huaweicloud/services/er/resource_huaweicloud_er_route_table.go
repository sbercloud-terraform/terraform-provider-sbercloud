package er

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/er/v3/associations"
	"github.com/chnsz/golangsdk/openstack/er/v3/propagations"
	"github.com/chnsz/golangsdk/openstack/er/v3/routetables"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceRouteTable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRouteTableCreate,
		UpdateContext: resourceRouteTableUpdate,
		ReadContext:   resourceRouteTableRead,
		DeleteContext: resourceRouteTableDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceRouteTableImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the ER instance and route table are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the ER instance to which the route table belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the route table.`,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5\\w.-]*$"), "The name only english and "+
						"chinese letters, digits, underscore (_), hyphens (-) and dots (.) are allowed."),
				),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the ER route table.`,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 255),
					validation.StringMatch(regexp.MustCompile(`^[^<>]*$`),
						"The angle brackets (< and >) are not allowed."),
				),
			},
			"tags": common.TagsSchema(),
			// Attributes
			"is_default_association": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether this route table is the default association route table.`,
			},
			"is_default_propagation": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether this route table is the default propagation route table.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the route table.`,
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

func buildRouteTableCreateOpts(d *schema.ResourceData) routetables.CreateOpts {
	return routetables.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}
}

func resourceRouteTableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ErV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	opts := buildRouteTableCreateOpts(d)
	resp, err := routetables.Create(client, instanceId, opts)
	if err != nil {
		return diag.Errorf("error creating route table: %s", err)
	}
	d.SetId(resp.ID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      routeTableStatusRefreshFunc(client, instanceId, d.Id(), []string{"available"}),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRouteTableRead(ctx, d, meta)
}

func routeTableStatusRefreshFunc(client *golangsdk.ServiceClient, instanceId, routeTableId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := routetables.Get(client, instanceId, routeTableId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return resp, "COMPLETED", nil
			}

			return nil, "", err
		}
		log.Printf("[DEBUG] The details of the route table (%s) is: %#v", routeTableId, resp)

		if utils.StrSliceContains([]string{"failed"}, resp.Status) {
			return resp, "", fmt.Errorf("unexpected status '%s'", resp.Status)
		}
		if utils.StrSliceContains(targets, resp.Status) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func resourceRouteTableRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	routeTableId := d.Id()
	resp, err := routetables.Get(client, instanceId, routeTableId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "ER route table")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("is_default_association", resp.IsDefaultAssociation),
		d.Set("is_default_propagation", resp.IsDefaultPropagation),
		d.Set("tags", utils.TagsToMap(resp.Tags)),
		d.Set("status", resp.Status),
		d.Set("created_at", resp.CreatedAt),
		d.Set("updated_at", resp.UpdatedAt),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving route table (%s) fields: %s", routeTableId, mErr)
	}
	return nil
}

func updateRouteTableBasicInfo(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		instanceId   = d.Get("instance_id").(string)
		routeTableId = d.Id()
	)

	opts := routetables.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: utils.String(d.Get("description").(string)),
	}

	_, err := routetables.Update(client, instanceId, routeTableId, opts)
	if err != nil {
		return fmt.Errorf("error updating route table (%s): %s", routeTableId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      routeTableStatusRefreshFunc(client, instanceId, routeTableId, []string{"available"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

func resourceRouteTableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	if d.HasChanges("name", "description") {
		if err = updateRouteTableBasicInfo(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		err = utils.UpdateResourceTags(client, d, "route-table", d.Id())
		if err != nil {
			return diag.Errorf("error updating route table tags: %s", err)
		}
	}

	return resourceRouteTableRead(ctx, d, meta)
}

func releaseRouteTableAssociations(client *golangsdk.ServiceClient, instanceId, routeTableId string) error {
	resp, err := associations.List(client, instanceId, routeTableId, associations.ListOpts{})
	if err != nil {
		return fmt.Errorf("error getting association list from the specified route table (%s): %s", routeTableId, err)
	}
	for _, association := range resp {
		opts := associations.DeleteOpts{
			AttachmentId: association.AttachmentId,
		}
		err := associations.Delete(client, instanceId, routeTableId, opts)
		if err != nil {
			return fmt.Errorf("error disable the association: %s", err)
		}
	}

	return nil
}

func releaseRouteTablePropagations(client *golangsdk.ServiceClient, instanceId, routeTableId string) error {
	resp, err := propagations.List(client, instanceId, routeTableId, propagations.ListOpts{})
	if err != nil {
		return fmt.Errorf("error getting association list from the specified route table (%s): %s", routeTableId, err)
	}
	for _, propagation := range resp {
		opts := propagations.DeleteOpts{
			AttachmentId: propagation.AttachmentId,
		}
		err := propagations.Delete(client, instanceId, routeTableId, opts)
		if err != nil {
			return fmt.Errorf("error disable the propagation: %s", err)
		}
	}

	return nil
}

func resourceRouteTableDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	routeTableId := d.Id()
	// Before delete route table, release all associations and propagations.
	err = releaseRouteTableAssociations(client, instanceId, routeTableId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = releaseRouteTablePropagations(client, instanceId, routeTableId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = routetables.Delete(client, instanceId, routeTableId)
	if err != nil {
		return diag.Errorf("error deleting route table (%s): %s", routeTableId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      routeTableStatusRefreshFunc(client, instanceId, routeTableId, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRouteTableImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid format for import ID, want '<instance_id>/<route_table_id>', but '%s'", d.Id())
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}
