package apig

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApiCheck_basic(t *testing.T) {
	var (
		name      = acceptance.RandomAccResourceName()
		uniqeName = acceptance.RandomAccResourceName()
	)

	// Avoid CheckDestroy because this resource is a one-time resource.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccApiCheck_basic_step1(name),
				ExpectError: regexp.MustCompile(`The instance does not exist`),
			},
			// Check whether the API name already exists in the same group.
			{
				Config:      testAccApiCheck_basic_step2(name),
				ExpectError: regexp.MustCompile(fmt.Sprintf("The API name already exists, api_name:%s", name)),
			},
			// Check whether the API path already exists in the same group.
			{
				Config:      testAccApiCheck_basic_step3(name),
				ExpectError: regexp.MustCompile(fmt.Sprintf("The API already exists, api_name:%s", name)),
			},
			{
				// Check the API name does not exist in the same group.
				Config: testAccApiCheck_basic_step4(name, uniqeName),
			},
			{
				// Check the API path does not exist in the same group.
				Config: testAccApiCheck_basic_step5(name, uniqeName),
			},
		},
	})
}

func testAccApiCheck_base(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_apig_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "sbercloud_apig_api" "test" {
  instance_id      = "%[1]s"
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
`, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccApiCheck_basic_step1(name string) string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
resource "sbercloud_apig_api_check" "test" {
  instance_id = "%[1]s"
  type        = "name"
  name        = "%[2]s"
}
`, randomId, name)
}

func testAccApiCheck_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_apig_api_check" "test" {
  instance_id = "%[2]s"
  type        = "name"
  name        = sbercloud_apig_api.test.name
  group_id    = sbercloud_apig_group.test.id

  depends_on = [sbercloud_apig_api.test]
}
`, testAccApiCheck_base(name), acceptance.SBC_APIG_DEDICATED_INSTANCE_ID)
}

func testAccApiCheck_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_apig_api_check" "test" {
  instance_id = "%[2]s"
  type        = "path"
  group_id    = sbercloud_apig_group.test.id
  req_method  = sbercloud_apig_api.test.request_method
  req_uri     = sbercloud_apig_api.test.request_path
  match_mode  = "NORMAL"
}
`, testAccApiCheck_base(name), acceptance.SBC_APIG_DEDICATED_INSTANCE_ID)
}

func testAccApiCheck_basic_step4(name, uniqeName string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_apig_api_check" "check_name" {
  instance_id = "%[2]s"
  type        = "name"
  name        =  "%[3]s"
  group_id    = sbercloud_apig_group.test.id
}
`, testAccApiCheck_base(name), acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, uniqeName)
}

func testAccApiCheck_basic_step5(name, uniqeName string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_apig_api_check" "check_path" {
  instance_id = "%[2]s"
  type        = "path"
  group_id    = sbercloud_apig_group.test.id
  req_method  = sbercloud_apig_api.test.request_method
  req_uri     = "/test/%[3]s"
  match_mode  = "NORMAL"
}
`, testAccApiCheck_base(name), acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, uniqeName)
}
