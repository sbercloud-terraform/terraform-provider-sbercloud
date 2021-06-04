package sbercloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/bss/v2/orders"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// GetRegion returns the region that was specified in the resource. If a
// region was not set, the provider-level region is checked. The provider-level
// region can either be set by the region argument or by SBC_REGION_NAME.
func GetRegion(d *schema.ResourceData, config *config.Config) string {
	if v, ok := d.GetOk("region"); ok {
		return v.(string)
	}

	return config.Region
}

// GetEnterpriseProjectID returns the enterprise_project_id that was specified in the resource.
// If it was not set, the provider-level value is checked. The provider-level value can
// either be set by the `enterprise_project_id` argument or by SBC_ENTERPRISE_PROJECT_ID.
func GetEnterpriseProjectID(d *schema.ResourceData, config *config.Config) string {
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		return v.(string)
	}

	return config.EnterpriseProjectID
}

// UnsubscribePrePaidResource impl the action of unsubscribe resource
func UnsubscribePrePaidResource(d *schema.ResourceData, config *config.Config, resourceIDs []string) error {
	bssV2Client, err := config.BssV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating SberCloud bss V2 client: %s", err)
	}

	unsubscribeOpts := orders.UnsubscribeOpts{
		ResourceIds:     resourceIDs,
		UnsubscribeType: 1,
	}
	_, err = orders.Unsubscribe(bssV2Client, unsubscribeOpts).Extract()
	return err
}
