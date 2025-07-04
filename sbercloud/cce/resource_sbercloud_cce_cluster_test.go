package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/chnsz/golangsdk/openstack/cce/v3/clusters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCCEClusterV3_basic(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "Available"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "cce.s1.small"),
					resource.TestCheckResourceAttr(resourceName, "container_network_type", "overlay_l2"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCCEClusterV3_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "new description"),
				),
			},
		},
	})
}

func TestAccCCEClusterV3_withEip(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_withEip(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode", "rbac"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"eip",
				},
			},
		},
	})
}

func TestAccCCEClusterV3_withEpsId(t *testing.T) {
	var cluster clusters.Clusters

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterV3_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV3Exists(resourceName, &cluster),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.SBC_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

func testAccCheckCCEClusterV3Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	cceClient, err := config.CceV3Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud CCE client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_cce_cluster" {
			continue
		}

		_, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Cluster still exists")
		}
	}

	return nil
}

func testAccCheckCCEClusterV3Exists(n string, cluster *clusters.Clusters) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		cceClient, err := config.CceV3Client(acceptance.SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud CCE client: %s", err)
		}

		found, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("Cluster not found")
		}

		*cluster = *found

		return nil
	}
}

func testAccCCEClusterV3_Base(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_subnet" "test" {
  name          = "%s"
  cidr          = "192.168.0.0/16"
  gateway_ip    = "192.168.0.1"

  //dns is required for cce node installing
  primary_dns   = "100.125.13.59"
  secondary_dns = "8.8.8.8"
  vpc_id        = sbercloud_vpc.test.id
}
`, rName, rName)
}

func testAccCCEClusterV3_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = sbercloud_vpc.test.id
  subnet_id              = sbercloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
}
`, testAccCCEClusterV3_Base(rName), rName)
}

func testAccCCEClusterV3_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = sbercloud_vpc.test.id
  subnet_id              = sbercloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  description            = "new description"
}
`, testAccCCEClusterV3_Base(rName), rName)
}

func testAccCCEClusterV3_withEip(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "sbercloud_cce_cluster" "test" {
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = sbercloud_vpc.test.id
  subnet_id              = sbercloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  authentication_mode    = "rbac"
  eip                    = sbercloud_vpc_eip.test.address
}
`, testAccCCEClusterV3_Base(rName), rName)
}

func testAccCCEClusterV3_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = sbercloud_vpc.test.id
  subnet_id              = sbercloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  enterprise_project_id  = "%s"
}

`, testAccCCEClusterV3_Base(rName), rName, acceptance.SBC_ENTERPRISE_PROJECT_ID)
}

func testAccCluster_turbo(rName string, eniNum int) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_vpc_subnet" "eni_test" {
  count      = %[3]d

  name       = "%[2]s-eni-${count.index}"
  cidr       = cidrsubnet(sbercloud_vpc.test.cidr, 8, count.index + 1)
  gateway_ip = cidrhost(cidrsubnet(sbercloud_vpc.test.cidr, 8, count.index + 1), 1)
  vpc_id     = sbercloud_vpc.test.id
}

resource "sbercloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = sbercloud_vpc.test.id
  subnet_id              = sbercloud_vpc_subnet.test.id
  container_network_type = "eni"
  enable_dist_mgt        = true
  eni_subnet_id          = join(",", sbercloud_vpc_subnet.eni_test[*].ipv4_subnet_id)
}

output "is_eni_subnet_id_different" {
  value = length(setsubtract(split(",", sbercloud_cce_cluster.test.eni_subnet_id),
  sbercloud_vpc_subnet.eni_test[*].ipv4_subnet_id)) != 0
}
`, acceptance.TestVpc(rName), rName, eniNum)
}
