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

func TestAccDwsCluster_basic(t *testing.T) {
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDwsClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDwsCluster_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDwsClusterExists(),
				),
			},
		},
	})
}

func testAccDwsCluster_basic(name string) string {
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

resource "sbercloud_networking_secgroup" "secgroup" {
  name = "%s"
  description = "terraform security group"
}

resource "sbercloud_dws_cluster" "cluster" {
  node_type = "dws2.m6.4xlarge.8"
  number_of_node = 3
  network_id = sbercloud_vpc_subnet.test.id
  vpc_id = sbercloud_vpc.test.id
  security_group_id = sbercloud_networking_secgroup.secgroup.id
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  name = "%s"
  user_name = "test_cluster_admin"
  user_pwd = "cluster123@!"

  timeouts {
    create = "30m"
    delete = "30m"
  }
}
	`, name, name, name, name)
}

func testAccCheckDwsClusterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.DwsV1Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_dws_cluster" {
			continue
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		_, err = client.Get(
			url, nil,
			&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err == nil {
			return fmt.Errorf("sbercloud_dws_cluster still exists at %s", url)
		}
	}

	return nil
}

func testAccCheckDwsClusterExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.DwsV1Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["sbercloud_dws_cluster.cluster"]
		if !ok {
			return fmt.Errorf("Error checking sbercloud_dws_cluster.cluster exist, err=not found sbercloud_dws_cluster.cluster")
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return fmt.Errorf("Error checking sbercloud_dws_cluster.cluster exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(
			url, nil,
			&golangsdk.RequestOpts{MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmt.Errorf("sbercloud_dws_cluster.cluster is not exist")
			}
			return fmt.Errorf("Error checking sbercloud_dws_cluster.cluster exist, err=send request failed: %s", err)
		}
		return nil
	}
}
