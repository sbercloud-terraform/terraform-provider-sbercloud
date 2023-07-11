package dcs2

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/services/internal/core"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	access_key = ""
	secret_key = ""
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
			"param_id": {
				Type:     schema.TypeString,
				Required: true,
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
	}
}

func resourceDcsConfigParamsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	//cfg := meta.(*config.Config)

	var diags diag.Diagnostics

	s := core.Signer{Key: access_key, Secret: secret_key}

	project_id := d.Get("project_id").(string)
	instance_id := d.Get("instance_id").(string)
	param_id := d.Get("param_id").(string)
	param_name := d.Get("param_name").(string)
	param_value := d.Get("param_value").(string)
	urlString := fmt.Sprintf("https://dcs.ru-moscow-1.hc.sbercloud.ru/v1.0/%s/instances/%s/configs", project_id, instance_id)

	params := DCSParams{
		RedisConfig: []struct {
			ParamId    string `json:"param_id"`
			ParamName  string `json:"param_name"`
			ParamValue string `json:"param_value"`
		}{
			{ParamId: param_id, ParamName: param_name, ParamValue: param_value}}}

	order_bytes, err := json.Marshal(params)
	if err != nil {
		diag.FromErr(err)
	}
	r, err := http.NewRequest("PUT", urlString, strings.NewReader(string(order_bytes)))
	if err != nil {
		diag.FromErr(err)
	}

	errsign := s.Sign(r)
	if errsign != nil {
		diag.FromErr(err)

		log.Fatal(errsign)
	}
	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(resp.StatusCode))
	resourceDcsConfigParamsRead(ctx, d, meta)
	return diags
}

func resourceDcsConfigParamsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceDcsConfigParamsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}
