package dc

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dc/v3/interfaces"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type (
	InterfaceType string
	ServiceType   string
	RouteMode     string
	AddressType   string
)

const (
	InterfaceTypePrivate InterfaceType = "private"

	ServiceTypeVpc  ServiceType = "VPC"
	ServiceTypeVgw  ServiceType = "VGW"
	ServiceTypeGdww ServiceType = "GDWW"
	ServiceTypeLgw  ServiceType = "LGW"

	RouteModeStatic RouteMode = "static"
	RouteModeBgp    RouteMode = "bgp"

	AddressTypeIpv4 AddressType = "ipv4"
	AddressTypeIpv6 AddressType = "ipv6"
)

// @API DC DELETE /v3/{project_id}/dcaas/virtual-interfaces/{interfaceId}
// @API DC GET /v3/{project_id}/dcaas/virtual-interfaces/{interfaceId}
// @API DC PUT /v3/{project_id}/dcaas/virtual-interfaces/{interfaceId}
// @API DC POST /v3/{project_id}/dcaas/virtual-interfaces
func ResourceVirtualInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualInterfaceCreate,
		ReadContext:   resourceVirtualInterfaceRead,
		UpdateContext: resourceVirtualInterfaceUpdate,
		DeleteContext: resourceVirtualInterfaceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the virtual interface is located.",
			},
			"direct_connect_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the direct connection associated with the virtual interface.",
			},
			"vgw_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the virtual gateway to which the virtual interface is connected.",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5\\w-.]*$"),
						"Only chinese and english letters, digits, hyphens (-), underscores (_) and dots (.) are "+
							"allowed."),
					validation.StringLenBetween(1, 64),
				),
				Description: "The name of the virtual interface.",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(InterfaceTypePrivate),
				}, false),
				Description: "The type of the virtual interface.",
			},
			"route_mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(RouteModeStatic),
					string(RouteModeBgp),
				}, false),
				Description: "The route mode of the virtual interface.",
			},
			"vlan": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 3999),
				Description:  "The VLAN for constom side.",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ingress bandwidth size of the virtual interface.",
			},
			"remote_ep_group": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The CIDR list of remote subnets.",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[^<>]*$`),
						"The angle brackets (< and >) are not allowed."),
					validation.StringLenBetween(0, 128),
				),
				Description: "The description of the virtual interface.",
			},
			"service_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(ServiceTypeVpc),
					string(ServiceTypeVgw),
					string(ServiceTypeGdww),
					string(ServiceTypeLgw),
				}, false),
				Description: "The service type of the virtual interface.",
			},
			"local_gateway_v4_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"remote_gateway_v4_ip"},
				Description:  "The IPv4 address of the virtual interface in cloud side.",
			},
			"remote_gateway_v4_ip": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"remote_gateway_v6_ip"},
				Description:   "The IPv4 address of the virtual interface in client side.",
			},
			"address_family": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(AddressTypeIpv4),
					string(AddressTypeIpv6),
				}, false),
				Description: "The address family type of the virtual interface.",
			},
			"local_gateway_v6_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"remote_gateway_v6_ip"},
				ExactlyOneOf: []string{"local_gateway_v4_ip"},
				Description:  "The IPv6 address of the virtual interface in cloud side.",
			},
			"remote_gateway_v6_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The IPv6 address of the virtual interface in client side.",
			},
			"asn": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntNotInSlice([]int{64512}),
				Description:  "The local BGP ASN in client side.",
			},
			"bgp_md5": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The (MD5) password for the local BGP.",
			},
			"enable_bfd": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the Bidirectional Forwarding Detection (BFD) function.",
			},
			"enable_nqa": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the Network Quality Analysis (NQA) function.",
			},
			"lag_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the link aggregation group (LAG) associated with the virtual interface.",
			},
			"resource_tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The project ID of another tenant which is used to create virtual interface across tenant.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The enterprise project ID to which the virtual interface belongs.",
			},
			// Attributes
			"device_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The attributed device ID.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the virtual interface.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the virtual interface.",
			},
			"tags": common.TagsSchema(),
		},
	}
}

func buildVirtualInterfaceCreateOpts(d *schema.ResourceData, cfg *config.Config) interfaces.CreateOpts {
	return interfaces.CreateOpts{
		VgwId:               d.Get("vgw_id").(string),
		Type:                d.Get("type").(string),
		RouteMode:           d.Get("route_mode").(string),
		Vlan:                d.Get("vlan").(int),
		Bandwidth:           d.Get("bandwidth").(int),
		RemoteEpGroup:       utils.ExpandToStringList(d.Get("remote_ep_group").([]interface{})),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		DirectConnectId:     d.Get("direct_connect_id").(string),
		ServiceType:         d.Get("service_type").(string),
		LocalGatewayV4Ip:    d.Get("local_gateway_v4_ip").(string),
		RemoteGatewayV4Ip:   d.Get("remote_gateway_v4_ip").(string),
		AddressFamily:       d.Get("address_family").(string),
		LocalGatewayV6Ip:    d.Get("local_gateway_v6_ip").(string),
		RemoteGatewayV6Ip:   d.Get("remote_gateway_v6_ip").(string),
		BgpAsn:              d.Get("asn").(int),
		BgpMd5:              d.Get("bgp_md5").(string),
		EnableBfd:           d.Get("enable_bfd").(bool),
		EnableNqa:           d.Get("enable_nqa").(bool),
		LagId:               d.Get("lag_id").(string),
		ResourceTenantId:    d.Get("resource_tenant_id").(string),
		EnterpriseProjectId: common.GetEnterpriseProjectID(d, cfg),
	}
}

func resourceVirtualInterfaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DcV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DC v3 client: %s", err)
	}

	opts := buildVirtualInterfaceCreateOpts(d, cfg)
	resp, err := interfaces.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating virtual interface: %s", err)
	}
	d.SetId(resp.ID)

	// create tags
	if err := utils.CreateResourceTags(client, d, "dc-vif", d.Id()); err != nil {
		return diag.Errorf("error setting tags of DC virtual interface %s: %s", d.Id(), err)
	}

	return resourceVirtualInterfaceRead(ctx, d, meta)
}

func resourceVirtualInterfaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DcV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DC v3 client: %s", err)
	}

	interfaceId := d.Id()
	resp, err := interfaces.Get(client, interfaceId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "virtual interface")
	}
	log.Printf("[DEBUG] The response of virtual interface is: %#v", resp)

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("vgw_id", resp.VgwId),
		d.Set("type", resp.Type),
		d.Set("route_mode", resp.RouteMode),
		d.Set("vlan", resp.Vlan),
		d.Set("bandwidth", resp.Bandwidth),
		d.Set("remote_ep_group", resp.RemoteEpGroup),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("direct_connect_id", resp.DirectConnectId),
		d.Set("service_type", resp.ServiceType),
		d.Set("local_gateway_v4_ip", resp.LocalGatewayV4Ip),
		d.Set("remote_gateway_v4_ip", resp.RemoteGatewayV4Ip),
		d.Set("address_family", resp.AddressFamily),
		d.Set("local_gateway_v6_ip", resp.LocalGatewayV6Ip),
		d.Set("remote_gateway_v6_ip", resp.RemoteGatewayV6Ip),
		d.Set("asn", resp.BgpAsn),
		d.Set("bgp_md5", resp.BgpMd5),
		d.Set("enable_bfd", resp.EnableBfd),
		d.Set("enable_nqa", resp.EnableNqa),
		d.Set("lag_id", resp.LagId),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("device_id", resp.DeviceId),
		d.Set("status", resp.Status),
		d.Set("created_at", resp.CreatedAt),
		utils.SetResourceTagsToState(d, client, "dc-vif", d.Id()),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving virtual interface fields: %s", err)
	}
	return nil
}

func closeVirtualInterfaceNetworkDetection(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		interfaceId = d.Id()
		opts        = interfaces.UpdateOpts{}
	)

	// At the same time, only one of BFD and NQA is enabled.
	if d.HasChange("enable_bfd") && !d.Get("enable_bfd").(bool) {
		opts.EnableBfd = utils.Bool(false)
	} else if d.HasChange("enable_nqa") && !d.Get("enable_nqa").(bool) {
		opts.EnableNqa = utils.Bool(false)
	}
	if reflect.DeepEqual(opts, interfaces.UpdateOpts{}) {
		return nil
	}

	_, err := interfaces.Update(client, interfaceId, opts)
	if err != nil {
		return fmt.Errorf("error closing network detection of the virtual interface (%s): %s", interfaceId, err)
	}
	return nil
}

func openVirtualInterfaceNetworkDetection(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		interfaceId     = d.Id()
		detectionOpened = false
		opts            = interfaces.UpdateOpts{}
	)

	if d.HasChange("enable_bfd") && d.Get("enable_bfd").(bool) {
		detectionOpened = true
		opts.EnableBfd = utils.Bool(true)
	}
	if d.HasChange("enable_nqa") && d.Get("enable_nqa").(bool) {
		// The enable requests of BFD and NQA cannot be sent at the same time.
		if detectionOpened {
			return fmt.Errorf("BFD and NQA cannot be enabled at the same time")
		}
		opts.EnableNqa = utils.Bool(true)
	}
	if reflect.DeepEqual(opts, interfaces.UpdateOpts{}) {
		return nil
	}

	_, err := interfaces.Update(client, interfaceId, opts)
	if err != nil {
		return fmt.Errorf("error opening network detection of the virtual interface (%s): %s", interfaceId, err)
	}
	return nil
}

func resourceVirtualInterfaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DcV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DC v3 client: %s", err)
	}

	if d.HasChanges("name", "description", "bandwidth", "remote_ep_group") {
		var (
			interfaceId = d.Id()

			opts = interfaces.UpdateOpts{
				Name:          d.Get("name").(string),
				Description:   utils.String(d.Get("description").(string)),
				Bandwidth:     d.Get("bandwidth").(int),
				RemoteEpGroup: utils.ExpandToStringList(d.Get("remote_ep_group").([]interface{})),
			}
		)

		_, err := interfaces.Update(client, interfaceId, opts)
		if err != nil {
			return diag.Errorf("error closing network detection of the virtual interface (%s): %s", interfaceId, err)
		}
	}
	if d.HasChanges("enable_bfd", "enable_nqa") {
		// BFD and NQA cannot be enabled at the same time.
		// When BFD (NQA) is enabled and NQA (BFD) is disabled, we need to disable BFD (NQA) first, and then enable NQA (BFD).
		// If the disable and enable requests are sent at the same time, an error will be reported.
		if err = closeVirtualInterfaceNetworkDetection(client, d); err != nil {
			return diag.FromErr(err)
		}
		if err = openVirtualInterfaceNetworkDetection(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// update tags
	tagErr := utils.UpdateResourceTags(client, d, "dc-vif", d.Id())
	if tagErr != nil {
		return diag.Errorf("error updating tags of DC virtual interface %s: %s", d.Id(), tagErr)
	}

	return resourceVirtualInterfaceRead(ctx, d, meta)
}

func resourceVirtualInterfaceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DcV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DC v3 client: %s", err)
	}

	interfaceId := d.Id()
	err = interfaces.Delete(client, interfaceId)
	if err != nil {
		return diag.Errorf("error deleting virtual interface (%s): %s", interfaceId, err)
	}

	return nil
}
