package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDisPartitionV2DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDisStreamV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDisPartitionV2_basic(acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDisPartitionV2Exists(),
				),
			},
		},
	})
}

func testAccCheckDisPartitionV2Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["data.sbercloud_dis_partition.partition"]
		if !ok {
			return fmt.Errorf("Error checking sbercloud_dis_partition.partition exist, err=not found this resource")
		}

		if _, ok := rs.Primary.Attributes["partitions.0.id"]; !ok {
			return fmt.Errorf("expect partitions to be set")
		}

		return nil
	}
}

func testAccDisPartitionV2_basic(random string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_dis_partition" "partition" {
  stream_name = "${sbercloud_dis_stream.stream.stream_name}"
}
`, testAccDisStreamV2_basic(random))
}
