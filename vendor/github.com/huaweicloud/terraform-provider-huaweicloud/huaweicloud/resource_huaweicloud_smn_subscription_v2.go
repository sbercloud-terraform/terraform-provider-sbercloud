package huaweicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/smn/v2/subscriptions"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubscriptionCreate,
		Read:   resourceSubscriptionRead,
		Delete: resourceSubscriptionDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"email", "sms", "http", "https",
				}, false),
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subscription_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.SmnV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud smn client: %s", err)
	}
	topicUrn := d.Get("topic_urn").(string)
	createOpts := subscriptions.CreateOps{
		Endpoint: d.Get("endpoint").(string),
		Protocol: d.Get("protocol").(string),
		Remark:   d.Get("remark").(string),
	}
	logp.Printf("[DEBUG] Create Options: %#v", createOpts)

	subscription, err := subscriptions.Create(client, createOpts, topicUrn).Extract()
	if err != nil {
		return fmtp.Errorf("Error getting subscription from result: %s", err)
	}
	logp.Printf("[DEBUG] Create : subscription.SubscriptionUrn %s", subscription.SubscriptionUrn)
	if subscription.SubscriptionUrn != "" {
		d.SetId(subscription.SubscriptionUrn)
		d.Set("subscription_urn", subscription.SubscriptionUrn)
		return resourceSubscriptionRead(d, meta)
	}

	return fmtp.Errorf("Unexpected conversion error in resourceSubscriptionCreate.")
}

func resourceSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.SmnV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud smn client: %s", err)
	}

	logp.Printf("[DEBUG] Deleting subscription %s", d.Id())

	id := d.Id()
	result := subscriptions.Delete(client, id)
	if result.Err != nil {
		return result.Err
	}

	logp.Printf("[DEBUG] Successfully deleted subscription %s", id)
	return nil
}

func resourceSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.SmnV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud smn client: %s", err)
	}

	logp.Printf("[DEBUG] Getting subscription %s", d.Id())

	id := d.Id()
	subscriptionslist, err := subscriptions.List(client).Extract()
	if err != nil {
		return fmtp.Errorf("Error Get subscriptionslist: %s", err)
	}
	logp.Printf("[DEBUG] list : subscriptionslist %#v", subscriptionslist)
	for _, subscription := range subscriptionslist {
		if subscription.SubscriptionUrn == id {
			logp.Printf("[DEBUG] subscription: %#v", subscription)
			d.Set("topic_urn", subscription.TopicUrn)
			d.Set("endpoint", subscription.Endpoint)
			d.Set("protocol", subscription.Protocol)
			d.Set("subscription_urn", subscription.SubscriptionUrn)
			d.Set("owner", subscription.Owner)
			d.Set("remark", subscription.Remark)
			d.Set("status", subscription.Status)
		}
	}

	logp.Printf("[DEBUG] Successfully get subscription %s", id)
	return nil
}
