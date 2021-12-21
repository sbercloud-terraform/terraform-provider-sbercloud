package huaweicloud

import (
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/networking/v2/subnets"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func dataSourceNetworkingSubnetV2() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourceNetworkingSubnetV2Read,
		DeprecationMessage: "use huaweicloud_vpc_subnet data source instead",

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"dhcp_enabled": {
				Type:          schema.TypeBool,
				ConflictsWith: []string{"dhcp_disabled"},
				Optional:      true,
			},

			"dhcp_disabled": {
				Type:          schema.TypeBool,
				ConflictsWith: []string{"dhcp_enabled"},
				Optional:      true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"HW_PROJECT_ID",
					"OS_TENANT_ID",
				}, ""),
				Description: descriptions["tenant_id"],
			},

			"ip_version": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value != 4 && value != 6 {
						errors = append(errors, fmtp.Errorf(
							"Only 4 and 6 are supported values for 'ip_version'"))
					}
					return
				},
			},

			"gateway_ip": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			// Computed values
			"allocation_pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"enable_dhcp": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"dns_nameservers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"host_routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ipv6_address_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: utils.ValidateSubnetV2IPv6Mode,
			},
			"ipv6_ra_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: utils.ValidateSubnetV2IPv6Mode,
			},
		},
	}
}

func dataSourceNetworkingSubnetV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	listOpts := subnets.ListOpts{}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if _, ok := d.GetOk("dhcp_enabled"); ok {
		enableDHCP := true
		listOpts.EnableDHCP = &enableDHCP
	}

	if _, ok := d.GetOk("dhcp_disabled"); ok {
		enableDHCP := false
		listOpts.EnableDHCP = &enableDHCP
	}

	if v, ok := d.GetOk("network_id"); ok {
		listOpts.NetworkID = v.(string)
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		listOpts.TenantID = v.(string)
	}

	if v, ok := d.GetOk("ip_version"); ok {
		listOpts.IPVersion = v.(int)
	}

	if v, ok := d.GetOk("gateway_ip"); ok {
		listOpts.GatewayIP = v.(string)
	}

	if v, ok := d.GetOk("cidr"); ok {
		listOpts.CIDR = v.(string)
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		listOpts.ID = v.(string)
	}

	if v, ok := d.GetOk("ipv6_address_mode"); ok {
		listOpts.IPv6AddressMode = v.(string)
	}

	if v, ok := d.GetOk("ipv6_ra_mode"); ok {
		listOpts.IPv6RAMode = v.(string)
	}

	pages, err := subnets.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return fmtp.Errorf("Unable to retrieve subnets: %s", err)
	}

	allSubnets, err := subnets.ExtractSubnets(pages)
	if err != nil {
		return fmtp.Errorf("Unable to extract subnets: %s", err)
	}

	if len(allSubnets) < 1 {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allSubnets) > 1 {
		return fmtp.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	subnet := allSubnets[0]

	logp.Printf("[DEBUG] Retrieved Subnet %s: %+v", subnet.ID, subnet)
	d.SetId(subnet.ID)

	d.Set("name", subnet.Name)
	d.Set("tenant_id", subnet.TenantID)
	d.Set("network_id", subnet.NetworkID)
	d.Set("cidr", subnet.CIDR)
	d.Set("ip_version", subnet.IPVersion)
	d.Set("ipv6_address_mode", subnet.IPv6AddressMode)
	d.Set("ipv6_ra_mode", subnet.IPv6RAMode)
	d.Set("gateway_ip", subnet.GatewayIP)
	d.Set("enable_dhcp", subnet.EnableDHCP)
	d.Set("region", GetRegion(d, config))

	err = d.Set("dns_nameservers", subnet.DNSNameservers)
	if err != nil {
		logp.Printf("[DEBUG] Unable to set dns_nameservers: %s", err)
	}

	// Set the host_routes
	var hostRoutes []map[string]interface{} = make([]map[string]interface{}, len(subnet.HostRoutes))
	for i, v := range subnet.HostRoutes {
		routes := make(map[string]interface{})
		routes["destination_cidr"] = v.DestinationCIDR
		routes["next_hop"] = v.NextHop
		hostRoutes[i] = routes
	}
	if err = d.Set("host_routes", hostRoutes); err != nil {
		return fmtp.Errorf("Saving host_routes failed: %s", err)
	}

	// Set the allocation_pools
	var allocationPools []map[string]interface{}
	for _, v := range subnet.AllocationPools {
		pool := make(map[string]interface{})
		pool["start"] = v.Start
		pool["end"] = v.End

		allocationPools = append(allocationPools, pool)
	}
	err = d.Set("allocation_pools", allocationPools)
	if err != nil {
		logp.Printf("[DEBUG] Unable to set allocation_pools: %s", err)
	}

	return nil
}
