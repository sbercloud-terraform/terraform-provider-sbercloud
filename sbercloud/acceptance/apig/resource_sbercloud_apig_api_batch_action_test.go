package apig

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApigApiBatchAction_basic(t *testing.T) {
	var (
		name                  = acceptance.RandomAccResourceName()
		rcWithOnlineName      = "sbercloud_apig_api_batch_action.batch_online_apis_for_env"
		rcWithGroupOnlineName = "sbercloud_apig_api_batch_action.batch_online_apis_for_group"
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
				Config: testApigApiBatchAction_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rcWithOnlineName, "action", "online"),
					resource.TestCheckResourceAttr(rcWithOnlineName, "apis.#", "2"),
					resource.TestCheckResourceAttrSet(rcWithOnlineName, "id"),
				),
			},
			{
				Config: testApigApiBatchAction_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rcWithGroupOnlineName, "action", "online"),
					resource.TestCheckResourceAttrSet(rcWithGroupOnlineName, "group_id"),
					resource.TestCheckResourceAttrSet(rcWithGroupOnlineName, "id"),
				),
			},
		},
	})
}

func testApigApiBatchAction_base(name string) string {
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
  count = 2

  instance_id      = sbercloud_apig_instance.test.id
  group_id         = sbercloud_apig_group.test.id
  name             = format("%[2]s_%%d", count.index)
  type             = "Private"
  request_protocol = "HTTP"
  request_method   = "GET"
  request_path     = format("/mock/test%%d", count.index)

  mock {
    status_code = 200
  }
}
`, acceptance.TestBaseNetwork(name), name)
}

func testApigApiBatchAction_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_apig_api_batch_action" "batch_online_apis_for_env" {
  instance_id = sbercloud_apig_instance.test.id
  action      = "online"
  env_id      = sbercloud_apig_environment.test.id
  remark      = "Test batch action"
  apis        = sbercloud_apig_api.test[*].id
}

resource "sbercloud_apig_api_batch_action" "batch_offline_apis_for_env" {
  instance_id = sbercloud_apig_instance.test.id
  action      = "offline"
  env_id      = sbercloud_apig_environment.test.id
  remark      = "Test batch action"
  apis        = sbercloud_apig_api.test[*].id

  depends_on = [
    sbercloud_apig_api_batch_action.batch_online_apis_for_env,
  ]
}
`, testApigApiBatchAction_base(name))
}

func testApigApiBatchAction_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_apig_api_batch_action" "batch_online_apis_for_group" {
  instance_id = sbercloud_apig_instance.test.id
  action      = "online"
  env_id      = sbercloud_apig_environment.test.id
  group_id    = sbercloud_apig_group.test.id
  remark      = "Test batch action by group"

  depends_on = [
    sbercloud_apig_api.test,
  ]
}

resource "sbercloud_apig_api_batch_action" "batch_offline_apis_for_group" {
  instance_id = sbercloud_apig_instance.test.id
  action      = "offline"
  env_id      = sbercloud_apig_environment.test.id
  group_id    = sbercloud_apig_group.test.id
  remark      = "Test batch action by group"

  depends_on = [
    sbercloud_apig_api.test,
    sbercloud_apig_api_batch_action.batch_online_apis_for_group,
  ]
}
`, testApigApiBatchAction_base(name))
}
