package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/vpnaas/siteconnections"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccVpnSiteConnectionV2_basic(t *testing.T) {
	var conn siteconnections.Connection
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_vpnaas_site_connection.conn_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteConnectionV2_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteConnectionV2Exists(resourceName, &conn),
					resource.TestCheckResourceAttrPtr(resourceName, "name", &conn.Name),
					resource.TestCheckResourceAttrPtr(resourceName, "vpnservice_id", &conn.VPNServiceID),
					resource.TestCheckResourceAttrPtr(resourceName, "ikepolicy_id", &conn.IKEPolicyID),
					resource.TestCheckResourceAttrPtr(resourceName, "ipsecpolicy_id", &conn.IPSecPolicyID),
					resource.TestCheckResourceAttrPtr(resourceName, "peer_ep_group_id", &conn.PeerEPGroupID),
					resource.TestCheckResourceAttrPtr(resourceName, "local_ep_group_id", &conn.LocalEPGroupID),
					resource.TestCheckResourceAttrPtr(resourceName, "local_id", &conn.LocalID),

					resource.TestCheckResourceAttr(resourceName, "admin_state_up", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
		},
	})
}

func testAccCheckSiteConnectionV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_vpnaas_site_connection" {
			continue
		}
		_, err = siteconnections.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Site connection (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckSiteConnectionV2Exists(n string, conn *siteconnections.Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud networking client: %s", err)
		}

		var found *siteconnections.Connection

		found, err = siteconnections.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*conn = *found

		return nil
	}
}

func testAccSiteConnectionV2_basic(name string) string {
	return fmt.Sprintf(`
	resource "sbercloud_vpc" "test" {
		name = "%s"
		cidr = "192.168.0.0/16"
}

	resource "sbercloud_vpc_subnet" "test" {
		name          = "%s"
		cidr          = "192.168.0.0/24"
		gateway_ip    = "192.168.0.1"
		primary_dns   = "100.125.1.250"
		secondary_dns = "100.125.21.250"
		vpc_id        = sbercloud_vpc.test.id
	}

	resource "sbercloud_vpnaas_service" "service_1" {
		name = "%s"
		router_id = sbercloud_vpc.test.id
	}

	resource "sbercloud_vpnaas_ipsec_policy" "policy_1" {
	}

	resource "sbercloud_vpnaas_ike_policy" "policy_2" {
	}

	resource "sbercloud_vpnaas_endpoint_group" "group_1" {
		type = "cidr"
		endpoints = ["10.2.0.0/24", "10.3.0.0/24"]
	}

	resource "sbercloud_vpnaas_endpoint_group" "group_2" {
		type = "cidr"
		endpoints = [sbercloud_vpc_subnet.test.cidr]
	}

	resource "sbercloud_vpnaas_site_connection" "conn_1" {
		name = "%s"
		ikepolicy_id = sbercloud_vpnaas_ike_policy.policy_2.id
		ipsecpolicy_id = sbercloud_vpnaas_ipsec_policy.policy_1.id
		vpnservice_id = sbercloud_vpnaas_service.service_1.id
		psk = "secret"
		peer_address = "192.168.10.1"
		peer_id = "192.168.10.1"
		local_ep_group_id = sbercloud_vpnaas_endpoint_group.group_2.id
		peer_ep_group_id = sbercloud_vpnaas_endpoint_group.group_1.id

		tags = {
			foo = "bar"
			key = "value"
		}

		depends_on = ["sbercloud_vpc_subnet.test"]
	}
	`, name, name, name, name)
}
