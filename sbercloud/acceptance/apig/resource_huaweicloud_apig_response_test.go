package apig

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/responses"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getResponseFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return responses.Get(client, state.Primary.Attributes["instance_id"], state.Primary.Attributes["group_id"],
		state.Primary.ID).Extract()
}

func TestAccResponse_basic(t *testing.T) {
	var (
		resp responses.Response

		rName      = "sbercloud_apig_response.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&resp,
		getResponseFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResponse_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
				),
			},
			{
				Config: testAccResponse_basic(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResponseImportStateFunc(rName),
			},
		},
	})
}

func TestAccResponse_customRules(t *testing.T) {
	var (
		resp responses.Response

		rName = "sbercloud_apig_response.test"
		name  = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&resp,
		getResponseFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResponse_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
				),
			},
			{
				// Add two custom rule.
				Config: testAccResponse_rules(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "rule.#", "2"),
				),
			},
			{
				// Remove one and update another.
				Config: testAccResponse_rulesUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "rule.#", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResponseImportStateFunc(rName),
			},
		},
	})
}

func testAccResponseImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rsName)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["group_id"] == "" ||
			rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("missing some attributes, want '{instance_id}/{group_id}/{name}', but '%s/%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["group_id"], rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["group_id"],
			rs.Primary.Attributes["name"]), nil
	}
}

func testAccResponse_basic(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_apig_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "sbercloud_apig_response" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  group_id    = sbercloud_apig_group.test.id

  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 401
  }
}
`, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccResponse_rules(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_apig_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "sbercloud_apig_response" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  group_id    = sbercloud_apig_group.test.id

  rule {
    error_type  = "ACCESS_DENIED"
    body        = "{\"error_code\":\"$context.error.code\",\"error_msg\":\"$context.error.message\"}"
    status_code = 400
  }
  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 401
  }
}
`, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccResponse_rulesUpdate(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_apig_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "sbercloud_apig_response" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  group_id    = sbercloud_apig_group.test.id

  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 403
  }
}
`, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}
