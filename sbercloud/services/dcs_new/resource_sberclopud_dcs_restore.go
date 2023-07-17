package dcs_new

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
	"time"
)

const (
	access_key = ""
	secret_key = ""
)

type body struct {
	remark    string `json:"name"`
	backup_id string `json:"password"`
}

func ResourceDcsRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsRestoreCreate,
		ReadContext:   resourceDcsRestoreRead,
		//UpdateContext: //,
		DeleteContext: resourceDcsRestoreDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
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
			"remark": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Default:  "restore instance",
				Optional: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDcsRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	s := core.Signer{Key: access_key, Secret: secret_key}

	project_id := d.Get("project_id").(string)
	instance_id := d.Get("instance_id").(string)
	remark := d.Get("remark").(string)
	backup_id := d.Get("backup_id").(string)

	url := fmt.Sprintf("https://dcs.ru-moscow-1.hc.sbercloud.ru/v2/%s/instances/%s/restores", project_id, instance_id)

	reqBody := body{
		remark:    remark,
		backup_id: backup_id,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(reqBytes)))
	if err != nil {
		return diag.FromErr(err)
	}

	err = s.Sign(req)
	if err != nil {
		return diag.FromErr(err)
		log.Fatal(err)
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(resp.StatusCode))
	resourceDcsRestoreRead(ctx, d, meta)

	return diags
}

func resourceDcsRestoreRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceDcsRestoreDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}
