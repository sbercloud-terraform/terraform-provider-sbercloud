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
	"io"
)

type ParamsAttributes struct {
	ParamId    string `json:"param_id"`
	ParamName  string `json:"param_name"`
	ParamValue string `json:"param_value"`
}
type ParamsConfig struct {
	Config []ParamsAttributes `json:"redis_config"`
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
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

var (
	paramsID = map[string]string{
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
		"active-expire-num":                             "35",
	}
)

func resourceDcsConfigParamsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DcsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DCS Client(v2): %s", err)
	}

	projectId := d.Get("project_id").(string)
	instanceId := d.Get("instance_id").(string)
	parameters := d.Get("parameters").(map[string]interface{})
	paramsConf := ParamsConfig{}
	for name, value := range parameters {
		id, ok := paramsID[name]
		if !ok {
			return diag.Errorf("No such parameter %s", name)
		}
		attributes := ParamsAttributes{ParamId: id, ParamName: name, ParamValue: value.(string)}
		paramsConf.Config = append(paramsConf.Config, attributes)
	}

	urlString := fmt.Sprintf("https://dcs.ru-moscow-1.hc.sbercloud.ru/v1.0/%s/instances/%s/configs", projectId, instanceId)

	var rst golangsdk.Result
	resp, err := client.Put(urlString, paramsConf, &rst, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.Header.Get("X-Request-Id"))
	return resourceDcsConfigParamsRead(ctx, d, meta)
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
	parameters := d.Get("parameters").(map[string]interface{})
	paramsConf := ParamsConfig{}
	for name, value := range parameters {
		id, ok := paramsID[name]
		if !ok {
			return diag.Errorf("No such parameter %s", name)
		}
		attributes := ParamsAttributes{ParamId: id, ParamName: name, ParamValue: value.(string)}
		paramsConf.Config = append(paramsConf.Config, attributes)
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

	var respBody ParamsConfig
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		return diag.Errorf("unmarshall problem: %s", err)
	}
	err = d.Set("parameters", findParams(paramsConf, respBody))
	if err != nil {
		return diag.Errorf("error setting attributes: %s", err)
	}
	return diags
}

func resourceDcsConfigParamsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return diag.Diagnostics{}
}

func findParams(reqParams, respParams ParamsConfig) map[string]interface{} {
	reqParamNames := make([]string, len(respParams.Config))
	for i := 0; i < len(reqParams.Config); i++ {
		reqParamNames[i] = reqParams.Config[i].ParamName
	}
	res := make(map[string]interface{})
	for _, val := range respParams.Config {
		for _, paramName := range reqParamNames {
			if paramName == val.ParamName {
				res[paramName] = val.ParamValue
			}
		}
	}
	return res
}
