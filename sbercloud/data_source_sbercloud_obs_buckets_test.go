package sbercloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	//"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"
)

// В хуавей пречер осуществляется через методы в пакете acceptance, нужно переделпть под сберклауд
func TestAccDataSourceObsBuckets_basic(t *testing.T) {
	dataSourceName := "data.sbercloud_obs_buckets.buckets"
	//name := acceptance.RandomAccResourceNameWithDash()
	//dc := acceptance.InitDataSourceCheck(dataSourceName)

	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			//acceptance.TestAccPreCheck(t)
			//acceptance.TestAccPreCheckOBS(t)
			testAccPreCheckOBS(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBuckets_conf(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(name),
					resource.TestCheckResourceAttr(dataSourceName, "buckets.0.bucket", name),
				),
			},
		},
	})
}

func testAccObsBuckets_conf(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "private"
}

data "sbercloud_obs_buckets" "buckets" {
  bucket = "%s"

  depends_on = [sbercloud_obs_bucket.bucket]
}
`, name, name)
}
