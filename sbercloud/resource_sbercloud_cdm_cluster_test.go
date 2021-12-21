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

func TestAccCdmClusterV1_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdmClusterV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdmClusterV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdmClusterV1Exists(),
				),
			},
		},
	})
}

func testAccCdmClusterV1_basic(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_availability_zones" "test" {}

data "sbercloud_vpc" "test" {
 name = "vpc-default"
}

data "sbercloud_vpc_subnet" "test" {
 name = "subnet-default"
}

data "sbercloud_cdm_flavors" "test" {}

resource "sbercloud_networking_secgroup" "secgroup" {
 name        = "%s"
 description = "terraform security group acceptance test"
}

resource "sbercloud_cdm_cluster" "cluster" {
 availability_zone = data.sbercloud_availability_zones.test.names[0]
 flavor_id         = data.sbercloud_cdm_flavors.test.flavors[0].id
 name              = "%s"
 security_group_id = sbercloud_networking_secgroup.secgroup.id
 subnet_id         = data.sbercloud_vpc_subnet.test.id
 vpc_id            = data.sbercloud_vpc.test.id
 version           = "2.8.2"
}`, rName, rName)
}

func testAccCheckCdmClusterV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.CdmV11Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_cdm_cluster" {
			continue
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
				"X-Language":   "en-us",
			}})
		if err == nil {
			return fmt.Errorf("sbercloud_cdm_cluster still exists at %s", url)
		}
	}

	return nil
}

func testAccCheckCdmClusterV1Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.CdmV11Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["sbercloud_cdm_cluster.cluster"]
		if !ok {
			return fmt.Errorf("Error checking sbercloud_cdm_cluster.cluster exist, err=not found this resource")
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return fmt.Errorf("Error checking sbercloud_cdm_cluster.cluster exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
				"X-Language":   "en-us",
			}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmt.Errorf("sbercloud_cdm_cluster.cluster is not exist")
			}
			return fmt.Errorf("Error checking sbercloud_cdm_cluster.cluster exist, err=send request failed: %s", err)
		}
		return nil
	}
}
