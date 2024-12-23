package apig

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAppcodes_basic(t *testing.T) {
	var (
		rName = "data.sbercloud_apig_appcodes.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		rNameNotFound = "data.sbercloud_apig_appcodes.not_found"
		dcNotFound    = acceptance.InitDataSourceCheck(rNameNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppcodes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "appcodes.#", "1"),
					dcNotFound.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameNotFound, "appcodes.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceAppcodes_basic_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

data "sbercloud_apig_instances" "test" {
  instance_id = "%[2]s"
}

locals {
  instance_id = data.sbercloud_apig_instances.test.instances[0].id
}

resource "sbercloud_apig_application" "test" {
  count = 2

  instance_id = local.instance_id
  name        = format("%[3]s_%%d", count.index)
}

resource "sbercloud_apig_appcode" "test" {
  instance_id    = local.instance_id
  application_id = sbercloud_apig_application.test[0].id
}
`, acceptance.TestBaseNetwork(name), acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccDataSourceAppcodes_basic() string {
	return fmt.Sprintf(`
%s

data "sbercloud_apig_appcodes" "test" {
  depends_on = [
    sbercloud_apig_appcode.test
  ]

  instance_id    = local.instance_id
  application_id = sbercloud_apig_application.test[0].id
}

data "sbercloud_apig_appcodes" "not_found" {
  instance_id    = local.instance_id
  application_id = sbercloud_apig_application.test[1].id
}
`, testAccDataSourceAppcodes_basic_base())
}
