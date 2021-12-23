package vpc

import (
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/routes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceVPCRouteV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcRouteV2Create,
		Read:   resourceVpcRouteV2Read,
		Delete: resourceVpcRouteV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ //request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nexthop": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: utils.ValidateCIDR,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVpcRouteV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	vpcRouteClient, err := config.NetworkingV2Client(config.GetRegion(d))

	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud vpc route client: %s", err)
	}

	createOpts := routes.CreateOpts{
		Type:        d.Get("type").(string),
		NextHop:     d.Get("nexthop").(string),
		Destination: d.Get("destination").(string),
		VPC_ID:      d.Get("vpc_id").(string),
	}

	n, err := routes.Create(vpcRouteClient, createOpts).Extract()

	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud VPC route: %s", err)
	}
	d.SetId(n.RouteID)

	logp.Printf("[INFO] Vpc Route ID: %s", n.RouteID)

	d.SetId(n.RouteID)

	return resourceVpcRouteV2Read(d, meta)

}

func resourceVpcRouteV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	vpcRouteClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud Vpc route client: %s", err)
	}

	n, err := routes.Get(vpcRouteClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmtp.Errorf("Error retrieving Huaweicloud Vpc route: %s", err)
	}

	d.Set("type", n.Type)
	d.Set("nexthop", n.NextHop)
	d.Set("destination", n.Destination)
	d.Set("vpc_id", n.VPC_ID)
	d.Set("region", config.GetRegion(d))

	return nil
}

func resourceVpcRouteV2Delete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*config.Config)
	vpcRouteClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud vpc route: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcRouteDelete(vpcRouteClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error deleting Huaweicloud Vpc route: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForVpcRouteDelete(vpcRouteClient *golangsdk.ServiceClient, routeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := routes.Get(vpcRouteClient, routeId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[INFO] Successfully deleted Huaweicloud vpc route %s", routeId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = routes.Delete(vpcRouteClient, routeId).ExtractErr()
		logp.Printf("[DEBUG] Value if error: %#v", err)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[INFO] Successfully deleted Huaweicloud vpc route %s", routeId)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "ACTIVE", nil
				}
			}
			return r, "ACTIVE", err
		}

		return r, "ACTIVE", nil
	}
}
