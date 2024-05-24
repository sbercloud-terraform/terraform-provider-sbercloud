package elb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/elb/v3/listeners"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceListenerV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceListenerV3Create,
		ReadContext:   resourceListenerV3Read,
		UpdateContext: resourceListenerV3Update,
		DeleteContext: resourceListenerV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP", "HTTP", "HTTPS",
				}, false),
			},

			"protocol_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				AtLeastOneOf: []string{"protocol_port", "port_ranges"},
			},

			"port_ranges": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_port": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"end_port": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"default_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"http2_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"forward_eip": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"forward_port": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"forward_request_port": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"forward_host": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"access_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"ip_group"},
				ValidateFunc: validation.StringInSlice([]string{
					"white", "black",
				}, true),
			},

			"ip_group": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"access_policy"},
			},

			"server_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"sni_certificate": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"ca_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tls_ciphers_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"idle_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"request_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"response_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"advanced_forwarding_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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

			"gzip_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"tags": common.TagsSchema(),
		},
	}
}

func resourceListenerV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	http2Enable := d.Get("http2_enable").(bool)
	var sniContainerRefs []string
	if raw, ok := d.GetOk("sni_certificate"); ok {
		for _, v := range raw.([]interface{}) {
			sniContainerRefs = append(sniContainerRefs, v.(string))
		}
	}
	enhanceL7policy := d.Get("advanced_forwarding_enabled").(bool)
	createOpts := listeners.CreateOpts{
		Protocol:               listeners.Protocol(d.Get("protocol").(string)),
		ProtocolPort:           d.Get("protocol_port").(int),
		LoadbalancerID:         d.Get("loadbalancer_id").(string),
		Name:                   d.Get("name").(string),
		DefaultPoolID:          d.Get("default_pool_id").(string),
		Description:            d.Get("description").(string),
		DefaultTlsContainerRef: d.Get("server_certificate").(string),
		CAContainerRef:         d.Get("ca_certificate").(string),
		TlsCiphersPolicy:       d.Get("tls_ciphers_policy").(string),
		PortRanges:             buildPortRanges(d.Get("port_ranges").(*schema.Set).List()),
		SniContainerRefs:       sniContainerRefs,
		Http2Enable:            &http2Enable,
		EnhanceL7policy:        &enhanceL7policy,
		ProtectionStatus:       d.Get("protection_status").(string),
		ProtectionReason:       d.Get("protection_reason").(string),
	}
	if v, ok := d.GetOk("gzip_enable"); ok {
		gzipEnable := v.(bool)
		createOpts.GzipEnable = &gzipEnable
	}
	if v, ok := d.GetOk("idle_timeout"); ok {
		idleTimeout := v.(int)
		createOpts.KeepaliveTimeout = &idleTimeout
	}
	if v, ok := d.GetOk("request_timeout"); ok {
		requestTimeout := v.(int)
		createOpts.ClientTimeout = &requestTimeout
	}
	if v, ok := d.GetOk("response_timeout"); ok {
		responseTimeout := v.(int)
		createOpts.MemberTimeout = &responseTimeout
	}
	if v, ok := d.GetOk("access_policy"); ok {
		createOpts.IpGroup = &listeners.IpGroup{
			Enable:    true,
			Type:      v.(string),
			IpGroupId: d.Get("ip_group").(string),
		}
	}
	forwardEip := d.Get("forward_eip").(bool)
	forwardPort := d.Get("forward_port").(bool)
	forwardRequestPort := d.Get("forward_request_port").(bool)
	forwardHost := d.Get("forward_host").(bool)
	createOpts.InsertHeaders = &listeners.InsertHeaders{
		ForwardedELBIP:   &forwardEip,
		ForwardedPort:    &forwardPort,
		ForwardedForPort: &forwardRequestPort,
		ForwardedHost:    &forwardHost,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Wait for LoadBalancer to become active before continuing
	loadBalancerID := createOpts.LoadbalancerID
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForElbV3LoadBalancer(ctx, elbClient, loadBalancerID, "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Attempting to create listener")
	listener, err := listeners.Create(elbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating listener: %s", err)
	}

	// Wait for LoadBalancer to become active again before continuing
	err = waitForElbV3LoadBalancer(ctx, elbClient, loadBalancerID, "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(listener.ID)

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		elbV2Client, err := cfg.ElbV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating ELB v2.0 client: %s", err)
		}
		tagList := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(elbV2Client, "listeners", listener.ID, tagList).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of ELB listener %s: %s", listener.ID, tagErr)
		}
	}

	return resourceListenerV3Read(ctx, d, meta)
}

func resourceListenerV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	listener, err := listeners.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "listener")
	}

	log.Printf("[DEBUG] Retrieved listener %s: %#v", d.Id(), listener)

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", listener.Name),
		d.Set("protocol", listener.Protocol),
		d.Set("description", listener.Description),
		d.Set("protocol_port", listener.ProtocolPort),
		d.Set("default_pool_id", listener.DefaultPoolID),
		d.Set("http2_enable", listener.Http2Enable),
		d.Set("forward_eip", listener.InsertHeaders.ForwardedELBIP),
		d.Set("forward_port", listener.InsertHeaders.ForwardedPort),
		d.Set("forward_request_port", listener.InsertHeaders.ForwardedForPort),
		d.Set("forward_host", listener.InsertHeaders.ForwardedHost),
		d.Set("sni_certificate", listener.SniContainerRefs),
		d.Set("server_certificate", listener.DefaultTlsContainerRef),
		d.Set("ca_certificate", listener.CAContainerRef),
		d.Set("tls_ciphers_policy", listener.TlsCiphersPolicy),
		d.Set("idle_timeout", listener.KeepaliveTimeout),
		d.Set("request_timeout", listener.ClientTimeout),
		d.Set("response_timeout", listener.MemberTimeout),
		d.Set("loadbalancer_id", listener.Loadbalancers[0].ID),
		d.Set("advanced_forwarding_enabled", listener.EnhanceL7policy),
		d.Set("protection_status", listener.ProtectionStatus),
		d.Set("protection_reason", listener.ProtectionReason),
		d.Set("gzip_enable", listener.GzipEnable),
	)

	var portRanges []map[string]interface{}
	for _, v := range listener.PortRanges {
		portRanges = append(portRanges, map[string]interface{}{
			"start_port": v.StartPort,
			"end_port":   v.EndPort,
		})
	}
	mErr = multierror.Append(mErr, d.Set("port_ranges", portRanges))

	if listener.IpGroup != (listeners.IpGroup{}) {
		mErr = multierror.Append(mErr,
			d.Set("access_policy", listener.IpGroup.Type),
			d.Set("ip_group", listener.IpGroup.IpGroupId),
		)
	} else {
		mErr = multierror.Append(mErr,
			d.Set("access_policy", ""),
			d.Set("ip_group", ""),
		)
	}

	// fetch tags
	if resourceTags, err := tags.Get(elbV2Client, "listeners", d.Id()).Extract(); err == nil {
		tagMap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	} else {
		log.Printf("[WARN] fetching tags of ELB listener failed: %s", err)
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Dedicated ELB listener fields: %s", err)
	}

	return nil
}

func resourceListenerV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	updateListenerChanges := []string{"name", "description", "ca_certificate", "default_pool_id", "idle_timeout",
		"request_timeout", "response_timeout", "server_certificate", "access_policy", "ip_group", "forward_eip",
		"forward_port", "forward_request_port", "forward_host", "tls_ciphers_policy", "sni_certificate",
		"http2_enable", "gzip_enable", "advanced_forwarding_enabled", "protection_status", "protection_reason"}
	if d.HasChanges(updateListenerChanges...) {
		err := updateListener(ctx, d, elbClient)
		if err != nil {
			return err
		}
	}

	// update tags
	if d.HasChange("tags") {
		elbV2Client, err := cfg.ElbV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating ELB 2.0 client: %s", err)
		}
		tagErr := utils.UpdateResourceTags(elbV2Client, d, "listeners", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of ELB listener:%s, err:%s", d.Id(), tagErr)
		}
	}

	return resourceListenerV3Read(ctx, d, meta)
}

func updateListener(ctx context.Context, d *schema.ResourceData, elbClient *golangsdk.ServiceClient) diag.Diagnostics {
	var updateOpts listeners.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		desc := d.Get("description").(string)
		updateOpts.Description = &desc
	}
	if d.HasChange("idle_timeout") {
		idleTimeout := d.Get("idle_timeout").(int)
		updateOpts.KeepaliveTimeout = &idleTimeout
	}
	if d.HasChange("request_timeout") {
		requestTimeout := d.Get("request_timeout").(int)
		updateOpts.ClientTimeout = &requestTimeout
	}
	if d.HasChange("response_timeout") {
		responseTimeout := d.Get("response_timeout").(int)
		updateOpts.MemberTimeout = &responseTimeout
	}
	if d.HasChange("default_pool_id") {
		updateOpts.DefaultPoolID = d.Get("default_pool_id").(string)
	}
	if d.HasChanges("access_policy", "ip_group") {
		updateOpts.IpGroup = &listeners.IpGroupUpdate{
			Type:      d.Get("access_policy").(string),
			IpGroupId: d.Get("ip_group").(string),
		}
	}
	if d.HasChanges("forward_eip", "forward_port", "forward_request_port", "forward_host") {
		forwardEip := d.Get("forward_eip").(bool)
		forwardPort := d.Get("forward_port").(bool)
		forwardRequestPort := d.Get("forward_request_port").(bool)
		forwardHost := d.Get("forward_host").(bool)
		updateOpts.InsertHeaders = &listeners.InsertHeaders{
			ForwardedELBIP:   &forwardEip,
			ForwardedPort:    &forwardPort,
			ForwardedForPort: &forwardRequestPort,
			ForwardedHost:    &forwardHost,
		}
	}
	if d.HasChange("ca_certificate") {
		caCert := d.Get("ca_certificate").(string)
		updateOpts.CAContainerRef = &caCert
	}
	if d.HasChange("tls_ciphers_policy") {
		tlsCiphersPolicy := d.Get("tls_ciphers_policy").(string)
		updateOpts.TlsCiphersPolicy = &tlsCiphersPolicy
	}
	if d.HasChange("server_certificate") {
		serverCert := d.Get("server_certificate").(string)
		updateOpts.DefaultTlsContainerRef = &serverCert
	}
	if d.HasChange("sni_certificate") {
		var sniContainerRefs []string
		if raw, ok := d.GetOk("sni_certificate"); ok {
			for _, v := range raw.([]interface{}) {
				sniContainerRefs = append(sniContainerRefs, v.(string))
			}
		}
		updateOpts.SniContainerRefs = &sniContainerRefs
	}
	if d.HasChange("http2_enable") {
		http2 := d.Get("http2_enable").(bool)
		updateOpts.Http2Enable = &http2
	}
	if d.HasChange("gzip_enable") {
		gzipEnable := d.Get("gzip_enable").(bool)
		updateOpts.GzipEnable = &gzipEnable
	}
	if d.HasChange("advanced_forwarding_enabled") {
		enhanceL7policy := d.Get("advanced_forwarding_enabled").(bool)
		updateOpts.EnhanceL7policy = &enhanceL7policy
	}
	if d.HasChange("protection_status") {
		updateOpts.ProtectionStatus = d.Get("protection_status").(string)
	}
	if d.HasChange("protection_reason") {
		protectionReason := d.Get("protection_reason").(string)
		updateOpts.ProtectionReason = &protectionReason
	}

	// Wait for LoadBalancer to become active before continuing
	loadBalancerID := d.Get("loadbalancer_id").(string)
	timeout := d.Timeout(schema.TimeoutUpdate)
	err := waitForElbV3LoadBalancer(ctx, elbClient, loadBalancerID, "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updating listener %s with options: %#v", d.Id(), updateOpts)
	_, err = listeners.Update(elbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating listener %s: %s", d.Id(), err)
	}

	// Wait for LoadBalancer to become active again before continuing
	err = waitForElbV3LoadBalancer(ctx, elbClient, loadBalancerID, "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceListenerV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	// Wait for LoadBalancer to become active before continuing
	loadBalancerID := d.Get("loadbalancer_id").(string)
	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForElbV3LoadBalancer(ctx, elbClient, loadBalancerID, "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Deleting listener %s", d.Id())
	if d.Get("force_delete").(bool) {
		if err = listeners.ForceDelete(elbClient, d.Id()).ExtractErr(); err != nil {
			return diag.Errorf("error deleting listener %s: %s", d.Id(), err)
		}
	} else {
		if err = listeners.Delete(elbClient, d.Id()).ExtractErr(); err != nil {
			return diag.Errorf("error deleting listener %s: %s", d.Id(), err)
		}
	}

	// Wait for LoadBalancer to become active again before continuing
	err = waitForElbV3LoadBalancer(ctx, elbClient, loadBalancerID, "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for Listener to delete
	err = waitForELBV3Listener(ctx, elbClient, d.Id(), "DELETED", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForELBV3Listener(ctx context.Context, elbClient *golangsdk.ServiceClient, id string, target string,
	pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for listener %s to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceELBV3ListenerRefreshFunc(elbClient, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("error: listener %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("error waiting for listener %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceELBV3ListenerRefreshFunc(elbClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		listener, err := listeners.Get(elbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		// The listener resource has no Status attribute, so a successful Get is the best we can do
		return listener, "ACTIVE", nil
	}
}

func buildPortRanges(rawPortRanges []interface{}) []listeners.PortRange {
	if len(rawPortRanges) == 0 {
		return nil
	}
	portRanges := make([]listeners.PortRange, 0)
	for _, rawPortRange := range rawPortRanges {
		if portRange, ok := rawPortRange.(map[string]interface{}); ok {
			portRanges = append(portRanges, listeners.PortRange{
				StartPort: portRange["start_port"].(int),
				EndPort:   portRange["end_port"].(int),
			})
		}
	}
	return portRanges
}
