package vpc

import (
	"time"

	"github.com/chnsz/golangsdk"
	bandwidthsv1 "github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"
	"github.com/chnsz/golangsdk/openstack/networking/v2/bandwidths"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceVpcBandWidthV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcBandWidthV2Create,
		Read:   resourceVpcBandWidthV2Read,
		Update: resourceVpcBandWidthV2Update,
		Delete: resourceVpcBandWidthV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
			},
			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(5, 2000),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVpcBandWidthV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	NetworkingV1Client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating networking client: %s", err)
	}

	size := d.Get("size").(int)

	createOpts := bandwidths.CreateOpts{
		Name: d.Get("name").(string),
		Size: &size,
	}

	epsID := config.GetEnterpriseProjectID(d)
	if epsID != "" {
		createOpts.EnterpriseProjectId = epsID
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	b, err := bandwidths.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating Bandwidth: %s", err)
	}

	logp.Printf("[DEBUG] Waiting for Bandwidth (%s) to become available.", b.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"NORMAL"},
		Pending:    []string{"CREATING"},
		Refresh:    waitForBandwidth(NetworkingV1Client, b.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf(
			"Error waiting for Bandwidth (%s) to become ACTIVE for creation: %s",
			b.ID, err)
	}
	d.SetId(b.ID)

	return resourceVpcBandWidthV2Read(d, meta)
}

func resourceVpcBandWidthV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating networking client: %s", err)
	}

	var bandwidthOpts bandwidths.Bandwidth

	if d.HasChange("name") {
		bandwidthOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("size") {
		bandwidthOpts.Size = d.Get("size").(int)
	}

	if bandwidthOpts != (bandwidths.Bandwidth{}) {
		updateOpts := bandwidths.UpdateOpts{
			Bandwidth: bandwidthOpts,
		}
		_, err := bandwidths.Update(networkingClient, d.Id(), updateOpts)
		if err != nil {
			return fmtp.Errorf("Error updating Huaweicloud BandWidth (%s): %s", d.Id(), err)
		}
	}

	return resourceVpcBandWidthV2Read(d, meta)
}

func resourceVpcBandWidthV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating networking client: %s", err)
	}

	b, err := bandwidthsv1.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "bandwidth")
	}

	mErr := multierror.Append(
		d.Set("name", b.Name),
		d.Set("size", b.Size),
		d.Set("enterprise_project_id", b.EnterpriseProjectID),

		d.Set("share_type", b.ShareType),
		d.Set("bandwidth_type", b.BandwidthType),
		d.Set("charge_mode", b.ChargeMode),
		d.Set("status", b.Status),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.Errorf("Error setting vault fields: %s", err)
	}

	return nil
}

func resourceVpcBandWidthV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	NetworkingV1Client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating networking client: %s", err)
	}

	err = bandwidths.Delete(networkingClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud Bandwidth: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"NORMAL"},
		Target:     []string{"DELETED"},
		Refresh:    waitForBandwidth(NetworkingV1Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error deleting Bandwidth: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForBandwidth(networkingClient *golangsdk.ServiceClient, Id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		b, err := bandwidthsv1.Get(networkingClient, Id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return b, "DELETED", nil
			}
			return nil, "", err
		}

		logp.Printf("[DEBUG] HuaweiCloud Bandwidth (%s) current status: %s", b.ID, b.Status)
		return b, b.Status, nil
	}
}
