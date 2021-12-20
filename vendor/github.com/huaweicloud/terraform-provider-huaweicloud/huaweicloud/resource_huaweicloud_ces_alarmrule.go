package huaweicloud

import (
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/golangsdk/openstack/cloudeyeservice/alarmrule"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const nameCESAR = "CES-AlarmRule"

var cesAlarmActions = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	ForceNew: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"notification", "autoscaling",
				}, false),
			},

			"notification_list": {
				Type:     schema.TypeList,
				MaxItems: 5,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	},
}

func ResourceAlarmRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlarmRuleCreate,
		Read:   resourceAlarmRuleRead,
		Update: resourceAlarmRuleUpdate,
		Delete: resourceAlarmRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

			"alarm_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alarm_description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"metric": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:     schema.TypeString,
							Required: true,
						},

						"metric_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"dimensions": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 3,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"condition": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntInSlice([]int{0, 1, 300, 1200, 3600, 14400, 86400}),
						},

						"filter": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"max", "min", "average", "sum", "variance",
							}, false),
						},

						"comparison_operator": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								">=", ">", "<=", "<", "=",
							}, false),
						},

						"value": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 5),
						},

						"unit": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"alarm_actions": &cesAlarmActions,
			"ok_actions":    &cesAlarmActions,

			"alarm_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"alarm_level": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(1, 4),
			},

			"alarm_action_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"alarm_state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// deprecated
			"insufficientdata_actions": {
				Type:       schema.TypeList,
				Optional:   true,
				Deprecated: "insufficientdata_actions is deprecated",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"notification_list": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 5,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func getMetricOpts(d *schema.ResourceData) (alarmrule.MetricOpts, error) {
	mos, ok := d.Get("metric").([]interface{})
	if !ok {
		return alarmrule.MetricOpts{}, fmtp.Errorf("Error converting opt of metric:%v", d.Get("metric"))
	}
	mo := mos[0].(map[string]interface{})

	mod := mo["dimensions"].([]interface{})
	dopts := make([]alarmrule.DimensionOpts, len(mod))
	for i, v := range mod {
		v1 := v.(map[string]interface{})
		dopts[i] = alarmrule.DimensionOpts{
			Name:  v1["name"].(string),
			Value: v1["value"].(string),
		}
	}
	return alarmrule.MetricOpts{
		Namespace:  mo["namespace"].(string),
		MetricName: mo["metric_name"].(string),
		Dimensions: dopts,
	}, nil
}

func getAlarmAction(d *schema.ResourceData, name string) []alarmrule.ActionOpts {
	aos := d.Get(name).([]interface{})
	if len(aos) == 0 {
		return nil
	}
	opts := make([]alarmrule.ActionOpts, len(aos))
	for i, v := range aos {
		v1 := v.(map[string]interface{})

		v2 := v1["notification_list"].([]interface{})
		nl := make([]string, len(v2))
		for j, v3 := range v2 {
			nl[j] = v3.(string)
		}

		opts[i] = alarmrule.ActionOpts{
			Type:             v1["type"].(string),
			NotificationList: nl,
		}
	}
	return opts
}

func getAlarmCondition(d *schema.ResourceData) alarmrule.ConditionOpts {
	var opts alarmrule.ConditionOpts

	rawCondition := d.Get("condition").([]interface{})
	if len(rawCondition) == 1 {
		condition := rawCondition[0].(map[string]interface{})

		opts.Period = condition["period"].(int)
		opts.Filter = condition["filter"].(string)
		opts.ComparisonOperator = condition["comparison_operator"].(string)
		opts.Value = condition["value"].(int)
		opts.Unit = condition["unit"].(string)
		opts.Count = condition["count"].(int)
	}

	return opts
}

func resourceAlarmRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CesV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Cloud Eye Service client: %s", err)
	}

	metric, err := getMetricOpts(d)
	if err != nil {
		return err
	}

	createOpts := alarmrule.CreateOpts{
		AlarmName:               d.Get("alarm_name").(string),
		AlarmDescription:        d.Get("alarm_description").(string),
		AlarmLevel:              d.Get("alarm_level").(int),
		Metric:                  metric,
		Condition:               getAlarmCondition(d),
		AlarmActions:            getAlarmAction(d, "alarm_actions"),
		OkActions:               getAlarmAction(d, "ok_actions"),
		InsufficientdataActions: getAlarmAction(d, "insufficientdata_actions"),
		AlarmEnabled:            d.Get("alarm_enabled").(bool),
		AlarmActionEnabled:      d.Get("alarm_action_enabled").(bool),
		EnterpriseProjectID:     GetEnterpriseProjectID(d, config),
	}
	logp.Printf("[DEBUG] Create %s Options: %#v", nameCESAR, createOpts)

	r, err := alarmrule.Create(client, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating %s: %s", nameCESAR, err)
	}
	logp.Printf("[DEBUG] Create %s: %#v", nameCESAR, *r)

	d.SetId(r.AlarmID)

	return resourceAlarmRuleRead(d, meta)
}

func resourceAlarmRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CesV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Cloud Eye Service client: %s", err)
	}

	r, err := alarmrule.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "alarmrule")
	}
	logp.Printf("[DEBUG] Retrieved %s %s: %#v", nameCESAR, d.Id(), r)

	m, err := utils.ConvertStructToMap(r, map[string]string{"notificationList": "notification_list"})
	if err != nil {
		return err
	}

	alarmMetric := make([]interface{}, 1)
	alarmMetric[0] = m["metric"]
	alarmCondition := make([]interface{}, 1)
	alarmCondition[0] = m["condition"]

	mErr := multierror.Append(nil,
		d.Set("alarm_name", m["alarm_name"]),
		d.Set("alarm_description", m["alarm_description"]),
		d.Set("alarm_level", m["alarm_level"]),
		d.Set("metric", alarmMetric),
		d.Set("condition", alarmCondition),
		d.Set("alarm_actions", m["alarm_actions"]),
		d.Set("ok_actions", m["ok_actions"]),
		d.Set("alarm_enabled", m["alarm_enabled"]),
		d.Set("alarm_action_enabled", m["alarm_action_enabled"]),
		d.Set("alarm_state", m["alarm_state"]),
		d.Set("update_time", m["update_time"]),
		d.Set("enterprise_project_id", m["enterprise_project_id"]),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}

	return nil
}

func resourceAlarmRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CesV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Cloud Eye Service client: %s", err)
	}

	arId := d.Id()

	if d.HasChange("alarm_enabled") {
		enabled := d.Get("alarm_enabled").(bool)
		enableOpts := alarmrule.EnableOpts{
			AlarmEnabled: enabled,
		}
		logp.Printf("[DEBUG] Updating %s %s to %#v", nameCESAR, arId, enabled)

		timeout := d.Timeout(schema.TimeoutUpdate)
		//lintignore:R006
		err = resource.Retry(timeout, func() *resource.RetryError {
			err := alarmrule.Enable(client, arId, enableOpts).ExtractErr()
			if err != nil {
				return checkForRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmtp.Errorf("Error updating %s %s: %s", nameCESAR, arId, err)
		}
	}

	updateOpts := alarmrule.UpdateOpts{}
	changed := false
	if d.HasChanges("alarm_name", "alarm_description", "alarm_level", "alarm_action_enabled") {
		description := d.Get("alarm_description").(string)
		actionEnabled := d.Get("alarm_action_enabled").(bool)

		updateOpts.Name = d.Get("alarm_name").(string)
		updateOpts.AlarmLevel = d.Get("alarm_level").(int)
		updateOpts.Description = &description
		updateOpts.ActionEnabled = &actionEnabled
		changed = true
	}

	if d.HasChange("condition") {
		condition := getAlarmCondition(d)
		// unit field is not supported in Update
		condition.Unit = ""
		updateOpts.Condition = &condition
		changed = true
	}

	if changed {
		logp.Printf("[DEBUG] Updating %s %s opts: %#v", nameCESAR, arId, updateOpts)
		err := alarmrule.Update(client, arId, updateOpts).ExtractErr()
		if err != nil {
			return fmtp.Errorf("Error updating %s %s: %s", nameCESAR, arId, err)
		}
	}

	return resourceAlarmRuleRead(d, meta)
}

func resourceAlarmRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CesV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Cloud Eye Service client: %s", err)
	}

	arId := d.Id()
	logp.Printf("[DEBUG] Deleting %s %s", nameCESAR, arId)

	timeout := d.Timeout(schema.TimeoutDelete)
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		err := alarmrule.Delete(client, arId).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if utils.IsResourceNotFound(err) {
			logp.Printf("[INFO] deleting an unavailable %s: %s", nameCESAR, arId)
			return nil
		}
		return fmtp.Errorf("Error deleting %s %s: %s", nameCESAR, arId, err)
	}

	return nil
}
