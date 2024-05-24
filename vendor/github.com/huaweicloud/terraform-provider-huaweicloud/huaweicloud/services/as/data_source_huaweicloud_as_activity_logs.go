package as

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/activitylogs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API AS GET /autoscaling-api/v1/{project_id}/scaling_activity_log/{scaling_group_id}
func DataSourceActivityLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceActivityLogsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"activity_logs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"removed_instances": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deleted_instances": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"added_instances": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"current_instance_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"desire_instance_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"changes_instance_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceActivityLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating AS v1 client: %s", err)
	}

	var (
		groupID = d.Get("scaling_group_id").(string)
		opts    = activitylogs.ListOpts{
			StartTime: d.Get("start_time").(string),
			EndTime:   d.Get("end_time").(string),
		}
	)

	activityLogList, err := activitylogs.List(client, groupID, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "AS activity logs")
	}

	activityLogs := make([]map[string]interface{}, 0, len(activityLogList))
	for _, activityLog := range activityLogList {
		if val, ok := d.GetOk("status"); ok && val.(string) != activityLog.Status {
			continue
		}
		activityLogMap := map[string]interface{}{
			"id":                      activityLog.ID,
			"status":                  activityLog.Status,
			"start_time":              activityLog.StartTime,
			"end_time":                activityLog.EndTime,
			"removed_instances":       activityLog.InstanceRemovedList,
			"deleted_instances":       activityLog.InstanceDeletedList,
			"added_instances":         activityLog.InstanceAddedList,
			"current_instance_number": activityLog.InstanceValue,
			"desire_instance_number":  activityLog.DesireValue,
			"changes_instance_number": activityLog.ScalingValue,
			"description":             activityLog.Description,
		}
		activityLogs = append(activityLogs, activityLogMap)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("activity_logs", activityLogs),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving activity logs data source fields: %s", mErr)
	}
	return nil
}
