package huaweicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/pools"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourcePoolV2() *schema.Resource {
	return &schema.Resource{
		Create: resourcePoolV2Create,
		Read:   resourcePoolV2Read,
		Update: resourcePoolV2Update,
		Delete: resourcePoolV2Delete,

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

			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "tenant_id is deprecated",
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP", "HTTP",
				}, false),
			},

			// One of loadbalancer_id or listener_id must be provided
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// One of loadbalancer_id or listener_id must be provided
			"listener_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"lb_method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP",
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
					},
				},
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourcePoolV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	var persistence pools.SessionPersistence
	if p, ok := d.GetOk("persistence"); ok {
		pV := (p.([]interface{}))[0].(map[string]interface{})

		persistence = pools.SessionPersistence{
			Type: pV["type"].(string),
		}

		if persistence.Type == "APP_COOKIE" {
			if pV["cookie_name"].(string) == "" {
				return fmtp.Errorf(
					"Persistence cookie_name needs to be set if using 'APP_COOKIE' persistence type")
			}
			persistence.CookieName = pV["cookie_name"].(string)
		} else {
			if pV["cookie_name"].(string) != "" {
				return fmtp.Errorf(
					"Persistence cookie_name can only be set if using 'APP_COOKIE' persistence type")
			}
		}
	}

	createOpts := pools.CreateOpts{
		TenantID:       d.Get("tenant_id").(string),
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Protocol:       pools.Protocol(d.Get("protocol").(string)),
		LoadbalancerID: d.Get("loadbalancer_id").(string),
		ListenerID:     d.Get("listener_id").(string),
		LBMethod:       pools.LBMethod(d.Get("lb_method").(string)),
		AdminStateUp:   &adminStateUp,
	}

	// Must omit if not set
	if persistence != (pools.SessionPersistence{}) {
		createOpts.Persistence = &persistence
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutCreate)
	lbID := createOpts.LoadbalancerID
	listenerID := createOpts.ListenerID
	if lbID != "" {
		err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
		if err != nil {
			return err
		}
	} else if listenerID != "" {
		// Wait for Listener to become active before continuing
		err = waitForLBV2Listener(lbClient, listenerID, "ACTIVE", nil, timeout)
		if err != nil {
			return err
		}
	}

	logp.Printf("[DEBUG] Attempting to create pool")
	var pool *pools.Pool
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		pool, err = pools.Create(lbClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmtp.Errorf("Error creating pool: %s", err)
	}

	// Wait for LoadBalancer to become active before continuing
	if lbID != "" {
		err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
	} else {
		// Pool exists by now so we can ask for lbID
		err = waitForLBV2viaPool(lbClient, pool.ID, "ACTIVE", timeout)
	}
	if err != nil {
		return err
	}

	d.SetId(pool.ID)

	return resourcePoolV2Read(d, meta)
}

func resourcePoolV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	pool, err := pools.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "pool")
	}

	logp.Printf("[DEBUG] Retrieved pool %s: %#v", d.Id(), pool)

	d.Set("lb_method", pool.LBMethod)
	d.Set("protocol", pool.Protocol)
	d.Set("description", pool.Description)
	d.Set("tenant_id", pool.TenantID)
	d.Set("admin_state_up", pool.AdminStateUp)
	d.Set("name", pool.Name)
	d.Set("region", GetRegion(d, config))

	if pool.Persistence.Type != "" {
		var persistence []map[string]interface{} = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["cookie_name"] = pool.Persistence.CookieName
		params["type"] = pool.Persistence.Type
		persistence[0] = params
		if err = d.Set("persistence", persistence); err != nil {
			return fmtp.Errorf("Load balance persistence set error: %s", err)
		}
	}

	return nil
}

func resourcePoolV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	var updateOpts pools.UpdateOpts
	if d.HasChange("lb_method") {
		updateOpts.LBMethod = pools.LBMethod(d.Get("lb_method").(string))
	}
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutUpdate)
	lbID := d.Get("loadbalancer_id").(string)
	if lbID != "" {
		err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
	} else {
		err = waitForLBV2viaPool(lbClient, d.Id(), "ACTIVE", timeout)
	}
	if err != nil {
		return err
	}

	logp.Printf("[DEBUG] Updating pool %s with options: %#v", d.Id(), updateOpts)
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err = pools.Update(lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmtp.Errorf("Unable to update pool %s: %s", d.Id(), err)
	}

	// Wait for LoadBalancer to become active before continuing
	if lbID != "" {
		err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
	} else {
		err = waitForLBV2viaPool(lbClient, d.Id(), "ACTIVE", timeout)
	}
	if err != nil {
		return err
	}

	return resourcePoolV2Read(d, meta)
}

func resourcePoolV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutDelete)
	lbID := d.Get("loadbalancer_id").(string)
	if lbID != "" {
		err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
		if err != nil {
			return err
		}
	}

	logp.Printf("[DEBUG] Attempting to delete pool %s", d.Id())
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = pools.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if lbID != "" {
		err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
	} else {
		// Wait for Pool to delete
		err = waitForLBV2Pool(lbClient, d.Id(), "DELETED", nil, timeout)
	}
	if err != nil {
		return err
	}

	return nil
}
