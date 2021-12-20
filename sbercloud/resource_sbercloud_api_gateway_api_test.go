package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/apigw/apis"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccApiGatewayAPI_basic(t *testing.T) {
	var resName = "sbercloud_api_gateway_api.acc_apigw_api"
	rName := fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApiGatewayApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigwAPI_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGatewayApiExists(resName),
					resource.TestCheckResourceAttr(resName, "name", rName),
					resource.TestCheckResourceAttr(resName, "group_name", rName),
					resource.TestCheckResourceAttr(resName, "cors", "false"),
					resource.TestCheckResourceAttr(resName, "auth_type", "NONE"),
					resource.TestCheckResourceAttr(resName, "backend_type", "HTTP"),
					resource.TestCheckResourceAttr(resName, "request_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resName, "request_method", "GET"),
					resource.TestCheckResourceAttr(resName, "request_uri", "/test/path1"),
					resource.TestCheckResourceAttr(resName, "http_backend.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resName, "http_backend.0.method", "GET"),
					resource.TestCheckResourceAttr(resName, "http_backend.0.uri", "/web/openapi"),
					resource.TestCheckResourceAttr(resName, "http_backend.0.timeout", "10000"),
				),
			},
			{
				Config: testAccApigwAPI_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGatewayApiExists(resName),
					resource.TestCheckResourceAttr(resName, "description", "updated by acc test"),
					resource.TestCheckResourceAttr(resName, "auth_type", "IAM"),
					resource.TestCheckResourceAttr(resName, "cors", "true"),
					resource.TestCheckResourceAttr(resName, "request_protocol", "BOTH"),
					resource.TestCheckResourceAttr(resName, "request_uri", "/test/path2"),
				),
			},
		},
	})
}

func testAccCheckApiGatewayApiDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	apigwClient, err := config.ApiGatewayV1Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud api gateway client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_api_gateway_api" {
			continue
		}

		_, err := apis.Get(apigwClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("api gateway API still exists")
		}
	}

	return nil
}

func testAccCheckApiGatewayApiExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource %s not found", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		apigwClient, err := config.ApiGatewayV1Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud api gateway client: %s", err)
		}

		found, err := apis.Get(apigwClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("apigateway API not found")
		}

		return nil
	}
}

func testAccApigwAPI_basic(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_api_gateway_group" "acc_apigw_group" {
  name        = "%s"
  description = "created by acc test"
}

resource "sbercloud_api_gateway_api" "acc_apigw_api" {
  group_id    = sbercloud_api_gateway_group.acc_apigw_group.id
  name        = "%s"
  visibility  = 2
  description = "created by acc test"
  tags        = ["tag1","tag2"]

  auth_type        = "NONE"
  backend_type     = "HTTP"
  request_protocol = "HTTPS"
  request_method   = "GET"
  request_uri      = "/test/path1"
  example_success_response = "this is a successful response"

  http_backend {
    protocol   = "HTTPS"
    method     = "GET"
    uri        = "/web/openapi"
    url_domain = "mysbercloud.com"
    timeout    = 10000
  }
}
`, rName, rName)
}

func testAccApigwAPI_update(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_api_gateway_group" "acc_apigw_group" {
  name        = "%s"
  description = "created by acc test"
}

resource "sbercloud_api_gateway_api" "acc_apigw_api" {
  group_id    = sbercloud_api_gateway_group.acc_apigw_group.id
  name        = "%s"
  visibility  = 2
  cors        = true
  description = "updated by acc test"
  tags        = ["tag1","tag2"]

  auth_type        = "IAM"
  backend_type     = "HTTP"
  request_protocol = "BOTH"
  request_method   = "GET"
  request_uri      = "/test/path2"
  example_success_response = "this is a successful response"

  http_backend {
    protocol   = "HTTPS"
    method     = "GET"
    uri        = "/web/openapi"
    url_domain = "mysbercloud.com"
    timeout    = 10000
  }
}
`, rName, rName)
}
