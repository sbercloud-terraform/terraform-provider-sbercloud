package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSFSFileSystemV2DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSFileSystemV2DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSFileSystemV2DataSourceID("data.sbercloud_sfs_file_system.shares"),
					resource.TestCheckResourceAttr("data.sbercloud_sfs_file_system.shares", "name", rName),
					resource.TestCheckResourceAttr("data.sbercloud_sfs_file_system.shares", "status", "available"),
					resource.TestCheckResourceAttr("data.sbercloud_sfs_file_system.shares", "size", "1"),
				),
			},
		},
	})
}

func testAccCheckSFSFileSystemV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find share file data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("share file data source ID not set ")
		}

		return nil
	}
}

func testAccSFSFileSystemV2DataSource_basic(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc" "vpc_default" {
  name = "vpc-default"
  enterprise_project_id = "0"
}

data "sbercloud_availability_zones" "myaz" {}

resource "sbercloud_sfs_file_system" "sfs_1" {
	share_proto = "NFS"
	size=1
	name="%s"
	availability_zone = data.sbercloud_availability_zones.myaz.names[0]
	access_to = data.sbercloud_vpc.vpc_default.id
  	access_type="cert"
  	access_level="rw"
	description="sfs_c2c_test-file"
}
data "sbercloud_sfs_file_system" "shares" {
  id = sbercloud_sfs_file_system.sfs_1.id
}
`, rName)
}
