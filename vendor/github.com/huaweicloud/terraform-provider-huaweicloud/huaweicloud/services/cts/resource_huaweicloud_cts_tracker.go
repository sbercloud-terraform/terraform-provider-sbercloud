package cts

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	client "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3"
	cts "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceCTSTracker() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCTSTrackerCreate,
		ReadContext:   resourceCTSTrackerRead,
		UpdateContext: resourceCTSTrackerUpdate,
		DeleteContext: resourceCTSTrackerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCTSTrackerImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"bucket_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"bucket_name"},
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 64),
					validation.StringMatch(regexp.MustCompile(`^[\.\-_A-Za-z0-9]+$`),
						"only letters, numbers, hyphens (-), underscores (_), and periods (.) are allowed"),
				),
			},
			"validate_file": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"bucket_name"},
			},
			"kms_id": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"bucket_name"},
			},
			"lts_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"organization_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transfer_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceCTSTrackerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	ctsClient, err := cfg.HcCtsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	resourceID := "system"
	tracker, err := getSystemTracker(ctsClient)
	if err == nil && tracker != nil {
		log.Print("[DEBUG] the system tracker already exists, update the configuration")
		if tracker.Id != nil {
			resourceID = *tracker.Id
		}

		d.SetId(resourceID)
		return resourceCTSTrackerUpdate(ctx, d, meta)
	}

	// return the error with non-404 code
	if _, ok := err.(golangsdk.ErrDefault404); !ok {
		return diag.Errorf("error retrieving CTS tracker: %s", err)
	}

	resourceID, err = createSystemTracker(d, ctsClient)
	if err != nil {
		return diag.Errorf("error creating CTS tracker: %s", err)
	}

	d.SetId(resourceID)

	if rawTag := d.Get("tags").(map[string]interface{}); len(rawTag) > 0 {
		tagList := expandResourceTags(rawTag)
		_, err = ctsClient.BatchCreateResourceTags(buildCreateTagOpt(tagList, resourceID))
		if err != nil {
			return diag.Errorf("error creating CTS tracker tags: %s", err)
		}
	}

	// disable status if necessary
	if enabled := d.Get("enabled").(bool); !enabled {
		if err := updateSystemTrackerStatus(ctsClient, "disabled"); err != nil {
			return diag.Errorf("failed to disable CTS system tracker: %s", err)
		}
	}
	return nil
}

func resourceCTSTrackerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	ctsClient, err := cfg.HcCtsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	// update status firstly
	if d.IsNewResource() || d.HasChange("enabled") {
		status := "enabled"
		if enabled := d.Get("enabled").(bool); !enabled {
			status = "disabled"
		}

		if err := updateSystemTrackerStatus(ctsClient, status); err != nil {
			return diag.Errorf("error updating CTS tracker status: %s", err)
		}
	}

	// update other configurations
	if d.IsNewResource() || d.HasChangeExcept("enabled") {
		obsInfo := cts.TrackerObsInfo{
			BucketName:     utils.String(d.Get("bucket_name").(string)),
			FilePrefixName: utils.String(d.Get("file_prefix").(string)),
		}

		trackerType := cts.GetUpdateTrackerRequestBodyTrackerTypeEnum().SYSTEM
		updateBody := cts.UpdateTrackerRequestBody{
			TrackerName:           "system",
			TrackerType:           trackerType,
			IsLtsEnabled:          utils.Bool(d.Get("lts_enabled").(bool)),
			IsOrganizationTracker: utils.Bool(d.Get("organization_enabled").(bool)),
			IsSupportValidate:     utils.Bool(d.Get("validate_file").(bool)),
			ObsInfo:               &obsInfo,
		}

		var encryption bool
		if v, ok := d.GetOk("kms_id"); ok {
			encryption = true
			updateBody.KmsId = utils.String(v.(string))
		}
		updateBody.IsSupportTraceFilesEncryption = &encryption

		log.Printf("[DEBUG] updating CTS tracker options: %#v", updateBody)
		updateOpts := cts.UpdateTrackerRequest{
			Body: &updateBody,
		}

		_, err = ctsClient.UpdateTracker(&updateOpts)
		if err != nil {
			return diag.Errorf("error updating CTS tracker: %s", err)
		}

		if d.HasChange("tags") {
			err = updateResourceTags(ctsClient, d)
			if err != nil {
				return diag.Errorf("error updating CTS tracker tags: %s", err)
			}
		}
	}

	return resourceCTSTrackerRead(ctx, d, meta)
}

func resourceCTSTrackerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ctsClient, err := cfg.HcCtsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	ctsTracker, err := getSystemTracker(ctsClient)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CTS tracker")
	}

	if ctsTracker.Id != nil {
		d.SetId(*ctsTracker.Id)
	} else {
		d.SetId("system")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("name", ctsTracker.TrackerName),
		d.Set("lts_enabled", ctsTracker.Lts.IsLtsEnabled),
		d.Set("organization_enabled", ctsTracker.IsOrganizationTracker),
		d.Set("validate_file", ctsTracker.IsSupportValidate),
		d.Set("kms_id", ctsTracker.KmsId),
	)

	if ctsTracker.ObsInfo != nil {
		bucketName := ctsTracker.ObsInfo.BucketName
		mErr = multierror.Append(
			mErr,
			d.Set("bucket_name", bucketName),
			d.Set("file_prefix", ctsTracker.ObsInfo.FilePrefixName),
		)

		if *bucketName != "" {
			mErr = multierror.Append(mErr, d.Set("transfer_enabled", true))
		} else {
			mErr = multierror.Append(mErr, d.Set("transfer_enabled", false))
		}
	}

	if ctsTracker.TrackerType != nil {
		mErr = multierror.Append(mErr, d.Set("type", formatValue(ctsTracker.TrackerType)))
	}
	if ctsTracker.Status != nil {
		status := formatValue(ctsTracker.Status)
		mErr = multierror.Append(
			mErr,
			d.Set("status", status),
			d.Set("enabled", status == "enabled"),
		)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCTSTrackerDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	ctsClient, err := cfg.HcCtsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	if err := updateSystemTrackerStatus(ctsClient, "disabled"); err != nil {
		return diag.Errorf("failed to disable CTS system tracker: %s", err)
	}

	obsInfo := cts.TrackerObsInfo{
		BucketName:     utils.String(""),
		FilePrefixName: utils.String(""),
	}

	updateBody := cts.UpdateTrackerRequestBody{
		TrackerName:                   "system",
		TrackerType:                   cts.GetUpdateTrackerRequestBodyTrackerTypeEnum().SYSTEM,
		IsLtsEnabled:                  utils.Bool(false),
		IsSupportValidate:             utils.Bool(false),
		IsSupportTraceFilesEncryption: utils.Bool(false),
		KmsId:                         utils.String(""),
		ObsInfo:                       &obsInfo,
	}

	log.Printf("[DEBUG] updating CTS tracker to default configuration: %#v", updateBody)
	updateOpts := cts.UpdateTrackerRequest{
		Body: &updateBody,
	}

	_, err = ctsClient.UpdateTracker(&updateOpts)
	if err != nil {
		return diag.Errorf("falied to reset CTS system tracker: %s", err)
	}

	oldRaw, _ := d.GetChange("tags")
	if oldTags := oldRaw.(map[string]interface{}); len(oldTags) > 0 {
		oldTagList := expandResourceTags(oldTags)
		_, err = ctsClient.BatchDeleteResourceTags(buildDeleteTagOpt(oldTagList, d.Id()))
		if err != nil {
			return diag.Errorf("falied to delete CTS system tracker tags: %s", err)
		}
	}

	return nil
}

func resourceCTSTrackerImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	name := d.Id()
	d.Set("name", name)
	return []*schema.ResourceData{d}, nil
}

func formatValue(i interface{}) string {
	jsonRaw, err := json.Marshal(i)
	if err != nil {
		log.Printf("[WARN] failed to marshal %#v: %s", i, err)
		return ""
	}

	return strings.Trim(string(jsonRaw), `"`)
}

func createSystemTracker(d *schema.ResourceData, ctsClient *client.CtsClient) (string, error) {
	obsInfo := cts.TrackerObsInfo{
		BucketName:     utils.String(d.Get("bucket_name").(string)),
		FilePrefixName: utils.String(d.Get("file_prefix").(string)),
	}

	trackerType := cts.GetCreateTrackerRequestBodyTrackerTypeEnum().SYSTEM
	reqBody := cts.CreateTrackerRequestBody{
		TrackerName:           "system",
		TrackerType:           trackerType,
		IsLtsEnabled:          utils.Bool(d.Get("lts_enabled").(bool)),
		IsOrganizationTracker: utils.Bool(d.Get("organization_enabled").(bool)),
		IsSupportValidate:     utils.Bool(d.Get("validate_file").(bool)),
		ObsInfo:               &obsInfo,
	}

	if v, ok := d.GetOk("kms_id"); ok {
		encryption := true
		reqBody.KmsId = utils.String(v.(string))
		reqBody.IsSupportTraceFilesEncryption = &encryption
	}

	log.Printf("[DEBUG] creating system CTS tracker options: %#v", reqBody)
	createOpts := cts.CreateTrackerRequest{
		Body: &reqBody,
	}

	resp, err := ctsClient.CreateTracker(&createOpts)
	if err != nil {
		return "", err
	}
	if resp.Id == nil {
		return "", fmt.Errorf("ID is not found in API response")
	}

	return *resp.Id, nil
}

func getSystemTracker(ctsClient *client.CtsClient) (*cts.TrackerResponseBody, error) {
	name := "system"
	listOpts := &cts.ListTrackersRequest{
		TrackerName: &name,
	}

	response, err := ctsClient.ListTrackers(listOpts)
	if err != nil {
		return nil, err
	}

	if response.Trackers == nil || len(*response.Trackers) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	allTrackers := *response.Trackers
	return &allTrackers[0], nil
}

func updateSystemTrackerStatus(c *client.CtsClient, status string) error {
	enabledStatus := new(cts.UpdateTrackerRequestBodyStatus)
	if err := enabledStatus.UnmarshalJSON([]byte(status)); err != nil {
		return fmt.Errorf("failed to parse status %s: %s", status, err)
	}

	trackerType := cts.GetUpdateTrackerRequestBodyTrackerTypeEnum().SYSTEM
	statusOpts := cts.UpdateTrackerRequestBody{
		TrackerName: "system",
		TrackerType: trackerType,
		Status:      enabledStatus,
	}
	statusReq := cts.UpdateTrackerRequest{
		Body: &statusOpts,
	}

	_, err := c.UpdateTracker(&statusReq)
	return err
}
