package sbercloud

import (
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dns/v2/zones"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func DataSourceDNSZoneV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDNSZoneV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 2147483647),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"masters": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
		},
	}
}

func dataSourceDNSZoneV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	dnsClient, err := config.DnsV2Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	listOpts := zones.ListOpts{
		Description: d.Get("description").(string),
		Email:       d.Get("email").(string),
		Name:        d.Get("name").(string),
		TTL:         d.Get("ttl").(int),
		Type:        d.Get("zone_type").(string),
	}

	pages, err := zones.List(dnsClient, listOpts).AllPages()
	if err != nil {
		return fmtp.Errorf("Unable to retrieve zones list: %s", err)
	}

	refinedZones, err := zones.ExtractZones(pages)
	if err != nil {
		return fmtp.Errorf("Unable to retrieve vpc routes: %s", err)
	}

	total := len(refinedZones)
	if total < 1 {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	if total > 1 {
		return fmtp.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	zoneInfo := refinedZones[0]
	logp.Printf("[DEBUG] Retrieved Zone %s: %+v", zoneInfo.ID, zoneInfo)
	d.SetId(zoneInfo.ID)

	d.Set("name", zoneInfo.Name)
	d.Set("email", zoneInfo.Email)
	d.Set("description", zoneInfo.Description)
	d.Set("ttl", zoneInfo.TTL)
	if err = d.Set("masters", zoneInfo.Masters); err != nil {
		return fmtp.Errorf("[DEBUG] Error saving masters to state for HuaweiCloud DNS zone (%s): %s", d.Id(), err)
	}
	d.Set("region", region)
	d.Set("zone_type", zoneInfo.ZoneType)
	d.Set("enterprise_project_id", zoneInfo.EnterpriseProjectID)

	// save tags
	if resourceType, err := utils.GetDNSZoneTagType(zoneInfo.ZoneType); err == nil {
		resourceTags, err := tags.Get(dnsClient, resourceType, d.Id()).Extract()
		if err == nil {
			tagmap := utils.TagsToMap(resourceTags.Tags)
			d.Set("tags", tagmap)
		} else {
			logp.Printf("[WARN] Error fetching HuaweiCloud DNS zone tags: %s", err)
		}
	}

	return nil
}
