package sbercloud

import (
	"fmt"
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/bss/v2/orders"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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

// CheckDeleted checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeleted(d *schema.ResourceData, err error, msg string) error {
	if _, ok := err.(golangsdk.ErrDefault404); ok {
		d.SetId("")
		return nil
	}

	return fmt.Errorf("%s: %s", msg, err)
}

func TestBaseNetwork(name string) string {
	return fmt.Sprintf(`
# base security group without default rules
%s

# base vpc and subnet
%s
`, TestSecGroup(name), TestVpc(name))
}

// TestSecGroup can be referred as `sbercloud_networking_secgroup.test`
func TestSecGroup(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_networking_secgroup" "test" {
  name                 = "%s"
  delete_default_rules = true
}
`, name)
}

// TestVpc can be referred as `sbercloud_vpc.test` and `sbercloud_vpc_subnet.test`
func TestVpc(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_subnet" "test" {
  name       = "%[1]s"
  vpc_id     = sbercloud_vpc.test.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}
`, name)
}
