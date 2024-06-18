// Generated by PMS #136
package waf

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceWafRulesCcProtection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWafRulesCcProtectionRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the policy to which the cc protection rules belong.`,
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the cc protection rule.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the cc protection rule.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the cc protection rule.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID to which the protection policy belongs.`,
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of cc protection rules.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the cc protection rule.`,
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the policy to which the cc protection rule belongs.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the cc protection rule.`,
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The status of the cc protection rule.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the cc protection rule.`,
						},
						"conditions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The matching condition list of the cc protection rule.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The field of the condition.`,
									},
									"subfield": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The subfield of the condition.`,
									},
									"logic": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The condition matching logic.`,
									},
									"content": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `The content of the match condition.`,
									},
									"reference_table_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The reference table ID.`,
									},
								},
							},
						},
						"action": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The protective action taken when the number of requests reaches the upper limit.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protective_action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The protective action type.`,
									},
									"detail": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The block page detail information.`,
										Elem:        rulesActionDetailElem(),
									},
								},
							},
						},
						"rate_limit_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The rate limit mode.`,
						},
						"user_identifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The user identifier.`,
						},
						"other_user_identifier": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The other user identifier.`,
						},
						"limit_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of requests allowed from a web visitor in a rate limiting period.`,
						},
						"limit_period": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The rate limiting period.`,
						},
						"unlock_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The allowable frequency.`,
						},
						"lock_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The lock time for resuming normal page access after blocking can be set.`,
						},
						"request_aggregation": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether to enable domain aggregation statistics.`,
						},
						"all_waf_instances": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether to enable global counting.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the cc protection rule.`,
						},
					},
				},
			},
		},
	}
}

func rulesActionDetailElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"block_page_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the returned page.`,
			},
			"page_content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The content of the returned page.`,
			},
		},
	}
}

type RulesCcProtectionDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newRulesCcProtectionDSWrapper(d *schema.ResourceData, meta interface{}) *RulesCcProtectionDSWrapper {
	return &RulesCcProtectionDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceWafRulesCcProtectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newRulesCcProtectionDSWrapper(d, meta)
	listCcRulesRst, err := wrapper.ListCcRules()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listCcRulesToSchema(listCcRulesRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API WAF GET /v1/{project_id}/waf/policy/{policy_id}/cc
func (w *RulesCcProtectionDSWrapper) ListCcRules() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "waf")
	if err != nil {
		return nil, err
	}

	d := w.ResourceData
	uri := "/v1/{project_id}/waf/policy/{policy_id}/cc"
	uri = strings.ReplaceAll(uri, "{policy_id}", d.Get("policy_id").(string))
	params := map[string]any{
		"enterprise_project_id": w.Get("enterprise_project_id"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("items", "offset", "limit", 100).
		Filter(
			filters.New().From("items").
				Where("id", "=", w.Get("rule_id")).
				Where("name", "=", w.Get("name")).
				Where("status", "=", w.GetToInt("status")),
		).
		Request().
		Result()
}

func (w *RulesCcProtectionDSWrapper) listCcRulesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("rules", schemas.SliceToList(body.Get("items"),
			func(rule gjson.Result) any {
				return map[string]any{
					"id":          rule.Get("id").Value(),
					"policy_id":   rule.Get("policyid").Value(),
					"name":        rule.Get("name").Value(),
					"status":      rule.Get("status").Value(),
					"description": rule.Get("description").Value(),
					"conditions": schemas.SliceToList(rule.Get("conditions"),
						func(condition gjson.Result) any {
							return map[string]any{
								"field":              condition.Get("category").Value(),
								"subfield":           condition.Get("index").Value(),
								"logic":              condition.Get("logic_operation").Value(),
								"content":            schemas.SliceToStrList(condition.Get("contents")),
								"reference_table_id": condition.Get("value_list_id").Value(),
							}
						},
					),
					"action": schemas.SliceToList(rule.Get("action"),
						func(action gjson.Result) any {
							return map[string]any{
								"protective_action": action.Get("category").Value(),
								"detail":            w.setIteActDet(action),
							}
						},
					),
					"rate_limit_mode":       rule.Get("tag_type").Value(),
					"user_identifier":       rule.Get("tag_index").Value(),
					"other_user_identifier": schemas.SliceToStrList(rule.Get("tag_condition.contents")),
					"limit_num":             rule.Get("limit_num").Value(),
					"limit_period":          rule.Get("limit_period").Value(),
					"unlock_num":            rule.Get("unlock_num").Value(),
					"lock_time":             rule.Get("lock_time").Value(),
					"request_aggregation":   rule.Get("domain_aggregation").Value(),
					"all_waf_instances":     rule.Get("region_aggregation").Value(),
					"created_at":            w.setItemsTimestamp(rule),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*RulesCcProtectionDSWrapper) setIteActDet(action gjson.Result) any {
	return schemas.SliceToList(action.Get("detail"), func(detail gjson.Result) any {
		return map[string]any{
			"block_page_type": detail.Get("response.content_type").Value(),
			"page_content":    detail.Get("response.content").Value(),
		}
	})
}

func (*RulesCcProtectionDSWrapper) setItemsTimestamp(data gjson.Result) string {
	rawDate := data.Get("timestamp").Int()
	return utils.FormatTimeStampRFC3339(rawDate/1000, false)
}
