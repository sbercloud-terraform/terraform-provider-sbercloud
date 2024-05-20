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
	"github.com/chnsz/golangsdk/openstack/elb/v3/pools"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourcePoolV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePoolV3Create,
		ReadContext:   resourcePoolV3Read,
		UpdateContext: resourcePoolV3Update,
		DeleteContext: resourcePoolV3Delete,
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
					"TCP", "UDP", "HTTP", "HTTPS", "QUIC",
				}, false),
			},
			"loadbalancer_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				AtLeastOneOf: []string{"loadbalancer_id", "listener_id", "type"},
			},
			"listener_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				AtLeastOneOf: []string{"loadbalancer_id", "listener_id", "type"},
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"loadbalancer_id", "listener_id", "type"},
			},
			"lb_method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP", "QUIC_CID",
				}, false),
			},
			"persistence": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SOURCE_IP", "HTTP_COOKIE", "APP_COOKIE",
							}, false),
						},
						"cookie_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
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
			"slow_start_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"slow_start_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"slow_start_enabled"},
			},
			"ip_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePoolV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	var persistence pools.SessionPersistence
	if p, ok := d.GetOk("persistence"); ok {
		persistence, err = buildPersistence(p)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	createOpts := pools.CreateOpts{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Protocol:         d.Get("protocol").(string),
		LoadbalancerID:   d.Get("loadbalancer_id").(string),
		ListenerID:       d.Get("listener_id").(string),
		LBMethod:         d.Get("lb_method").(string),
		ProtectionStatus: d.Get("protection_status").(string),
		ProtectionReason: d.Get("protection_reason").(string),
		Type:             d.Get("type").(string),
		VpcId:            d.Get("vpc_id").(string),
	}

	if v, ok := d.GetOk("slow_start_enabled"); ok {
		createOpts.SlowStart = &pools.SlowStart{
			Enable:   v.(bool),
			Duration: d.Get("slow_start_duration").(int),
		}
	}

	// Must omit if not set
	if persistence != (pools.SessionPersistence{}) {
		createOpts.Persistence = &persistence
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	pool, err := pools.Create(elbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating pool: %s", err)
	}

	d.SetId(pool.ID)

	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForElbV3Pool(ctx, elbClient, d.Id(), "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePoolV3Read(ctx, d, meta)
}

func resourcePoolV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	elbClient, err := cfg.ElbV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	pool, err := pools.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "pool")
	}

	log.Printf("[DEBUG] Retrieved pool %s: %#v", d.Id(), pool)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("lb_method", pool.LBMethod),
		d.Set("protocol", pool.Protocol),
		d.Set("description", pool.Description),
		d.Set("name", pool.Name),
		d.Set("type", pool.Type),
		d.Set("vpc_id", pool.VpcId),
		d.Set("protection_status", pool.ProtectionStatus),
		d.Set("protection_reason", pool.ProtectionReason),
		d.Set("slow_start_enabled", pool.SlowStart.Enable),
		d.Set("slow_start_duration", pool.SlowStart.Duration),
		d.Set("ip_version", pool.IpVersion),
	)

	if len(pool.Loadbalancers) != 0 {
		mErr = multierror.Append(mErr, d.Set("loadbalancer_id", pool.Loadbalancers[0].ID))
	}

	if len(pool.Listeners) != 0 {
		mErr = multierror.Append(mErr, d.Set("listener_id", pool.Listeners[0].ID))
	}

	if pool.Persistence.Type != "" {
		var persistence = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["cookie_name"] = pool.Persistence.CookieName
		params["type"] = pool.Persistence.Type
		params["timeout"] = pool.Persistence.PersistenceTimeout
		persistence[0] = params
		mErr = multierror.Append(mErr, d.Set("persistence", persistence))
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Dedicated ELB pool fields: %s", err)
	}

	return nil
}

func resourcePoolV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	var updateOpts pools.UpdateOpts
	if d.HasChange("lb_method") {
		updateOpts.LBMethod = d.Get("lb_method").(string)
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("persistence") {
		persistence, err := buildPersistence(d.Get("persistence"))
		if err != nil {
			return diag.FromErr(err)
		}
		updateOpts.Persistence = &persistence
	}
	if d.HasChange("protection_status") {
		updateOpts.ProtectionStatus = d.Get("protection_status").(string)
	}
	if d.HasChange("protection_reason") {
		protectionReason := d.Get("protection_reason").(string)
		updateOpts.ProtectionReason = &protectionReason
	}
	if d.HasChange("type") {
		updateOpts.Type = d.Get("type").(string)
	}
	if d.HasChange("vpc_id") {
		updateOpts.VpcId = d.Get("vpc_id").(string)
	}
	if d.HasChanges("slow_start_enabled", "slow_start_duration") {
		updateOpts.SlowStart = &pools.SlowStart{
			Enable:   d.Get("slow_start_enabled").(bool),
			Duration: d.Get("slow_start_duration").(int),
		}
	}

	log.Printf("[DEBUG] Updating pool %s with options: %#v", d.Id(), updateOpts)
	_, err = pools.Update(elbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to update pool %s: %s", d.Id(), err)
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForElbV3Pool(ctx, elbClient, d.Id(), "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePoolV3Read(ctx, d, meta)
}

func resourcePoolV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	log.Printf("[DEBUG] Attempting to delete pool %s", d.Id())
	err = pools.Delete(elbClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("unable to delete pool %s: %s", d.Id(), err)
	}

	// Wait for Pool to delete
	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForElbV3Pool(ctx, elbClient, d.Id(), "DELETED", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForElbV3Pool(ctx context.Context, elbClient *golangsdk.ServiceClient, id string, target string, pending []string,
	timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for pool %s to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceElbV3PoolRefreshFunc(elbClient, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("error: pool %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("error waiting for pool %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceElbV3PoolRefreshFunc(elbClient *golangsdk.ServiceClient, poolID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pool, err := pools.Get(elbClient, poolID).Extract()
		if err != nil {
			return nil, "", err
		}

		// The pool resource has no Status attribute, so a successful Get is the best we can do
		return pool, "ACTIVE", nil
	}
}

func buildPersistence(p interface{}) (pools.SessionPersistence, error) {
	pV := (p.([]interface{}))[0].(map[string]interface{})

	persistence := pools.SessionPersistence{
		Type: pV["type"].(string),
	}

	if persistence.Type == "APP_COOKIE" {
		if pV["cookie_name"].(string) == "" {
			return persistence, fmt.Errorf(
				"persistence cookie_name needs to be set if using 'APP_COOKIE' persistence type")
		}
		persistence.CookieName = pV["cookie_name"].(string)

		if pV["timeout"].(int) != 0 {
			return persistence, fmt.Errorf(
				"persistence timeout is invalid when type is set to 'APP_COOKIE'")
		}
	} else if pV["cookie_name"].(string) != "" {
		return persistence, fmt.Errorf(
			"persistence cookie_name can only be set if using 'APP_COOKIE' persistence type")
	}

	persistence.PersistenceTimeout = pV["timeout"].(int)
	return persistence, nil
}
