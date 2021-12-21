package huaweicloud

import (
	"strconv"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/networks"
	"github.com/chnsz/golangsdk/openstack/networking/v2/subnets"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func dataSourceNetworkingNetworkV2() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourceNetworkingNetworkV2Read,
		DeprecationMessage: "use huaweicloud_vpc_subnet data source instead",

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"matching_subnet_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HW_PROJECT_ID",
					"OS_TENANT_ID",
				}, ""),
				Description: descriptions["tenant_id"],
			},
			"admin_state_up": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shared": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkingNetworkV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	listOpts := networks.ListOpts{
		ID:       d.Get("network_id").(string),
		Name:     d.Get("name").(string),
		TenantID: d.Get("tenant_id").(string),
	}

	if v, ok := d.GetOk("status"); ok {
		listOpts.Status = v.(string)
	}

	pages, err := networks.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return err
	}

	allNetworks, err := networks.ExtractNetworks(pages)
	if err != nil {
		return fmtp.Errorf("Unable to retrieve networks: %s", err)
	}

	var refinedNetworks []networks.Network
	if cidr := d.Get("matching_subnet_cidr").(string); cidr != "" {
		for _, n := range allNetworks {
			for _, s := range n.Subnets {
				subnet, err := subnets.Get(networkingClient, s).Extract()
				if err != nil {
					if _, ok := err.(golangsdk.ErrDefault404); ok {
						continue
					}
					return fmtp.Errorf("Unable to retrieve network subnet: %s", err)
				}
				if cidr == subnet.CIDR {
					refinedNetworks = append(refinedNetworks, n)
				}
			}
		}
	} else {
		refinedNetworks = allNetworks
	}

	if len(refinedNetworks) < 1 {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedNetworks) > 1 {
		return fmtp.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	network := refinedNetworks[0]

	logp.Printf("[DEBUG] Retrieved Network %s: %+v", network.ID, network)
	d.SetId(network.ID)

	d.Set("name", network.Name)
	d.Set("admin_state_up", strconv.FormatBool(network.AdminStateUp))
	d.Set("shared", strconv.FormatBool(network.Shared))
	d.Set("tenant_id", network.TenantID)
	d.Set("region", GetRegion(d, config))

	return nil
}
