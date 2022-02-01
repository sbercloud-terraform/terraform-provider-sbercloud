package dis

import (
	"context"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dis/v2/streams"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const disSysTagKeyEnterpriseProjectId = "_sys_enterprise_project_id"

func ResourceDisStream() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDisStreamCreate,
		ReadContext:   resourceDisStreamRead,
		DeleteContext: resourceDisStreamDelete,
		UpdateContext: resourceDisStreamUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"stream_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"partition_count": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"retention_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      24,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(24, 72),
			},

			"stream_type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"data_type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"data_schema": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"csv_delimiter": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"compression_format": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"auto_scale_min_partition_count": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"auto_scale_max_partition_count"},
			},

			"auto_scale_max_partition_count": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"auto_scale_min_partition_count"},
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"tags": common.TagsSchema(),

			"created": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"readable_partition_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"writable_partition_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"stream_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"partitions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hash_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sequence_number_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceDisStreamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DisV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DIS v2 client, err=%s", err)
	}

	opts := streams.CreateOpts{
		StreamName:        d.Get("stream_name").(string),
		PartitionCount:    d.Get("partition_count").(int),
		StreamType:        d.Get("stream_type").(string),
		DataDuration:      d.Get("retention_period").(int),
		DataType:          d.Get("data_type").(string),
		DataSchema:        d.Get("data_schema").(string),
		CompressionFormat: d.Get("compression_format").(string),
		Tags:              utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("csv_delimiter"); ok {
		opts.CsvProperties = &streams.CsvProperty{Delimiter: v.(string)}
	}

	// scale partitions
	autoScaleMinPartitionCount := d.Get("auto_scale_min_partition_count").(int)
	autoScaleMaxPartitionCount := d.Get("auto_scale_max_partition_count").(int)
	if autoScaleMinPartitionCount > 0 && autoScaleMaxPartitionCount > 0 {
		opts.AutoScaleEnabled = utils.Bool(true)
		opts.AutoScaleMinPartitionCount = &autoScaleMinPartitionCount
		opts.AutoScaleMaxPartitionCount = &autoScaleMaxPartitionCount
	} else {
		opts.AutoScaleEnabled = utils.Bool(false)
	}

	enterpriseProjectID := config.GetEnterpriseProjectID(d)
	if enterpriseProjectID != "" {
		opts.SysTags = []tags.ResourceTag{
			{
				Key:   disSysTagKeyEnterpriseProjectId,
				Value: enterpriseProjectID,
			},
		}
	}

	logp.Printf("[DEBUG] Creating new Cluster: %#v", opts)
	_, createErr := streams.Create(client, opts)
	if createErr != nil {
		return fmtp.DiagErrorf("Error creating DIS streams: %s", createErr)
	}

	d.SetId(opts.StreamName)

	return resourceDisStreamRead(ctx, d, meta)
}

func resourceDisStreamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DisV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DIS v2 client, err=%s", err)
	}

	detail, dErr := streams.Get(client, d.Id(), streams.GetOpts{})
	if dErr != nil {
		return fmtp.DiagErrorf("Error query DisStream %q:%s", d.Id(), dErr)
	}

	mErr := multierror.Append(
		d.Set("stream_name", detail.StreamName),
		d.Set("auto_scale_max_partition_count", detail.AutoScaleMaxPartitionCount),
		d.Set("auto_scale_min_partition_count", detail.AutoScaleMinPartitionCount),
		d.Set("compression_format", detail.CompressionFormat),
		d.Set("csv_delimiter", detail.CsvProperties.Delimiter),
		d.Set("data_schema", detail.DataSchema),
		d.Set("data_type", detail.DataType),
		d.Set("retention_period", detail.RetentionPeriod),
		d.Set("stream_type", detail.StreamType),
		d.Set("tags", utils.TagsToMap(detail.Tags)),
		d.Set("created", detail.CreateTime),
		d.Set("readable_partition_count", detail.ReadablePartitionCount),
		d.Set("writable_partition_count", detail.WritablePartitionCount),
		d.Set("partition_count", detail.WritablePartitionCount),
		d.Set("status", detail.Status),
		d.Set("stream_id", detail.StreamId),
		queryAndSetPartitionsToState(client, d, detail.StreamName),
	)

	enterpriseProjectId := parseEnterpriseProjectIdFromSysTags(detail.SysTags)
	if enterpriseProjectId != "" && enterpriseProjectId != "0" {
		multierror.Append(mErr, d.Set("enterprise_project_id", enterpriseProjectId))
	}

	if setSdErr := mErr.ErrorOrNil(); setSdErr != nil {
		return fmtp.DiagErrorf("Error setting vault fields: %s", setSdErr)
	}

	return nil
}

func resourceDisStreamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DisV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DIS v2 client, err=%s", err)
	}

	name := d.Id()
	errResult := streams.Delete(client, name)
	if errResult.Err != nil {
		return fmtp.DiagErrorf("Delete DIS streams failed.stream_name:%s,error:%s", name, errResult.Err)
	}

	d.SetId("")
	return nil
}

func resourceDisStreamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DisV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DIS v2 client, err=%s", err)
	}
	name := d.Id()

	// Update partition count
	if d.HasChange("partition_count") {
		newValue := d.Get("partition_count").(int)
		updateOpts := streams.UpdatePartitionOpt{
			StreamName:           name,
			TargetPartitionCount: newValue,
		}
		_, extendErr := streams.UpdatePartition(client, name, updateOpts)
		if extendErr != nil {
			return fmtp.DiagErrorf("Update DIS stream failed.stream_name=%s,error=%s", name, extendErr)
		}

		checkErr := checkPartitionUpdateResult(ctx, client, name, newValue, d.Timeout(schema.TimeoutUpdate))
		if checkErr != nil {
			return fmtp.DiagErrorf("Update DIS stream failed.stream_name=%s,error=%s", name, checkErr)
		}
	}

	if d.HasChange("tags") {
		streamId := d.Get("stream_id").(string)
		tagErr := utils.UpdateResourceTags(client, d, "stream", streamId)
		if tagErr != nil {
			return fmtp.DiagErrorf("Error updating tags of DIS stream:%s,streamId=%s, err:%s", name, streamId, tagErr)
		}
	}

	return resourceDisStreamRead(ctx, d, meta)
}

func parseEnterpriseProjectIdFromSysTags(value []tags.ResourceTag) (r string) {
	if len(value) == 0 {
		return
	}

	for i := 0; i < len(value); i++ {
		item := value[i]
		if item.Key == disSysTagKeyEnterpriseProjectId {
			return item.Value
		}
	}
	return
}

func checkPartitionUpdateResult(ctx context.Context, client *golangsdk.ServiceClient, name string, targetValue int,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Done"},
		Refresh: func() (interface{}, string, error) {
			resp, err := streams.Get(client, name, streams.GetOpts{})
			if err != nil {
				return nil, "failed", err
			}
			logp.Printf("WritablePartitionCount=", resp.WritablePartitionCount, targetValue)
			if resp.WritablePartitionCount == targetValue {
				return resp, "Done", nil
			}
			return resp, "Pending", nil
		},
		Timeout:      timeout,
		PollInterval: 5 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("waiting for DIS stream (%s) to update partition failed: %s", name, err)
	}
	return nil
}

func queryAndSetPartitionsToState(client *golangsdk.ServiceClient, d *schema.ResourceData, streamName string) error {
	var result []map[string]interface{}
	opts := streams.GetOpts{}
	for {
		rst, gErr := streams.Get(client, streamName, opts)
		if gErr != nil {
			return fmtp.Errorf("Error query the partitions of DIS stream, err=%s", gErr)
		}

		for _, partition := range rst.Partitions {
			result = append(result, map[string]interface{}{
				"id":                    partition.PartitionId,
				"status":                partition.Status,
				"hash_range":            partition.HashRange,
				"sequence_number_range": partition.SequenceNumberRange,
			})
		}

		if !rst.HasMorePartitions {
			break
		}

		opts.StartPartitionId = rst.Partitions[len(rst.Partitions)-1].PartitionId
	}

	return d.Set("partitions", result)
}
