package apig

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApigApiDebug_basic(t *testing.T) {
	var (
		rName        = acceptance.RandomAccResourceName()
		resourceName = "sbercloud_apig_api_debug.test_with_fgs"
	)

	// Avoid CheckDestroy because this resource is a one-time action resource.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApigApiDebug_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "request"),
					resource.TestCheckResourceAttrSet(resourceName, "response"),
					resource.TestCheckResourceAttrSet(resourceName, "latency"),
				),
			},
		},
	})
}

func testAccApiDebug_base(name string) string {
	return fmt.Sprintf(`
data "sbercloud_availability_zones" "test" {}

%[1]s

data "sbercloud_apig_instances" "test" {
  instance_id = "%[3]s"
}

resource "sbercloud_apig_group" "test" {
  instance_id = try(data.sbercloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  name        = "%[2]s"
}

resource "sbercloud_apig_environment" "test" {
  instance_id = try(data.sbercloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  name        = "%[2]s"
}

# Create FunctionGraph function
resource "sbercloud_fgs_function" "test" {
  name        = "%[2]s"
  app         = "default"
  description = "created by acc test for API debug"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = <<EOF
# -*- coding: utf-8 -*-
import json
def handler(event, context):
    return {
        'statusCode': 200,
        'body': json.dumps({
            'message': 'Hello, FunctionGraph!',
            'event': event
        })
    }
EOF
}

# Create API that calls the FunctionGraph function
resource "sbercloud_apig_api" "test" {
  instance_id      = try(data.sbercloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  group_id         = sbercloud_apig_group.test.id
  name             = "%[2]s"
  type             = "Private"
  request_protocol = "HTTP"
  request_method   = "POST"
  request_path     = "/test/function"

  request_params {
    name     = "test_query"
    type     = "STRING"
    location = "QUERY"
    required = false
  }
  request_params {
    name     = "header-param"
    type     = "STRING"
    location = "HEADER"
    required = true
  }

  backend_params {
    type     = "REQUEST"
    name     = "backend_query"
    location = "QUERY"
    value    = "test_query"
  }
  backend_params {
    type     = "REQUEST"
    name     = "backend-header"
    location = "HEADER"
    value    = "header-param"
  }

  func_graph {
    function_urn    = sbercloud_fgs_function.test.urn
    version         = "latest"
    timeout         = 10000
    invocation_type = "sync"
  }

  lifecycle {
    ignore_changes = [func_graph, request_params]
  }
}

resource "sbercloud_apig_api_action" "test_with_online" {
  instance_id = try(data.sbercloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  api_id      = sbercloud_apig_api.test.id
  env_id      = sbercloud_apig_environment.test.id
  action      = "online"
}
`, acceptance.TestBaseNetwork(name), name, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID)
}

func testAccApigApiDebug_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_apig_api_debug" "test_with_fgs" {
  instance_id = try(data.sbercloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  api_id      = sbercloud_apig_api.test.id
  mode        = "DEVELOPER"
  scheme      = "HTTP"
  method      = "POST"
  path        = "/test/function"
  body        = "{\"test\": \"data\"}"

  header = jsonencode({
    "Content-Type": ["application/json"],
    "test_param": ["test_value"]
  })

  query = jsonencode({
    "test_query": ["test_value"]
  })

  depends_on = [sbercloud_apig_api_action.test_with_online]
}

resource "sbercloud_apig_api_action" "test_with_offline" {
  instance_id = try(data.sbercloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  api_id      = sbercloud_apig_api.test.id
  env_id      = sbercloud_apig_environment.test.id
  action      = "offline"

  depends_on = [sbercloud_apig_api_debug.test_with_fgs]
}
`, testAccApiDebug_base(rName))
}
