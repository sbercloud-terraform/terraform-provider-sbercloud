package dcs2

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dcs/v2/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/services/internal/core"
	"io"
	"net/http"
	"strings"
)

type DCSParams struct {
	RedisConfig []struct {
		ParamId    string `json:"param_id"`
		ParamName  string `json:"param_name"`
		ParamValue string `json:"param_value"`
	} `json:"redis_config"`
}

func ResourceDcsConfigParams() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsConfigParamsCreate,
		ReadContext:   resourceDcsConfigParamsRead,
		DeleteContext: resourceDcsConfigParamsDelete,
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
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"param_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"param_value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}
func resourceDcsConfigParamsRead3(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	return diags
}

var (
	ParamsMap = map[string]string{
		"timeout":                                       "1",
		"maxmemory-policy":                              "2",
		"hash-max-ziplist-entries":                      "3",
		"hash-max-ziplist-value":                        "4",
		"list-max-ziplist-entries":                      "5",
		"list-max-ziplist-value":                        "6",
		"set-max-intset-entries":                        "7",
		"zset-max-ziplist-entries":                      "8",
		"zset-max-ziplist-value":                        "9",
		"latency-monitor-threshold":                     "10",
		"maxclients":                                    "11",
		"reserved-memory":                               "12",
		"notify-keyspace-events":                        "13",
		"repl-backlog-size":                             "14",
		"repl-backlog-ttl":                              "15",
		"appendfsync":                                   "16",
		"appendonly":                                    "17",
		"slowlog-log-slower-than":                       "18",
		"slowlog-max-len":                               "19",
		"lua-time-limit":                                "20",
		"repl-timeout":                                  "21",
		"proto-max-bulk-len":                            "22",
		"master-read-only":                              "23",
		"client-output-buffer-slave-soft-limit":         "24",
		"client-output-buffer-slave-hard-limit":         "25",
		"client-output-buffer-limit-slave-soft-seconds": "26",
		//:"27",
		//:"28",
		//:"29",
		//:"30",
		//:"31",
		//:"32",
		//:"33",
		//:"34",
		"active-expire-num": "35",
	}
)

func resourceDcsConfigParamsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DcsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DCS Client(v2): %s", err)
	}

	projectId := d.Get("project_id").(string)
	instanceId := d.Get("instance_id").(string)
	parameters := d.Get("parameters").([]interface{})
	params := DCSParams{}
	for _, val := range parameters {
		par := val.(map[string]interface{})
		paramName := par["param_name"].(string)
		paramValue := par["param_value"].(string)
		paramID, ok := ParamsMap[paramName]
		if !ok {
			return diag.Errorf("No such parameter %s", paramName)
		}

		params.RedisConfig = append(params.RedisConfig,
			struct {
				ParamId    string `json:"param_id"`
				ParamName  string `json:"param_name"`
				ParamValue string `json:"param_value"`
			}{ParamId: paramID, ParamName: paramName, ParamValue: paramValue})
	}

	urlString := fmt.Sprintf("https://dcs.ru-moscow-1.hc.sbercloud.ru/v1.0/%s/instances/%s/configs", projectId, instanceId)

	var rst golangsdk.Result
	resp, err := client.Put(urlString, params, &rst, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.Header.Get("X-Request-Id"))

	return diags
}

func resourceDcsConfigParamsCreate2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	cfg := meta.(*config.Config)

	var diags diag.Diagnostics

	s := core.Signer{Key: cfg.AccessKey, Secret: cfg.SecretKey}

	projectId := d.Get("project_id").(string)
	instanceId := d.Get("instance_id").(string)
	parameters := d.Get("parameters").([]interface{})
	params := DCSParams{}
	// попробовать парсить через json.unmarshal()
	for _, val := range parameters {
		par := val.(map[string]interface{})
		paramName := par["param_name"].(string)
		paramValue := par["param_value"].(string)
		paramID, ok := ParamsMap[paramName]
		if !ok {
			return diag.Errorf("No such parameter %s", paramName)
		}

		params.RedisConfig = append(params.RedisConfig,
			struct {
				ParamId    string `json:"param_id"`
				ParamName  string `json:"param_name"`
				ParamValue string `json:"param_value"`
			}{ParamId: paramID, ParamName: paramName, ParamValue: paramValue})
	}

	urlString := fmt.Sprintf("https://dcs.ru-moscow-1.hc.sbercloud.ru/v1.0/%s/instances/%s/configs", projectId, instanceId)

	reqBody, err := json.Marshal(params)
	if err != nil {
		return diag.FromErr(err)
	}
	r, err := http.NewRequest("PUT", urlString, strings.NewReader(string(reqBody)))
	if err != nil {
		return diag.FromErr(err)
	}

	err = s.Sign(r)
	if err != nil {
		return diag.FromErr(err)
	}
	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Header.Get("X-Request-Id"))
	resourceDcsConfigParamsRead(ctx, d, meta)
	return diags
}

func resourceDcsConfigParamsRead2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	var diags diag.Diagnostics

	s := core.Signer{Key: cfg.AccessKey, Secret: cfg.SecretKey}

	projectId := d.Get("project_id").(string)
	instanceId := d.Get("instance_id").(string)
	parameters := d.Get("parameters").([]interface{})
	params := DCSParams{}
	for _, val := range parameters {
		par := val.(map[string]interface{})
		paramName := par["param_name"].(string)
		paramValue := par["param_value"].(string)
		paramID, ok := ParamsMap[paramName]
		if !ok {
			return diag.Errorf("No such parameter %s", paramName)
		}

		params.RedisConfig = append(params.RedisConfig,
			struct {
				ParamId    string `json:"param_id"`
				ParamName  string `json:"param_name"`
				ParamValue string `json:"param_value"`
			}{ParamId: paramID, ParamName: paramName, ParamValue: paramValue})
	}

	urlString := fmt.Sprintf("https://dcs.ru-moscow-1.hc.sbercloud.ru/v1.0/%s/instances/%s/configs", projectId, instanceId)

	_, err := json.Marshal(params)
	if err != nil {
		return diag.FromErr(err)
	}
	r, err := http.NewRequest("GET", urlString, strings.NewReader(string("")))
	if err != nil {
		return diag.FromErr(err)
	}

	err = s.Sign(r)
	if err != nil {
		return diag.FromErr(err)
	}
	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		return diag.FromErr(err)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}
	var respBody DCSParams
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		return diag.Errorf("unmarshall problem: %s", err)
	}
	err = d.Set("parameters", findParams(params, respBody))
	if err != nil {
		return diag.Errorf("error setting attributes: %s", err)
	}

	return diags

}

func resourceDcsConfigParamsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.DcsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DCS Client(v2): %s", err)
	}

	projectId := d.Get("project_id").(string)
	instanceId := d.Get("instance_id").(string)
	parameters := d.Get("parameters").([]interface{})
	params := DCSParams{}
	for _, val := range parameters {
		par := val.(map[string]interface{})
		paramName := par["param_name"].(string)
		paramValue := par["param_value"].(string)
		paramID, ok := ParamsMap[paramName]
		if !ok {
			return diag.Errorf("No such parameter %s", paramName)
		}

		params.RedisConfig = append(params.RedisConfig,
			struct {
				ParamId    string `json:"param_id"`
				ParamName  string `json:"param_name"`
				ParamValue string `json:"param_value"`
			}{ParamId: paramID, ParamName: paramName, ParamValue: paramValue})
	}

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
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var respBody DCSParams
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		return diag.Errorf("unmarshall problem: %s", err)
	}
	err = d.Set("parameters", findParams(params, respBody))
	if err != nil {
		return diag.Errorf("error setting attributes: %s", err)
	}
	return diags
}

func findParams(reqParams, respParams DCSParams) []interface{} {
	reqParamNames := make([]string, 0)
	for _, val := range reqParams.RedisConfig {
		reqParamNames = append(reqParamNames, val.ParamName)
	}
	res := make([]interface{}, 0)
	for _, val := range respParams.RedisConfig {
		for _, paramName := range reqParamNames {
			if paramName == val.ParamName {
				parameters := make(map[string]interface{})
				parameters["param_id"] = ParamsMap[val.ParamName]
				parameters["param_value"] = val.ParamValue
				parameters["param_name"] = val.ParamName
				res = append(res, parameters)
			}
		}
	}
	return res
}

func resourceDcsConfigParamsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return diag.Diagnostics{}
}
