package huaweicloud

import (
	"strings"

	"github.com/chnsz/golangsdk/openstack/obs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceObsBucketPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceObsBucketPolicyPut,
		Read:   resourceObsBucketPolicyRead,
		Update: resourceObsBucketPolicyPut,
		Delete: resourceObsBucketPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceObsBucketImport,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     utils.ValidateJsonString,
				DiffSuppressFunc: utils.SuppressEquivalentAwsPolicyDiffs,
			},
			"policy_format": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "obs",
			},
		},
	}
}

func resourceObsBucketPolicyPut(d *schema.ResourceData, meta interface{}) error {
	var err error
	var obsClient *obs.ObsClient
	config := meta.(*config.Config)

	format := d.Get("policy_format").(string)
	if format == "obs" {
		obsClient, err = config.ObjectStorageClientWithSignature(GetRegion(d, config))
	} else {
		obsClient, err = config.ObjectStorageClient(GetRegion(d, config))
	}
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	policy := d.Get("policy").(string)
	logp.Printf("[DEBUG] OBS bucket: %s, set policy: %s", bucket, policy)

	params := &obs.SetBucketPolicyInput{
		Bucket: bucket,
		Policy: policy,
	}
	if _, err = obsClient.SetBucketPolicy(params); err != nil {
		return getObsError("Error setting OBS bucket policy", bucket, err)
	}

	// seem bucket as the policy id
	d.SetId(bucket)
	return nil
}

func resourceObsBucketPolicyRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var obsClient *obs.ObsClient
	config := meta.(*config.Config)

	format := d.Get("policy_format").(string)
	logp.Printf("[DEBUG] obs bucket policy format: %s", format)
	if format == "obs" {
		obsClient, err = config.ObjectStorageClientWithSignature(GetRegion(d, config))
	} else {
		obsClient, err = config.ObjectStorageClient(GetRegion(d, config))
	}
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	// set bucket from the policy id
	d.Set("bucket", d.Id())

	logp.Printf("[DEBUG] read policy for obs bucket: %s", d.Id())
	output, err := obsClient.GetBucketPolicy(d.Id())

	var pol string
	if err == nil && output.Policy != "" {
		pol = output.Policy
	}
	if err := d.Set("policy", pol); err != nil {
		return err
	}

	return nil
}

func resourceObsBucketPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	var err error
	var obsClient *obs.ObsClient
	config := meta.(*config.Config)

	format := d.Get("policy_format").(string)
	if format == "obs" {
		obsClient, err = config.ObjectStorageClientWithSignature(GetRegion(d, config))
	} else {
		obsClient, err = config.ObjectStorageClient(GetRegion(d, config))
	}
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud OBS client: %s", err)
	}

	bucket := d.Get("bucket").(string)

	logp.Printf("[DEBUG] OBS bucket: %s, delete policy", bucket)
	_, err = obsClient.DeleteBucketPolicy(bucket)
	if err != nil {
		return getObsError("Error deleting policy of OBS bucket %s: %s", bucket, err)
	}

	return nil
}

func resourceObsBucketImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var policyFormat = "obs"
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) == 2 {
		policyFormat = parts[1]
	}

	d.SetId(parts[0])
	d.Set("policy_format", policyFormat)

	return []*schema.ResourceData{d}, nil
}
