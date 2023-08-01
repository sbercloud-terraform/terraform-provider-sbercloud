package dcs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"io"
	"net/http"
	"time"
)

type body struct {
	Remark    string `json:"remark"`
	Backup_id string `json:"backup_id"`
}
type respBody struct {
	RestoreId string `json:"restore_id"`
}

type ReadRespBody struct {
	RestoreRecordResponse []struct {
		Status             string      `json:"status"`
		Progress           string      `json:"progress"`
		RestoreId          string      `json:"restore_id"`
		BackupId           string      `json:"backup_id"`
		RestoreRemark      string      `json:"restore_remark"`
		BackupRemark       interface{} `json:"backup_remark"`
		CreatedAt          string      `json:"created_at"`
		UpdatedAt          string      `json:"updated_at"`
		RestoreName        string      `json:"restore_name"`
		BackupName         string      `json:"backup_name"`
		SourceInstanceId   string      `json:"sourceInstanceId"`
		SourceInstanceName string      `json:"sourceInstanceName"`
		ErrorCode          interface{} `json:"error_code"`
	} `json:"restore_record_response"`
	TotalNum int `json:"total_num"`
}

func ResourceDcsRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsRestoreCreate,
		ReadContext:   resourceDcsRestoreRead,
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
				ForceNew: true,
				Default:  "restore instance",
				Optional: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"restore_records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status":               {Type: schema.TypeString, Computed: true},
						"progress":             {Type: schema.TypeString, Computed: true},
						"restore_id":           {Type: schema.TypeString, Computed: true},
						"backup_id":            {Type: schema.TypeString, Computed: true},
						"restore_remark":       {Type: schema.TypeString, Computed: true},
						"backup_remark":        {Type: schema.TypeString, Computed: true},
						"created_at":           {Type: schema.TypeString, Computed: true},
						"updated_at":           {Type: schema.TypeString, Computed: true},
						"restore_name":         {Type: schema.TypeString, Computed: true},
						"backup_name":          {Type: schema.TypeString, Computed: true},
						"source_instance_id":   {Type: schema.TypeString, Computed: true},
						"source_instance_name": {Type: schema.TypeString, Computed: true},
						"error_code":           {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func resourceDcsRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	cfg := meta.(*config.Config)

	region := cfg.GetRegion(d)
	client, err := cfg.DcsV2Client(region)

	project_id := d.Get("project_id").(string)
	instance_id := d.Get("instance_id").(string)
	remark := d.Get("remark").(string)
	backup_id := d.Get("backup_id").(string)

	url := fmt.Sprintf("https://dcs.ru-moscow-1.hc.sbercloud.ru/v2/%s/instances/%s/restores", project_id, instance_id)

	reqBody := body{
		Remark:    remark,
		Backup_id: backup_id,
	}

	opts := golangsdk.RequestOpts{OkCodes: []int{200}, JSONBody: reqBody, KeepResponseBody: true}
	resp, err := client.Request("POST", url, &opts)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Failed to create DCS restore. Status code: %d. Error: %d", resp.StatusCode, resp.Status)
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	responseBody := respBody{}

	err = json.Unmarshal(respBytes, &responseBody)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(responseBody.RestoreId)
	diags = resourceDcsRestoreRead(ctx, d, meta)
	return diags
}

func resourceDcsRestoreRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	region := cfg.GetRegion(d)
	client, err := cfg.DcsV2Client(region)
	if err != nil {
		return diag.FromErr(err)
	}
	project_id := d.Get("project_id").(string)
	instance_id := d.Get("instance_id").(string)

	url := fmt.Sprintf("https://dcs.ru-moscow-1.hc.sbercloud.ru/v2/%s/instances/%s/restores", project_id, instance_id)

	opts := golangsdk.RequestOpts{OkCodes: []int{200}, KeepResponseBody: true}
	resp, err := client.Request("GET", url, &opts)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Failed to get DCS restores list. Status code: %d. Error: %d", resp.StatusCode, resp.Status)
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	responseBody := ReadRespBody{}

	err = json.Unmarshal(respBytes, &responseBody)
	if err != nil {
		return diag.FromErr(err)
	}

	result := flattenRecordsList(responseBody)
	if result == nil {
		return diag.Errorf("records list is empty")
	}

	return diag.FromErr(d.Set("restore_records", result))
}

func flattenRecordsList(respBody ReadRespBody) []map[string]interface{} {
	if len(respBody.RestoreRecordResponse) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(respBody.RestoreRecordResponse))
	for i, record := range respBody.RestoreRecordResponse {
		result[i] = map[string]interface{}{
			"status":               record.Status,
			"progress":             record.Progress,
			"restore_id":           record.RestoreId,
			"backup_id":            record.BackupId,
			"restore_remark":       record.RestoreRemark,
			"backup_remark":        record.BackupRemark,
			"created_at":           record.CreatedAt,
			"updated_at":           record.UpdatedAt,
			"restore_name":         record.RestoreName,
			"backup_name":          record.BackupName,
			"source_instance_id":   record.SourceInstanceId,
			"source_instance_name": record.SourceInstanceName,
			"error_code":           record.ErrorCode,
		}
	}
	return result
}

func resourceDcsRestoreDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("Delete")
	return diag.Diagnostics{}
}
