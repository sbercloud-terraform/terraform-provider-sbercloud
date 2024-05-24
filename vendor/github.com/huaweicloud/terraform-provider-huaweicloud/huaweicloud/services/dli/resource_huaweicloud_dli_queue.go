package dli

import (
	"context"
	"fmt"
	"log"
	"math"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dli/v1/queues"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var regexp4Name = regexp.MustCompile(`^[a-z0-9_]{1,128}$`)

const (
	CU16                  = 16
	CU64                  = 64
	CU256                 = 256
	resourceModeShared    = 0
	resourceModeExclusive = 1

	QueueTypeSQL         = "sql"
	QueueTypeGeneral     = "general"
	queueFeatureBasic    = "basic"
	queueFeatureAI       = "ai"
	queuePlatformX86     = "x86_64"
	queuePlatformAARCH64 = "aarch64"

	actionRestart  = "restart"
	actionScaleOut = "scale_out"
	actionScaleIn  = "scale_in"
)

func ResourceDliQueue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDliQueueCreate,
		ReadContext:   resourceDliQueueRead,
		UpdateContext: resourceDliQueueUpdate,
		DeleteContext: resourceDliQueueDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceQueueImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(regexp4Name,
					"only contain digits, lower letters, and underscores (_)"),
			},

			"queue_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      QueueTypeSQL,
				ValidateFunc: validation.StringInSlice([]string{QueueTypeSQL, QueueTypeGeneral}, false),
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"cu_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validCuCount,
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      queuePlatformX86,
				ValidateFunc: validation.StringInSlice([]string{queuePlatformX86, queuePlatformAARCH64}, false),
			},

			"resource_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{resourceModeShared, resourceModeExclusive}),
			},

			"feature": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{queueFeatureBasic, queueFeatureAI}, false),
			},

			"tags": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},

			"vpc_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"management_subnet_cidr": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "management_subnet_cidr is Deprecated",
			},

			"subnet_cidr": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "subnet_cidr is Deprecated",
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceDliQueueCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dliClient, err := cfg.DliV1Client(region)
	if err != nil {
		return diag.Errorf("creating dli client failed: %s", err)
	}

	queueName := d.Get("name").(string)

	log.Printf("[DEBUG] create dli queues queueName: %s", queueName)
	createOpts := queues.CreateOpts{
		QueueName:           queueName,
		QueueType:           d.Get("queue_type").(string),
		Description:         d.Get("description").(string),
		CuCount:             d.Get("cu_count").(int),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		Platform:            d.Get("platform").(string),
		ResourceMode:        d.Get("resource_mode").(int),
		Feature:             d.Get("feature").(string),
		Tags:                assembleTagsFromRecource("tags", d),
	}

	log.Printf("[DEBUG] create dli queues using parameters: %+v", createOpts)
	createResult := queues.Create(dliClient, createOpts)
	if createResult.Err != nil {
		return diag.Errorf("create dli queues failed: %s", createResult.Err)
	}

	// The resource ID (queue name) at this time is only used as a mark the resource, and the value will be refreshed
	// in the READ method.
	d.SetId(queueName)

	// This is a workaround to avoid issue: the queue is assigning, which is not available
	time.Sleep(4 * time.Minute) // lintignore:R018

	if v, ok := d.GetOk("vpc_cidr"); ok {
		err = updateVpcCidrOfQueue(dliClient, queueName, v.(string))
		if err != nil {
			return diag.Errorf("update cidr failed when creating dli queues: %s", err)
		}
	}

	return resourceDliQueueRead(ctx, d, meta)
}

func assembleTagsFromRecource(key string, d *schema.ResourceData) []tags.ResourceTag {
	if v, ok := d.GetOk(key); ok {
		tagRaw := v.(map[string]interface{})
		taglist := utils.ExpandResourceTags(tagRaw)
		return taglist
	}
	return nil
}

func resourceDliQueueRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DliV1Client, err=%s", err)
	}

	queueName := d.Get("name").(string)

	queryOpts := queues.ListOpts{
		QueueType: d.Get("queue_type").(string),
	}

	log.Printf("[DEBUG] query dli queues using parameters: %+v", queryOpts)

	queryAllResult := queues.List(client, queryOpts)
	if queryAllResult.Err != nil {
		return diag.Errorf("query queues failed: %s", queryAllResult.Err)
	}

	// filter by queue_name
	queueDetail, err := filterByQueueName(queryAllResult.Body, queueName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DLI queue")
	}
	d.SetId(queueDetail.ResourceId)

	log.Printf("[DEBUG]The detail of queue from SDK:%+v", queueDetail)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", queueDetail.QueueName),
		d.Set("queue_type", queueDetail.QueueType),
		d.Set("description", queueDetail.Description),
		d.Set("cu_count", queueDetail.CuCount),
		d.Set("enterprise_project_id", utils.StringIgnoreEmpty(queueDetail.EnterpriseProjectId)),
		d.Set("platform", queueDetail.Platform),
		d.Set("resource_mode", queueDetail.ResourceMode),
		d.Set("feature", queueDetail.Feature),
		d.Set("create_time", queueDetail.CreateTime),
		d.Set("vpc_cidr", queueDetail.CidrInVpc),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterByQueueName(body interface{}, queueName string) (r *queues.Queue, err error) {
	if queueList, ok := body.(*queues.ListResult); ok {
		log.Printf("[DEBUG]The list of queue from SDK:%+v", queueList)

		for _, v := range queueList.Queues {
			if v.QueueName == queueName {
				return &v, nil
			}
		}
		return nil, golangsdk.ErrDefault404{}
	}

	return nil, fmt.Errorf("sdk-client response type is wrong, expect type:*queues.ListResult,acutal Type:%T",
		body)
}

func resourceDliQueueDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DliV1Client, err=%s", err)
	}

	queueName := d.Get("name").(string)
	log.Printf("[DEBUG] Deleting dli Queue %q", queueName)

	result := queues.Delete(client, queueName)
	if result.Err != nil {
		return diag.Errorf("error deleting dli Queue %q, err=%s", queueName, result.Err)
	}

	return nil
}

// support cu_count scaling
func resourceDliQueueUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DliV1Client: %s", err)
	}

	queueName := d.Get("name").(string)
	opt := queues.ActionOpts{
		QueueName: queueName,
	}
	if d.HasChange("cu_count") {
		oldValue, newValue := d.GetChange("cu_count")
		cuChange := newValue.(int) - oldValue.(int)

		opt.CuCount = int(math.Abs(float64(cuChange)))
		opt.Action = buildScaleActionParam(oldValue.(int), newValue.(int))

		log.Printf("[DEBUG]DLI queue Update Option: %#v", opt)
		result := queues.ScaleOrRestart(client, opt)
		if result.Err != nil {
			return diag.Errorf("update dli queues failed, queueName=%s, error:%s", queueName, result.Err)
		}

		updateStateConf := &resource.StateChangeConf{
			Pending: []string{fmt.Sprintf("%d", oldValue)},
			Target:  []string{fmt.Sprintf("%d", newValue)},
			Refresh: func() (interface{}, string, error) {
				getResult := queues.Get(client, queueName)
				queueDetail := getResult.Body.(*queues.Queue4Get)
				return getResult, fmt.Sprintf("%d", queueDetail.CuCount), nil
			},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        30 * time.Second,
			PollInterval: 20 * time.Second,
		}
		_, err = updateStateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for dli.queue (%s) to be scale: %s", queueName, err)
		}
	}

	if d.HasChange("vpc_cidr") {
		cidr := d.Get("vpc_cidr").(string)
		err = updateVpcCidrOfQueue(client, queueName, cidr)
		if err != nil {
			return diag.Errorf("update cidr failed when updating dli queues: %s", err)
		}
	}

	return resourceDliQueueRead(ctx, d, meta)
}

func buildScaleActionParam(oldValue, newValue int) string {
	if oldValue > newValue {
		return actionScaleIn
	}
	return actionScaleOut
}

func validCuCount(val interface{}, key string) (warns []string, errs []error) {
	diviNum := 16
	warns, errs = validation.IntAtLeast(diviNum)(val, key)
	if len(errs) > 0 {
		return warns, errs
	}
	return validation.IntDivisibleBy(diviNum)(val, key)
}

func updateVpcCidrOfQueue(client *golangsdk.ServiceClient, queueName, cidr string) error {
	_, err := queues.UpdateCidr(client, queueName, queues.UpdateCidrOpts{Cidr: cidr})
	return err
}

func resourceQueueImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	err := d.Set("name", d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error saving resource name of the DLI queue: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
