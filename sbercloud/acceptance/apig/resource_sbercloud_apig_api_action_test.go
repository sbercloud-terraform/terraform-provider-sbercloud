package apig

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApigApiAction_basic(t *testing.T) {
	var (
		rName        = acceptance.RandomAccResourceName()
		resourceName = "sbercloud_apig_api_action.test_with_online"
	)

	// Avoid CheckDestroy because this resource is a one-time action resource.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApigApiAction_basic_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "publish_id"),
					resource.TestCheckResourceAttrSet(resourceName, "api_name"),
					resource.TestCheckResourceAttrSet(resourceName, "publish_time"),
					resource.TestCheckResourceAttrSet(resourceName, "version_id"),
				),
			},
			{
				Config: testAccApigApiAction_basic_step2(rName),
			},
		},
	})
}

func testAccApiAction_base(name string) string {
	return fmt.Sprintf(`
data "sbercloud_availability_zones" "test" {}

%[1]s

resource "sbercloud_apig_instance" "test" {
  vpc_id                = sbercloud_vpc.test.id
  subnet_id             = sbercloud_vpc_subnet.test.id
  security_group_id     = sbercloud_networking_secgroup.test.id
  enterprise_project_id = "0"
  availability_zones    = try(slice(data.sbercloud_availability_zones.test.names, 0, 1), [])
  edition               = "BASIC"
  name                  = "%[2]s"
  description           = "created by acc test for API offline action"
}

resource "sbercloud_apig_group" "test" {
  instance_id = sbercloud_apig_instance.test.id
  name        = "%[2]s"
}

resource "sbercloud_apig_environment" "test" {
  instance_id = sbercloud_apig_instance.test.id
  name        = "%[2]s"
}

resource "sbercloud_apig_api" "test" {
  instance_id      = sbercloud_apig_instance.test.id
  group_id         = sbercloud_apig_group.test.id
  name             = "%[2]s"
  type             = "Private"
  request_protocol = "HTTP"
  request_method   = "GET"
  request_path     = "/mock/test"

  mock {
    status_code = 200
  }
}
`, acceptance.TestBaseNetwork(name), name)
}

func testAccApigApiAction_basic_step1(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_apig_api_action" "test_with_online" {
  instance_id = sbercloud_apig_instance.test.id
  api_id      = sbercloud_apig_api.test.id
  env_id      = sbercloud_apig_environment.test.id
  action      = "online"
}
`, testAccApiAction_base(rName))
}

func testAccApigApiAction_basic_step2(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_apig_api_action" "test_with_online" {
  instance_id = sbercloud_apig_instance.test.id
  api_id      = sbercloud_apig_api.test.id
  env_id      = sbercloud_apig_environment.test.id
  action      = "online"
}

resource "sbercloud_apig_api_action" "test_with_offline" {
  instance_id = sbercloud_apig_instance.test.id
  api_id      = sbercloud_apig_api.test.id
  env_id      = sbercloud_apig_environment.test.id
  action      = "offline"

  depends_on = [sbercloud_apig_api_action.test_with_online]
}
`, testAccApiAction_base(rName))
}
