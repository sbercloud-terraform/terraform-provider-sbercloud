package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccGesGraphV1_basic(t *testing.T) {
	name := fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGesGraphV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGesGraphV1_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGesGraphV1Exists(),
				),
			},
		},
	})
}

func testAccGesGraphV1_basic(name string) string {
	return fmt.Sprintf(`
data "sbercloud_availability_zones" "test" {}

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

resource "sbercloud_networking_secgroup" "test" {
  name = "%s"
}

resource "sbercloud_ges_graph" "graph" {
  availability_zone = data.sbercloud_availability_zones.test.names[1]
  graph_size_type   = 0
  name              = "%s"
  region            = "%s"
  security_group_id = sbercloud_networking_secgroup.test.id
  subnet_id         = sbercloud_vpc_subnet.test.id
  vpc_id            = sbercloud_vpc.test.id
}
	`, name, name, name, name, SBC_REGION_NAME)
}

func testAccCheckGesGraphV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.GesV1Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_ges_graph" {
			continue
		}

		url, err := replaceVarsForTest(rs, "graphs/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err == nil {
			return fmt.Errorf("sbercloud_ges_graph still exists at %s", url)
		}
	}

	return nil
}

func testAccCheckGesGraphV1Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.GesV1Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["sbercloud_ges_graph.graph"]
		if !ok {
			return fmt.Errorf("Error checking sbercloud_ges_graph.graph exist, err=not found this resource")
		}

		url, err := replaceVarsForTest(rs, "graphs/{id}")
		if err != nil {
			return fmt.Errorf("Error checking sbercloud_ges_graph.graph exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmt.Errorf("sbercloud_ges_graph.graph is not exist")
			}
			return fmt.Errorf("Error checking sbercloud_ges_graph.graph exist, err=send request failed: %s", err)
		}
		return nil
	}
}
