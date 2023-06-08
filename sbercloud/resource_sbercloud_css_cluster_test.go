package sbercloud

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCssClusterV1_basic(t *testing.T) {
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_css_cluster.cluster"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCssClusterV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCssClusterV1_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssClusterV1Exists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s", name)),
					resource.TestCheckResourceAttr(resourceName, "expect_node_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "elasticsearch"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccCssClusterV1_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssClusterV1Exists(),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value"),
				),
			},
		},
	})
}

func TestAccCssClusterV1_security(t *testing.T) {
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_css_cluster.cluster"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCssClusterV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCssClusterV1_security(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssClusterV1Exists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s", name)),
					resource.TestCheckResourceAttr(resourceName, "expect_node_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "elasticsearch"),
					resource.TestCheckResourceAttr(resourceName, "security_mode", "true"),
				),
			},
		},
	})
}

func testAccCheckCssClusterV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.CssV1Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_css_cluster" {
			continue
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err == nil {
			return fmt.Errorf("sbercloud_css_cluster still exists at %s", url)
		}
	}

	return nil
}

func testAccCheckCssClusterV1Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.CssV1Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["sbercloud_css_cluster.cluster"]
		if !ok {
			return fmt.Errorf("Error checking sbercloud_css_cluster.cluster exist, err=not found this resource")
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return fmt.Errorf("Error checking sbercloud_css_cluster.cluster exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmt.Errorf("sbercloud_css_cluster.cluster is not exist")
			}
			return fmt.Errorf("Error checking sbercloud_css_cluster.cluster exist, err=send request failed: %s", err)
		}
		return nil
	}
}

func testAccCssClusterV1_base(name string) string {
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
}`, name, name, name)
}

func testAccCssClusterV1_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_css_cluster" "cluster" {
  name = "%s"
  engine_version  = "7.10.2"
  expect_node_num = 1

  node_config {
    flavor = "ess.spec-4u8g"
    network_info {
      security_group_id = sbercloud_networking_secgroup.test.id
      subnet_id = sbercloud_vpc_subnet.test.id
      vpc_id = sbercloud_vpc.test.id
    }
    volume {
      volume_type = "HIGH"
      size = 40
    }
    availability_zone = data.sbercloud_availability_zones.test.names[0]
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
	`, testAccCssClusterV1_base(name), name)
}

func testAccCssClusterV1_update(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_css_cluster" "cluster" {
  name = "%s"
  engine_version  = "7.10.2"
  expect_node_num = 1

  node_config {
    flavor = "ess.spec-4u8g"
    network_info {
      security_group_id = sbercloud_networking_secgroup.test.id
      subnet_id = sbercloud_vpc_subnet.test.id
      vpc_id = sbercloud_vpc.test.id
    }
    volume {
      volume_type = "HIGH"
      size = 40
    }
    availability_zone = data.sbercloud_availability_zones.test.names[0]
  }
  tags = {
    foo = "bar_update"
    key_update = "value"
  }
}
	`, testAccCssClusterV1_base(name), name)
}

func testAccCssClusterV1_security(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_css_cluster" "cluster" {
  name = "%s"
  engine_version  = "7.6.2"
  expect_node_num = 1
  security_mode   = true
  password        = "Test@passw0rd"

  node_config {
    flavor = "ess.spec-4u8g"
    network_info {
      security_group_id = sbercloud_networking_secgroup.test.id
      subnet_id = sbercloud_vpc_subnet.test.id
      vpc_id = sbercloud_vpc.test.id
    }
    volume {
      volume_type = "HIGH"
      size = 40
    }
    availability_zone = data.sbercloud_availability_zones.test.names[0]
  }
}
	`, testAccCssClusterV1_base(name), name)
}
