package huaweicloud

import (
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/vpcep/v1/services"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const (
	actionReceive string = "receive"
	actionReject  string = "reject"
)

var approvalActionStatusMap = map[string]string{
	actionReceive: "accepted",
	actionReject:  "rejected",
}

func ResourceVPCEndpointApproval() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPCEndpointApprovalCreate,
		Read:   resourceVPCEndpointApprovalRead,
		Update: resourceVPCEndpointApprovalUpdate,
		Delete: resourceVPCEndpointApprovalDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoints": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"packet_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceVPCEndpointApprovalCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	vpcepClient, err := config.VPCEPClient(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud VPC endpoint client: %s", err)
	}

	// check status of the VPC endpoint service
	serviceID := d.Get("service_id").(string)
	n, err := services.Get(vpcepClient, serviceID).Extract()
	if err != nil {
		return fmtp.Errorf("Error retrieving VPC endpoint service %s: %s", serviceID, err)
	}
	if n.Status != "available" {
		return fmtp.Errorf("Error the status of VPC endpoint service is %s, expected to be available", n.Status)
	}

	raw := d.Get("endpoints").(*schema.Set).List()
	err = doConnectionAction(d, vpcepClient, serviceID, actionReceive, raw)
	if err != nil {
		return fmtp.Errorf("Error receiving connections to VPC endpoint service %s: %s", serviceID, err)
	}

	d.SetId(serviceID)
	return resourceVPCEndpointApprovalRead(d, meta)
}

func resourceVPCEndpointApprovalRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	vpcepClient, err := config.VPCEPClient(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud VPC endpoint client: %s", err)
	}

	serviceID := d.Get("service_id").(string)
	if conns, err := flattenVPCEndpointConnections(vpcepClient, serviceID); err == nil {
		d.Set("connections", conns)
	}

	return nil
}

func resourceVPCEndpointApprovalUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	vpcepClient, err := config.VPCEPClient(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud VPC endpoint client: %s", err)
	}

	if d.HasChange("endpoints") {
		old, new := d.GetChange("endpoints")
		oldConnSet := old.(*schema.Set)
		newConnSet := new.(*schema.Set)
		received := newConnSet.Difference(oldConnSet)
		rejected := oldConnSet.Difference(newConnSet)

		serviceID := d.Get("service_id").(string)
		err = doConnectionAction(d, vpcepClient, serviceID, actionReceive, received.List())
		if err != nil {
			return fmtp.Errorf("Error receiving connections to VPC endpoint service %s: %s", serviceID, err)
		}

		err = doConnectionAction(d, vpcepClient, serviceID, actionReject, rejected.List())
		if err != nil {
			return fmtp.Errorf("Error rejecting connections to VPC endpoint service %s: %s", serviceID, err)
		}
	}
	return resourceVPCEndpointApprovalRead(d, meta)
}

func resourceVPCEndpointApprovalDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	vpcepClient, err := config.VPCEPClient(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud VPC endpoint client: %s", err)
	}

	serviceID := d.Get("service_id").(string)
	raw := d.Get("endpoints").(*schema.Set).List()
	err = doConnectionAction(d, vpcepClient, serviceID, actionReject, raw)
	if err != nil {
		return fmtp.Errorf("Error rejecting connections to VPC endpoint service %s: %s", serviceID, err)
	}

	d.SetId("")
	return nil
}

func doConnectionAction(d *schema.ResourceData, client *golangsdk.ServiceClient, serviceID, action string, raw []interface{}) error {
	if len(raw) == 0 {
		return nil
	}

	if _, ok := approvalActionStatusMap[action]; !ok {
		return fmtp.Errorf("approval action(%s) is invalid, only support %s or %s", action, actionReceive, actionReject)
	}

	targetStatus := approvalActionStatusMap[action]
	for _, v := range raw {
		// Each request accepts or rejects only one VPC endpoint
		epID := v.(string)
		connOpts := services.ConnActionOpts{
			Action:    action,
			Endpoints: []string{epID},
		}

		logp.Printf("[DEBUG] %s to endpoint %s from VPC endpoint service %s", action, epID, serviceID)
		if result := services.ConnAction(client, serviceID, connOpts); result.Err != nil {
			return result.Err
		}

		logp.Printf("[INFO] Waiting for VPC endpoint(%s) to become %s", epID, targetStatus)
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"creating", "pendingAcceptance"},
			Target:     []string{targetStatus},
			Refresh:    waitForVPCEndpointConnected(client, serviceID, epID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      3 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, stateErr := stateConf.WaitForState()
		if stateErr != nil {
			return fmtp.Errorf(
				"Error waiting for VPC endpoint(%s) to become %s: %s",
				epID, targetStatus, stateErr)
		}
	}

	return nil
}

func waitForVPCEndpointConnected(vpcepClient *golangsdk.ServiceClient, serviceId, endpointId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		listOpts := services.ListConnOpts{
			EndpointID: endpointId,
		}
		connections, err := services.ListConnections(vpcepClient, serviceId, listOpts)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return connections, "deleted", nil
			}
			return connections, "error", err
		}
		if len(connections) == 1 && connections[0].EndpointID == endpointId {
			return connections, connections[0].Status, nil
		}
		return connections, "deleted", nil
	}
}
