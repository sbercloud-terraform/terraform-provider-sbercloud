package dcs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dcs/v2/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

type ParamsAttributes struct {
	ParamName      string `json:"param_name"`
	ParamValue     string `json:"param_value"`
	ValueType      string `json:"value_type"`
	NeedRestart    bool   `json:"need_restart"`
	UserPermission string `json:"user_permission"`
}
type ParamsConfig struct {
	Config []ParamsAttributes `json:"redis_config"`
}

func ResourceDcsParameters() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsParametersCreate,
		ReadContext:   resourceDcsParametersRead,
		DeleteContext: resourceDcsParametersDelete,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"configuration_parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"need_restart": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"user_permission": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceDcsParametersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DcsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DCS Client(v2): %s", err)
	}

	projectId := d.Get("project_id").(string)
	instanceId := d.Get("instance_id").(string)
	parameters := buildParameters(d)

	urlString := fmt.Sprintf("https://dcs.ru-moscow-1.hc.sbercloud.ru/v1.0/%s/instances/%s/configs", projectId, instanceId)

	var rst golangsdk.Result
	resp, err := client.Put(urlString, parameters, &rst, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.Header.Get("X-Request-Id"))
	return resourceDcsParametersRead(ctx, d, meta)
}

func resourceDcsParametersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DcsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DCS Client(v2): %s", err)
	}

	projectId := d.Get("project_id").(string)
	instanceId := d.Get("instance_id").(string)
	parameters := buildParameters(d)

	urlString := fmt.Sprintf("https://dcs.ru-moscow-1.hc.sbercloud.ru/v1.0/%s/instances/%s/configs", projectId, instanceId)

	resp, err := client.Get(urlString, nil, &golangsdk.RequestOpts{
		OkCodes:          []int{200},
		MoreHeaders:      instances.RequestOpts.MoreHeaders,
		KeepResponseBody: true,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	var respBody ParamsConfig
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("parameters", findParameters(parameters, respBody))
	if err != nil {
		return diag.Errorf("error setting attributes: %s", err)
	}
	attributes := buildAttributes(respBody)
	err = d.Set("configuration_parameters", attributes)
	if err != nil {
		return diag.Errorf("error setting attributes: %s", err)
	}
	return diags
}

func resourceDcsParametersDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return diag.Diagnostics{}
}

func buildParameters(d *schema.ResourceData) ParamsConfig {
	parameters := d.Get("parameters").(map[string]interface{})
	paramsConf := ParamsConfig{}
	for name, value := range parameters {
		attributes := ParamsAttributes{ParamName: name, ParamValue: value.(string)}
		paramsConf.Config = append(paramsConf.Config, attributes)
	}
	return paramsConf
}

func findParameters(reqParams, respParams ParamsConfig) map[string]interface{} {
	res := make(map[string]interface{}, len(reqParams.Config))

	for _, val := range reqParams.Config {
		res[val.ParamName] = ""
	}

	for _, val := range respParams.Config {
		if _, ok := res[val.ParamName]; ok {
			res[val.ParamName] = val.ParamValue
		}
	}
	return res
}

func buildAttributes(paramsConf ParamsConfig) []map[string]interface{} {

	parameters := make([]map[string]interface{}, len(paramsConf.Config))

	for i, val := range paramsConf.Config {
		attributes := make(map[string]interface{}, 5)
		attributes["name"] = val.ParamName
		attributes["value"] = val.ParamValue
		attributes["type"] = val.ValueType
		attributes["need_restart"] = val.NeedRestart
		attributes["user_permission"] = val.UserPermission
		parameters[i] = attributes
	}
	return parameters
}
