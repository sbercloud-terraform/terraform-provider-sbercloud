// Generated by PMS #163
package css

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/filters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceCssLogstashPipelineActions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCssLogstashPipelineActionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the CSS logstash cluster.`,
			},
			"action_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the action.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the action.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the action.`,
			},
			"actions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the actions.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the action.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the action.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the action.`,
						},
						"conf_content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The configuration file content.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time.`,
						},
						"error_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The error message of the action.`,
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The message of the action.`,
						},
					},
				},
			},
		},
	}
}

type LogstashPipelineActionsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newLogstashPipelineActionsDSWrapper(d *schema.ResourceData, meta interface{}) *LogstashPipelineActionsDSWrapper {
	return &LogstashPipelineActionsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCssLogstashPipelineActionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newLogstashPipelineActionsDSWrapper(d, meta)
	listActionsRst, err := wrapper.ListActions()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listActionsToSchema(listActionsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/lgsconf/listactions
func (w *LogstashPipelineActionsDSWrapper) ListActions() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "css")
	if err != nil {
		return nil, err
	}

	uri := "/v1.0/{project_id}/clusters/{cluster_id}/lgsconf/listactions"
	uri = strings.ReplaceAll(uri, "{cluster_id}", w.Get("cluster_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Filter(
			filters.New().From("actions").
				Where("actionType", "=", w.Get("type")).
				Where("status", "=", w.Get("status")).
				Where("id", "=", w.Get("action_id")),
		).
		OkCode(200).
		Request().
		Result()
}

func (w *LogstashPipelineActionsDSWrapper) listActionsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("actions", schemas.SliceToList(body.Get("actions"),
			func(action gjson.Result) any {
				return map[string]any{
					"id":           action.Get("id").Value(),
					"type":         action.Get("actionType").Value(),
					"status":       action.Get("status").Value(),
					"conf_content": action.Get("confContent").Value(),
					"updated_at":   action.Get("updateAt").Value(),
					"error_msg":    action.Get("errorMsg").Value(),
					"message":      action.Get("message").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
