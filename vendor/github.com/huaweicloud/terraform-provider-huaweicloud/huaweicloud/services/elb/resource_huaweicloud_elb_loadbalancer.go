package elb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/elb/v3/loadbalancers"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceLoadBalancerV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerV3Create,
		ReadContext:   resourceLoadBalancerV3Read,
		UpdateContext: resourceLoadBalancerV3Update,
		DeleteContext: resourceLoadBalancerV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"availability_zone": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"cross_vpc_backend": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"ipv4_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the IPv4 subnet ID of the subnet where the load balancer resides",
			},

			"ipv6_network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the ID of the subnet where the load balancer resides",
			},

			"ipv6_bandwidth_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ipv4_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ipv4_eip_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ConflictsWith: []string{
					"iptype", "bandwidth_charge_mode", "bandwidth_size", "sharetype", "bandwidth_id",
				},
			},

			"iptype": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"ipv4_eip_id"},
			},

			"bandwidth_charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				RequiredWith: []string{
					"iptype", "bandwidth_size", "sharetype",
				},
				ConflictsWith: []string{"ipv4_eip_id", "bandwidth_id"},
			},

			"sharetype": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				RequiredWith: []string{
					"iptype",
				},
				ConflictsWith: []string{"ipv4_eip_id"},
			},

			"bandwidth_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
				RequiredWith: []string{
					"iptype", "bandwidth_charge_mode", "sharetype",
				},
				ConflictsWith: []string{"ipv4_eip_id", "bandwidth_id"},
			},

			"bandwidth_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				RequiredWith: []string{
					"iptype",
				},
				ConflictsWith: []string{"ipv4_eip_id", "bandwidth_size", "bandwidth_charge_mode"},
			},

			"l4_flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"l7_flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"backend_subnets": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"protection_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"protection_reason": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tags": common.TagsSchema(),

			// charge info: charging_mode, period_unit, period, auto_renew, auto_pay
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
			"auto_pay":      common.SchemaAutoPay(nil),

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"ipv4_port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv4_eip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv6_eip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv6_eip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv6_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"autoscaling_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"min_l7_flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				RequiredWith: []string{
					"l7_flavor_id",
				},
			},
		},
	}
}

func resourceLoadBalancerV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	iPTargetEnable := d.Get("cross_vpc_backend").(bool)
	createOpts := loadbalancers.CreateOpts{
		AvailabilityZoneList: utils.ExpandToStringListBySet(d.Get("availability_zone").(*schema.Set)),
		IPTargetEnable:       &iPTargetEnable,
		VpcID:                d.Get("vpc_id").(string),
		VipSubnetID:          d.Get("ipv4_subnet_id").(string),
		IpV6VipSubnetID:      d.Get("ipv6_network_id").(string),
		VipAddress:           d.Get("ipv4_address").(string),
		L4Flavor:             d.Get("l4_flavor_id").(string),
		L7Flavor:             d.Get("l7_flavor_id").(string),
		ProtectionStatus:     d.Get("protection_status").(string),
		ProtectionReason:     d.Get("protection_reason").(string),
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		EnterpriseProjectID:  common.GetEnterpriseProjectID(d, cfg),
	}

	if v, ok := d.GetOk("backend_subnets"); ok {
		createOpts.ElbSubnetIds = utils.ExpandToStringList(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("ipv6_bandwidth_id"); ok {
		createOpts.IPV6Bandwidth = &loadbalancers.BandwidthRef{
			ID: v.(string),
		}
	}
	if v, ok := d.GetOk("ipv4_eip_id"); ok {
		createOpts.PublicIPIds = []string{v.(string)}
	}
	if v, ok := d.GetOk("iptype"); ok {
		createOpts.PublicIP = &loadbalancers.PublicIP{
			IPVersion:   4,
			NetworkType: v.(string),
			Bandwidth: loadbalancers.Bandwidth{
				Id:         d.Get("bandwidth_id").(string),
				Name:       d.Get("name").(string),
				Size:       d.Get("bandwidth_size").(int),
				ChargeMode: d.Get("bandwidth_charge_mode").(string),
				ShareType:  d.Get("sharetype").(string),
			},
		}
	}
	if v, ok := d.GetOk("autoscaling_enabled"); ok {
		createOpts.AutoScaling = &loadbalancers.AutoScaling{
			Enable:      v.(bool),
			MinL7Flavor: d.Get("min_l7_flavor_id").(string),
		}
	}

	var loadBalancerID string
	if d.Get("charging_mode").(string) == "prePaid" {
		autoRenew, _ := strconv.ParseBool(d.Get("auto_renew").(string))
		prepaidOpts := loadbalancers.PrepaidOpts{
			PeriodType: d.Get("period_unit").(string),
			PeriodNum:  d.Get("period").(int),
			AutoRenew:  autoRenew,
		}
		if d.Get("auto_pay").(string) != "false" {
			prepaidOpts.AutoPay = true
		}
		createOpts.PrepaidOpts = &prepaidOpts

		log.Printf("[DEBUG] Create Options: %#v", createOpts)
		resp, err := loadbalancers.Create(elbClient, createOpts).ExtractPrepaid()
		if err != nil {
			return diag.Errorf("error creating prepaid LoadBalancer: %s", err)
		}

		// wait for the order to be completed.
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, resp.OrderID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("the order is not completed while creating ELB LoadBalancer (%s): %#v", resp.LoadBalancerID, err)
		}
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, resp.OrderID,
			d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		loadBalancerID = resourceId
	} else {
		log.Printf("[DEBUG] Create Options: %#v", createOpts)
		lb, err := loadbalancers.Create(elbClient, createOpts).Extract()
		if err != nil {
			return diag.Errorf("error creating LoadBalancer: %s", err)
		}

		loadBalancerID = lb.ID
	}

	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForElbV3LoadBalancer(ctx, elbClient, loadBalancerID, "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// set the ID on the resource
	d.SetId(loadBalancerID)

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		elbV2Client, err := cfg.ElbV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating ELB 2.0 client: %s", err)
		}
		tagList := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(elbV2Client, "loadbalancers", d.Id(), tagList).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of LoadBalancer %s: %s", d.Id(), tagErr)
		}
	}

	return resourceLoadBalancerV3Read(ctx, d, meta)
}

func resourceLoadBalancerV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	// client for fetching tags
	elbV2Client, err := cfg.ElbV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB 2.0 client: %s", err)
	}

	lb, err := loadbalancers.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "LoadBalancer")
	}

	log.Printf("[DEBUG] Retrieved LoadBalancer %s: %#v", d.Id(), lb)

	mErr := multierror.Append(nil,
		d.Set("name", lb.Name),
		d.Set("description", lb.Description),
		d.Set("availability_zone", lb.AvailabilityZoneList),
		d.Set("cross_vpc_backend", lb.IpTargetEnable),
		d.Set("vpc_id", lb.VpcID),
		d.Set("ipv4_subnet_id", lb.VipSubnetCidrID),
		d.Set("ipv6_network_id", lb.Ipv6VipVirsubnetID),
		d.Set("ipv4_address", lb.VipAddress),
		d.Set("ipv4_port_id", lb.VipPortID),
		d.Set("ipv6_address", lb.Ipv6VipAddress),
		d.Set("l4_flavor_id", lb.L4FlavorID),
		d.Set("l7_flavor_id", lb.L7FlavorID),
		d.Set("region", cfg.GetRegion(d)),
		d.Set("enterprise_project_id", lb.EnterpriseProjectID),
		d.Set("autoscaling_enabled", lb.AutoScaling.Enable),
		d.Set("min_l7_flavor_id", lb.AutoScaling.MinL7Flavor),
		d.Set("backend_subnets", lb.ElbVirsubnetIDs),
		d.Set("protection_status", lb.ProtectionStatus),
		d.Set("protection_reason", lb.ProtectionReason),
	)

	for _, eip := range lb.Eips {
		if eip.IpVersion == 4 {
			mErr = multierror.Append(mErr,
				d.Set("ipv4_eip_id", eip.EipID),
				d.Set("ipv4_eip", eip.EipAddress),
			)
		} else if eip.IpVersion == 6 {
			mErr = multierror.Append(mErr,
				d.Set("ipv6_eip_id", eip.EipID),
				d.Set("ipv6_eip", eip.EipAddress),
			)
		}
	}

	// set charging_mode according to billing_info
	if len(lb.BillingInfo) > 0 {
		mErr = multierror.Append(mErr,
			d.Set("charging_mode", "prePaid"),
		)
	} else {
		mErr = multierror.Append(mErr,
			d.Set("charging_mode", "postPaid"),
		)
	}

	// fetch tags
	if resourceTags, err := tags.Get(elbV2Client, "loadbalancers", d.Id()).Extract(); err == nil {
		tagMap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	} else {
		log.Printf("[WARN] Fetching tags of ELB LoadBalancer failed: %s", err)
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Dedicated ELB LoadBalancer fields: %s", err)
	}

	return nil
}

func resourceLoadBalancerV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	updateLoadBalancerChanges := []string{"name", "description", "cross_vpc_backend", "ipv4_subnet_id", "ipv6_network_id",
		"ipv6_bandwidth_id", "ipv4_address", "l4_flavor_id", "l7_flavor_id", "autoscaling_enabled", "min_l7_flavor_id",
		"protection_status", "protection_reason",
	}

	if d.HasChanges(updateLoadBalancerChanges...) {
		updateOpts := buildUpdateLoadBalancerBodyParams(d)
		err := updateLoadBalancer(ctx, d, cfg, updateOpts, elbClient)
		if err != nil {
			return err
		}
	}

	// backend_subnets and cross_vpc_backend can not be updated at the same time, otherwise, an error will occur
	if d.HasChange("backend_subnets") {
		updateOpts := loadbalancers.UpdateOpts{
			ElbSubnetIds: utils.ExpandToStringList(d.Get("backend_subnets").(*schema.Set).List()),
		}
		// if the value of vip_subnet_cidr_id and ipv6_vip_virsubnet_id are null, then they will be unbound
		if v, ok := d.GetOk("ipv4_subnet_id"); ok {
			updateOpts.VipSubnetID = utils.String(v.(string))
		}
		if v, ok := d.GetOk("ipv6_network_id"); ok {
			updateOpts.IpV6VipSubnetID = utils.String(v.(string))
		}
		_, err = loadbalancers.Update(elbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating ELB LoadBalancer: %s", err)
		}
		// Wait for LoadBalancer to become active before continuing
		err = waitForElbV3LoadBalancer(ctx, elbClient, d.Id(), "ACTIVE", nil, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return diag.Errorf("error updating the auto-renew of the LoadBalancer (%s): %s", d.Id(), err)
		}
	}
	// update tags
	if d.HasChange("tags") {
		elbV2Client, err := cfg.ElbV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating ELB 2.0 client: %s", err)
		}
		tagErr := utils.UpdateResourceTags(elbV2Client, d, "loadbalancers", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of LoadBalancer:%s, err:%s", d.Id(), tagErr)
		}
	}

	// update availability zone
	if d.HasChange("availability_zone") {
		if err = updateAvailabilityZone(ctx, cfg, elbClient, d); err != nil {
			return diag.Errorf("error updating availability zone of LoadBalancer:%s, err:%s", d.Id(), err)
		}
	}

	return resourceLoadBalancerV3Read(ctx, d, meta)
}

func updateAvailabilityZone(ctx context.Context, cfg *config.Config, elbClient *golangsdk.ServiceClient,
	d *schema.ResourceData) error {
	oldAvailabilityZone, newAvailabilityZone := d.GetChange("availability_zone")
	addList := newAvailabilityZone.(*schema.Set).Difference(oldAvailabilityZone.(*schema.Set)).List()
	rmList := oldAvailabilityZone.(*schema.Set).Difference(newAvailabilityZone.(*schema.Set)).List()
	err := addAvailabilityZone(ctx, cfg, elbClient, d, addList)
	if err != nil {
		return err
	}
	err = removeAvailabilityZone(ctx, cfg, elbClient, d, rmList)
	if err != nil {
		return err
	}
	return nil
}

func addAvailabilityZone(ctx context.Context, cfg *config.Config, elbClient *golangsdk.ServiceClient,
	d *schema.ResourceData, addList []interface{}) error {
	if len(addList) > 0 {
		updateOpts := loadbalancers.AvailabilityZoneOpts{
			AvailabilityZoneList: utils.ExpandToStringList(addList),
		}
		resp, err := loadbalancers.AddAvailabilityZone(elbClient, d.Id(), updateOpts).ExtractPrepaid()
		if err != nil {
			return err
		}
		if len(resp.OrderID) > 0 {
			// wait for the order to be completed.
			bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
			if err != nil {
				return err
			}
			err = common.WaitOrderComplete(ctx, bssClient, resp.OrderID, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
		}
		// Wait for LoadBalancer to become active before continuing
		err = waitForElbV3LoadBalancer(ctx, elbClient, d.Id(), "ACTIVE", nil, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}
	return nil
}

func removeAvailabilityZone(ctx context.Context, cfg *config.Config, elbClient *golangsdk.ServiceClient,
	d *schema.ResourceData, rmList []interface{}) error {
	if len(rmList) > 0 {
		updateOpts := loadbalancers.AvailabilityZoneOpts{
			AvailabilityZoneList: utils.ExpandToStringList(rmList),
		}
		resp, err := loadbalancers.RemoveAvailabilityZone(elbClient, d.Id(), updateOpts).ExtractPrepaid()
		if err != nil {
			return err
		}
		if len(resp.OrderID) > 0 {
			// wait for the order to be completed.
			bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
			if err != nil {
				return err
			}
			err = common.WaitOrderComplete(ctx, bssClient, resp.OrderID, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return err
			}
		}
		// Wait for LoadBalancer to become active before continuing
		err = waitForElbV3LoadBalancer(ctx, elbClient, d.Id(), "ACTIVE", nil, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}
	return nil
}

func buildUpdateLoadBalancerBodyParams(d *schema.ResourceData) loadbalancers.UpdateOpts {
	var updateOpts loadbalancers.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("cross_vpc_backend") {
		iPTargetEnable := d.Get("cross_vpc_backend").(bool)
		updateOpts.IPTargetEnable = &iPTargetEnable
	}
	if d.HasChange("ipv4_address") {
		updateOpts.VipAddress = d.Get("ipv4_address").(string)
	}
	if d.HasChange("l4_flavor_id") {
		updateOpts.L4Flavor = d.Get("l4_flavor_id").(string)
	}
	if d.HasChange("l7_flavor_id") {
		updateOpts.L7Flavor = d.Get("l7_flavor_id").(string)
	}
	if d.HasChange("ipv6_bandwidth_id") {
		if v, ok := d.GetOk("ipv6_bandwidth_id"); ok {
			bw := v.(string)
			updateOpts.IPV6Bandwidth = &loadbalancers.UBandwidthRef{
				ID: &bw,
			}
		} else {
			updateOpts.IPV6Bandwidth = &loadbalancers.UBandwidthRef{}
		}
	}
	if d.HasChange("protection_status") {
		updateOpts.ProtectionStatus = d.Get("protection_status").(string)
	}
	if d.HasChange("protection_reason") {
		protectionReason := d.Get("protection_reason").(string)
		updateOpts.ProtectionReason = &protectionReason
	}

	// always with below values as null is meaningful
	if v, ok := d.GetOk("ipv4_subnet_id"); ok {
		vipSubnetID := v.(string)
		updateOpts.VipSubnetID = &vipSubnetID
	}
	if v, ok := d.GetOk("ipv6_network_id"); ok {
		v6SubnetID := v.(string)
		updateOpts.IpV6VipSubnetID = &v6SubnetID
	}

	if d.HasChange("autoscaling_enabled") {
		autoscalingEnabled := d.Get("autoscaling_enabled").(bool)
		updateOpts.AutoScaling = &loadbalancers.AutoScaling{
			Enable: autoscalingEnabled,
		}
		if autoscalingEnabled {
			updateOpts.AutoScaling.MinL7Flavor = d.Get("min_l7_flavor_id").(string)
		} else {
			updateOpts.L4Flavor = d.Get("l4_flavor_id").(string)
			updateOpts.L7Flavor = d.Get("l7_flavor_id").(string)
			updateOpts.AutoScaling.MinL7Flavor = ""
		}
	} else if d.HasChange("min_l7_flavor_id") && d.Get("autoscaling_enabled").(bool) {
		if autoscalingEnabled := d.Get("autoscaling_enabled").(bool); autoscalingEnabled {
			updateOpts.AutoScaling.MinL7Flavor = d.Get("min_l7_flavor_id").(string)
		}
	}

	log.Printf("[DEBUG] Updating LoadBalancer %s with options: %#v", d.Id(), updateOpts)

	return updateOpts
}

func updateLoadBalancer(ctx context.Context, d *schema.ResourceData, cfg *config.Config, updateOpts loadbalancers.UpdateOpts,
	elbClient *golangsdk.ServiceClient) diag.Diagnostics {
	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutUpdate)
	err := waitForElbV3LoadBalancer(ctx, elbClient, d.Id(), "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("charging_mode").(string) == "prePaid" && d.HasChanges("l4_flavor_id", "l7_flavor_id") {
		prepaidOpts := loadbalancers.PrepaidOpts{
			PeriodType: d.Get("period_unit").(string),
			PeriodNum:  d.Get("period").(int),
		}
		if d.Get("auto_pay").(string) != "false" {
			prepaidOpts.AutoPay = true
		}
		updateOpts.PrepaidOpts = &prepaidOpts

		resp, err := loadbalancers.Update(elbClient, d.Id(), updateOpts).ExtractPrepaid()
		if err != nil {
			return diag.Errorf("error updating prepaid LoadBalancer: %s", err)
		}

		// wait for the order to be completed.
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, resp.OrderID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("the order is not completed while updating ELB LoadBalancer (%s): %#v",
				resp.LoadBalancerID, err)
		}
	} else {
		_, err = loadbalancers.Update(elbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating ELB LoadBalancer: %s", err)
		}
	}
	// Wait for LoadBalancer to become active before continuing
	err = waitForElbV3LoadBalancer(ctx, elbClient, d.Id(), "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceLoadBalancerV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	elbClient, err := cfg.ElbV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	log.Printf("[DEBUG] Deleting LoadBalancer %s", d.Id())

	if d.Get("charging_mode").(string) == "prePaid" {
		// Unsubscribe the prepaid LoadBalancer will automatically delete it
		if err = common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
			return diag.Errorf("error unsubscribing ELB LoadBalancer : %s", err)
		}
	} else {
		if d.Get("force_delete").(bool) {
			if err = loadbalancers.ForceDelete(elbClient, d.Id()).ExtractErr(); err != nil {
				return diag.Errorf("error deleting ELB LoadBalancer: %s", err)
			}
		} else {
			if err = loadbalancers.Delete(elbClient, d.Id()).ExtractErr(); err != nil {
				return diag.Errorf("error deleting ELB LoadBalancer: %s", err)
			}
		}
	}

	// Wait for LoadBalancer to become delete
	timeout := d.Timeout(schema.TimeoutDelete)
	pending := []string{"PENDING_UPDATE", "PENDING_DELETE", "ACTIVE"}
	err = waitForElbV3LoadBalancer(ctx, elbClient, d.Id(), "DELETED", pending, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	// delete the EIP if necessary
	eipID := d.Get("ipv4_eip_id").(string)
	if _, ok := d.GetOk("iptype"); ok && eipID != "" {
		eipClient, err := cfg.NetworkingV1Client(region)
		if err == nil {
			if eipErr := eips.Delete(eipClient, eipID).ExtractErr(); eipErr != nil {
				if _, ok := err.(golangsdk.ErrDefault404); !ok {
					eipDiag := diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "failed to delete EIP",
						Detail:   fmt.Sprintf("failed to delete EIP %s: %s", eipID, eipErr),
					}
					diags = append(diags, eipDiag)
				}
			}
		} else {
			clientDiag := diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "failed to create VPC client",
				Detail:   fmt.Sprintf("failed to create VPC client: %s", err),
			}
			diags = append(diags, clientDiag)
		}
	}

	return diags
}

func waitForElbV3LoadBalancer(ctx context.Context, elbClient *golangsdk.ServiceClient,
	id string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for LoadBalancer %s to become %s", id, target)

	stateConf := &resource.StateChangeConf{
		Target:       []string{target},
		Pending:      pending,
		Refresh:      resourceElbV3LoadBalancerRefreshFunc(elbClient, id),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("error: LoadBalancer %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("error waiting for LoadBalancer %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceElbV3LoadBalancerRefreshFunc(elbClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		lb, err := loadbalancers.Get(elbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}
		return lb, lb.ProvisioningStatus, nil
	}
}
