package huaweicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/natgateways"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceNatGatewayV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNatGatewayV2Create,
		Read:   resourceNatGatewayV2Read,
		Update: resourceNatGatewayV2Update,
		Delete: resourceNatGatewayV2Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"spec": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"1", "2", "3", "4",
				}, false),
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ExactlyOneOf: []string{"router_id"},
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ExactlyOneOf: []string{"internal_network_id"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// deprecated
			"internal_network_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "use subnet_id instead",
			},
			"router_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "use vpc_id instead",
			},
		},
	}
}

func resourceNatGatewayV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	natClient, err := config.NatGatewayClient(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud nat client: %s", err)
	}

	var vpcID, subnetID string
	if v1, ok := d.GetOk("vpc_id"); ok {
		vpcID = v1.(string)
	} else {
		vpcID = d.Get("router_id").(string)
	}
	if v2, ok := d.GetOk("subnet_id"); ok {
		subnetID = v2.(string)
	} else {
		subnetID = d.Get("internal_network_id").(string)
	}

	createOpts := &natgateways.CreateOpts{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Spec:                d.Get("spec").(string),
		RouterID:            vpcID,
		InternalNetworkID:   subnetID,
		EnterpriseProjectID: GetEnterpriseProjectID(d, config),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	natGateway, err := natgateways.Create(natClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creatting Nat Gateway: %s", err)
	}

	logp.Printf("[DEBUG] Waiting for HuaweiCloud Nat Gateway (%s) to become available.", natGateway.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitForNatGatewayActive(natClient, natGateway.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud Nat Gateway: %s", err)
	}

	d.SetId(natGateway.ID)

	return resourceNatGatewayV2Read(d, meta)
}

func resourceNatGatewayV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	natClient, err := config.NatGatewayClient(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud nat client: %s", err)
	}

	natGateway, err := natgateways.Get(natClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Nat Gateway")
	}

	d.Set("name", natGateway.Name)
	d.Set("description", natGateway.Description)
	d.Set("spec", natGateway.Spec)
	d.Set("vpc_id", natGateway.RouterID)
	d.Set("subnet_id", natGateway.InternalNetworkID)
	d.Set("status", natGateway.Status)
	d.Set("region", GetRegion(d, config))
	d.Set("enterprise_project_id", natGateway.EnterpriseProjectID)

	return nil
}

func resourceNatGatewayV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	natClient, err := config.NatGatewayClient(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud nat client: %s", err)
	}

	var updateOpts natgateways.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("spec") {
		updateOpts.Spec = d.Get("spec").(string)
	}

	logp.Printf("[DEBUG] Update Options: %#v", updateOpts)

	_, err = natgateways.Update(natClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating Nat Gateway: %s", err)
	}

	return resourceNatGatewayV2Read(d, meta)
}

func resourceNatGatewayV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	natClient, err := config.NatGatewayClient(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud nat client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForNatGatewayDelete(natClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud Nat Gateway: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForNatGatewayActive(client *golangsdk.ServiceClient, nId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := natgateways.Get(client, nId).Extract()
		if err != nil {
			return nil, "", err
		}

		logp.Printf("[DEBUG] HuaweiCloud Nat Gateway: %+v", n)
		if n.Status == "ACTIVE" {
			return n, "ACTIVE", nil
		}

		return n, "", nil
	}
}

func waitForNatGatewayDelete(client *golangsdk.ServiceClient, nId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[DEBUG] Attempting to delete HuaweiCloud Nat Gateway %s.\n", nId)

		n, err := natgateways.Get(client, nId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Nat gateway %s", nId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		err = natgateways.Delete(client, nId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Nat Gateway %s", nId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		logp.Printf("[DEBUG] HuaweiCloud Nat Gateway %s still active.\n", nId)
		return n, "ACTIVE", nil
	}
}
