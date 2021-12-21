package huaweicloud

import (
	"github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func DataSourceBandWidth() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBandWidthRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(5, 2000),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
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

func dataSourceBandWidthRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud vpc client: %s", err)
	}

	listOpts := bandwidths.ListOpts{
		ShareType: "WHOLE",
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		listOpts.EnterpriseProjectID = v.(string)
	}

	allBWs, err := bandwidths.List(vpcClient, listOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Unable to list huaweicloud bandwidths: %s", err)
	}
	if len(allBWs) == 0 {
		return fmtp.Errorf("No huaweicloud bandwidth was found")
	}

	// Filter bandwidths by "name"
	var bandList []bandwidths.BandWidth
	name := d.Get("name").(string)
	for _, band := range allBWs {
		if name == band.Name {
			bandList = append(bandList, band)
		}
	}
	if len(bandList) == 0 {
		return fmtp.Errorf("No huaweicloud bandwidth was found by name: %s", name)
	}

	// Filter bandwidths by "size"
	result := bandList[0]
	if v, ok := d.GetOk("size"); ok {
		var found bool
		for _, band := range bandList {
			if v.(int) == band.Size {
				found = true
				result = band
				break
			}
		}
		if !found {
			return fmtp.Errorf("No huaweicloud bandwidth was found by size: %d", v.(int))
		}
	}

	logp.Printf("[DEBUG] Retrieved huaweicloud bandwidth %s: %+v", result.ID, result)
	d.SetId(result.ID)
	d.Set("name", result.Name)
	d.Set("size", result.Size)
	d.Set("enterprise_project_id", result.EnterpriseProjectID)

	d.Set("share_type", result.ShareType)
	d.Set("bandwidth_type", result.BandwidthType)
	d.Set("charge_mode", result.ChargeMode)
	d.Set("status", result.Status)
	return nil
}
