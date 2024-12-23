package apig

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/plugins"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getPluginFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	return plugins.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccPlugin_basic(t *testing.T) {
	var (
		plugin plugins.Plugin

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		rNameForCors = "sbercloud_apig_plugin.cors"
		rcForCors    = acceptance.InitResourceCheck(rNameForCors, &plugin, getPluginFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcForCors.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Check whether illegal type ​​can be intercepted normally (create phase).
				Config:      testAccPlugin_basic_step1(name),
				ExpectError: regexp.MustCompile("error creating the plugin"),
			},
			{
				Config: testAccPlugin_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcForCors.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameForCors, "instance_id", acceptance.SBC_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(rNameForCors, "name", name),
					resource.TestCheckResourceAttr(rNameForCors, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(rNameForCors, "type", "cors"),
					resource.TestCheckResourceAttrSet(rNameForCors, "created_at"),
				),
			},
			{
				// Check whether illegal content value ​​can be intercepted normally (update phase).
				Config:      testAccPlugin_basic_step3(name),
				ExpectError: regexp.MustCompile("error updating the plugin"),
			},
			{
				Config: testAccPlugin_basic_step4(updateName),
				Check: resource.ComposeTestCheckFunc(
					rcForCors.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameForCors, "name", updateName),
					resource.TestCheckResourceAttr(rNameForCors, "description", "Updated by acc test"),
					resource.TestCheckResourceAttrSet(rNameForCors, "updated_at"),
				),
			},
			{
				ResourceName:      rNameForCors,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPluginImportStateFunc(rNameForCors),
			},
		},
	})
}

func testAccPluginImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rsName)
		}
		if rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}

func testAccPlugin_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_apig_plugin" "cors" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = "Created by acc test"
  type        = "INVALID_TYPE"
  content     = jsonencode(
    {
      allow_origin      = "*"
      allow_methods     = "GET,PUT,DELETE,HEAD,PATCH"
      allow_headers     = "Content-Type,Accept,Cache-Control"
      expose_headers    = "X-Request-Id,X-Apig-Latency"
      max_age           = 12700
      allow_credentials = true
    }
  )
}
`, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccPlugin_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_apig_plugin" "cors" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = "Created by acc test"
  type        = "cors"
  content     = jsonencode(
    {
      allow_origin      = "*"
      allow_methods     = "GET,PUT,DELETE,HEAD,PATCH"
      allow_headers     = "Content-Type,Accept,Cache-Control"
      expose_headers    = "X-Request-Id,X-Apig-Latency"
      max_age           = 12700
      allow_credentials = true
    }
  )
}
`, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccPlugin_basic_step3(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_apig_plugin" "cors" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = "Created by acc test"
  type        = "cors"
  content     = "INVALID_CONTENT"
}
`, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccPlugin_basic_step4(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_apig_plugin" "cors" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = "Updated by acc test" # Description cannot be updated to a empty value.
  type        = "cors"
  content     = jsonencode(
    {
      allow_origin      = "*.terraform.test.com"
      allow_methods     = "POST,PATCH"
      allow_headers     = "Content-Type,Accept,Accept-Ranges"
      expose_headers    = "X-Request-Id,X-Apig-Auth-Type"
      max_age           = 800
      allow_credentials = false
    }
  )
}
`, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}

func TestAccPlugin_httpResponse(t *testing.T) {
	var (
		plugin plugins.Plugin

		name = acceptance.RandomAccResourceName()

		rNameForHttpResponse = "sbercloud_apig_plugin.http_response"
		rcForHttpResponse    = acceptance.InitResourceCheck(rNameForHttpResponse, &plugin, getPluginFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcForHttpResponse.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPlugin_httpResponse_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcForHttpResponse.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameForHttpResponse, "type", "set_resp_headers"),
					resource.TestCheckResourceAttrSet(rNameForHttpResponse, "created_at"),
				),
			},
			{
				Config: testAccPlugin_httpResponse_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcForHttpResponse.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rNameForHttpResponse, "updated_at"),
				),
			},
			{
				ResourceName:      rNameForHttpResponse,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPluginImportStateFunc(rNameForHttpResponse),
			},
		},
	})
}

func testAccPlugin_httpResponse_step1(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_apig_plugin" "http_response" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  type        = "set_resp_headers"
  content     = jsonencode(
    {
      response_headers = [{
        name       = "X-Custom-Pwd"
        value      = "**********"
        value_type = "custom_value"
        action     = "override"
      },
      {
        name       = "X-Custom-Debug-Step"
        value      = "Beta"
        value_type = "custom_value"
        action     = "skip"
      },
      {
        name       = "X-Custom-Config"
        value      = "<HTTP response test>"
        action     = "append"
        value_type = "custom_value"
      },
      {
        name       = "X-Custom-Id"
        value      = ""
        value_type = "custom_value"
        action     = "delete"
      },
      {
        name       = "X-Custom-Log-Level"
        value      = "DEBUG"
        value_type = "custom_value"
        action     = "add"
      },
      {
        name       = "Sys-Param"
        value      = "$context.cacheStatus"
        value_type = "system_parameter"
        action     = "add"
      }]
    }
  )
}
`, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccPlugin_httpResponse_step2(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_apig_plugin" "http_response" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  type        = "set_resp_headers"
  content     = jsonencode(
    {
      response_headers = [{
        name       = "X-Custom-Pwd"
        value      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
        value_type = "custom_value"
        action     = "delete"
      },
      {
        name       = "X-Custom-Log-PATH"
        value      = "/tmp/debug.log"
        value_type = "custom_value"
        action     = "add"
      },
      {
        name       = "Sys-Param-updated"
        value      = "$context.cacheStatus"
        value_type = "system_parameter"
        action     = "append"
      }]
    }
  )
}
`, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}

func TestAccPlugin_rateLimit(t *testing.T) {
	var (
		plugin plugins.Plugin

		name = acceptance.RandomAccResourceName()

		rNameForRateLimit = "sbercloud_apig_plugin.rate_limit"
		rcForRateLimit    = acceptance.InitResourceCheck(rNameForRateLimit, &plugin, getPluginFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcForRateLimit.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPlugin_rateLimit_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcForRateLimit.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameForRateLimit, "type", "rate_limit"),
					resource.TestCheckResourceAttrSet(rNameForRateLimit, "created_at"),
				),
			},
			{
				Config: testAccPlugin_rateLimit_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcForRateLimit.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rNameForRateLimit, "updated_at"),
				),
			},
			{
				ResourceName:      rNameForRateLimit,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPluginImportStateFunc(rNameForRateLimit),
			},
		},
	})
}

func testAccPlugin_rateLimit_step1(name string) string {
	return fmt.Sprintf(`
data "sbercloud_identity_users" "test" {}

resource "sbercloud_apig_application" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "sbercloud_apig_plugin" "rate_limit" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  type        = "rate_limit"
  content     = jsonencode(
    {
      "scope": "basic",
      "default_time_unit": "minute",
      "default_interval": 1,
      "api_limit": 25,
      "app_limit": 10,
      "user_limit": 15,
      "ip_limit": 25,
      "algorithm": "counter",
      "specials": [
        {
          "type": "app",
          "policies": [
            {
              "key": "${sbercloud_apig_application.test.id}",
              "limit": 10
            }
          ]
        },
        {
          "type": "user",
          "policies": [
            {
              "key": "${data.sbercloud_identity_users.test.users[0].id}",
              "limit": 10
            }
          ]
        }
      ],
      "parameters": [
        {
          "type": "path",
          "name": "reqPath",
          "value": "reqPath"
        },
        {
          "type": "method",
          "name": "method",
          "value": "method"
        },
        {
          "type": "system",
          "name": "serverName",
          "value": "serverName"
        }
      ],
      "rules": [
        {
          "rule_name": "rule-0001",
          "match_regex": "[\"AND\",[\"method\",\"~=\",\"POST\"],[\"method\",\"~=\",\"PATCH\"]]",
          "time_unit": "minute",
          "interval": 1,
          "limit": 20
        },
        {
          "rule_name": "rule-0002",
          "match_regex": "[\"reqPath\",\"~~\",\"/terraform/test/*/\"]",
          "time_unit": "minute",
          "interval": 1,
          "limit": 10
        },
        {
          "rule_name": "rule-0003",
          "match_regex": "[\"serverName\",\"==\",\"terraform\"]",
          "time_unit": "minute",
          "interval": 1,
          "limit": 15
        },
        {
          "rule_name": "rule-0004",
          "match_regex": "[\"method\",\"in\",\"PATCH\"]",
          "time_unit": "minute",
          "interval": 1,
          "limit": 5
        }
      ]
    }
  )
}
`, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccPlugin_rateLimit_step2(name string) string {
	return fmt.Sprintf(`
data "sbercloud_identity_users" "test" {}

resource "sbercloud_apig_application" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "sbercloud_apig_plugin" "rate_limit" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  type        = "rate_limit"
  content     = jsonencode(
    {
      "scope": "basic",
      "default_time_unit": "minute",
      "default_interval": 2,
      "api_limit": 35,
      "app_limit": 15,
      "user_limit": 25,
      "ip_limit": 30,
      "algorithm": "haht",
      "specials": [
        {
          "type": "app",
          "policies": [
            {
              "key": "${sbercloud_apig_application.test.id}",
              "limit": 15
            }
          ]
        },
        {
          "type": "user",
          "policies": [
            {
              "key": "${data.sbercloud_identity_users.test.users[0].id}",
              "limit": 15
            }
          ]
        }
      ],
      "parameters": [
        {
          "type": "path",
          "name": "reqPath",
          "value": "reqPath"
        },
        {
          "type": "method",
          "name": "method",
          "value": "method"
        },
        {
          "type": "system",
          "name": "serverName",
          "value": "serverName"
        }
      ],
      "rules": [
        {
          "rule_name": "rule-0001",
          "match_regex": "[\"AND\",[\"method\",\"~=\",\"POST\"],[\"method\",\"~=\",\"PATCH\"]]",
          "time_unit": "minute",
          "interval": 1,
          "limit": 25
        },
        {
          "rule_name": "rule-0002",
          "match_regex": "[\"reqPath\",\"~~\",\"/terraform/test/*/\"]",
          "time_unit": "minute",
          "interval": 2,
          "limit": 15
        },
        {
          "rule_name": "rule-0003",
          "match_regex": "[\"serverName\",\"==\",\"terraform\"]",
          "time_unit": "minute",
          "interval": 1,
          "limit": 20
        },
        {
          "rule_name": "rule-0004",
          "match_regex": "[\"method\",\"in\",\"PATCH\"]",
          "time_unit": "minute",
          "interval": 1,
          "limit": 15
        }
      ]
    }
  )
}
`, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}

func TestAccPlugin_kafkaLog(t *testing.T) {
	var (
		plugin plugins.Plugin

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		baseConfig = testAccPlugin_kafkaLog_base(name)

		rNameForKafkaLog = "sbercloud_apig_plugin.kafka_log"
		rcForKafkaLog    = acceptance.InitResourceCheck(rNameForKafkaLog, &plugin, getPluginFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcForKafkaLog.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPlugin_kafkaLog_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rcForKafkaLog.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameForKafkaLog, "type", "kafka_log"),
					resource.TestCheckResourceAttrSet(rNameForKafkaLog, "created_at"),
				),
			},
			{
				Config: testAccPlugin_kafkaLog_step2(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rcForKafkaLog.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rNameForKafkaLog, "updated_at"),
				),
			},
			{
				ResourceName:      rNameForKafkaLog,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPluginImportStateFunc(rNameForKafkaLog),
			},
		},
	})
}

func testAccPlugin_kafkaLog_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "sbercloud_dms_kafka_flavors" "test" {
  type = "cluster"
}

locals {
  query_results     = data.sbercloud_dms_kafka_flavors.test
  flavor            = data.sbercloud_dms_kafka_flavors.test.flavors[0]
  connect_addresses = split(",", sbercloud_dms_kafka_instance.test.connect_address)
  connect_port      = sbercloud_dms_kafka_instance.test.port
}

resource "sbercloud_dms_kafka_instance" "test" {
  name              = "%[2]s"
  vpc_id            = sbercloud_vpc.test.id
  network_id        = sbercloud_vpc_subnet.test.id
  security_group_id = sbercloud_networking_secgroup.test.id

  flavor_id          = local.flavor.id
  storage_spec_code  = local.flavor.ios[0].storage_spec_code
  availability_zones = local.flavor.ios[0].availability_zones
  engine_version     = element(local.query_results.versions, length(local.query_results.versions)-1)
  storage_space      = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  broker_num         = 3
  enable_auto_topic  = true

  access_user      = "user"
  password         = "Kafkatest@123"
  manager_user     = "kafka-user"
  manager_password = "Kafkatest@123"

  lifecycle {
    ignore_changes = [
      availability_zones, manager_password, password,
    ]
  }
}

resource "sbercloud_dms_kafka_topic" "test" {
  instance_id = sbercloud_dms_kafka_instance.test.id
  name        = "%[2]s"
  partitions  = 1
}
`, acceptance.TestBaseNetwork(name), name)
}

func testAccPlugin_kafkaLog_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_apig_plugin" "kafka_log" {
  instance_id = "%[2]s"
  name        = "%[3]s"
  type        = "kafka_log"
  content     = jsonencode(
    {
      "broker_list": [for v in local.connect_addresses: format("%%s:%%d", v, local.connect_port)],
      "topic": "${sbercloud_dms_kafka_topic.test.name}",
      "key": "",
      "max_retry_count": 0,
      "retry_backoff": 1,
      "sasl_config": {
        "security_protocol": "PLAINTEXT",
        "sasl_mechanisms": "PLAIN",
        "sasl_username": "",
        "sasl_password": "",
        "ssl_ca_content": ""
      },
      "meta_config": {
        "system": {
          "start_time": false,
          "request_id": false,
          "client_ip": false,
          "api_id": false,
          "user_name": false,
          "app_id": false,
          "access_model1": false,
          "request_time": true,
          "http_status": true,
          "server_protocol": false,
          "scheme": true,
          "request_method": true,
          "host": false,
          "api_uri_mode": false,
          "uri": false,
          "request_size": false,
          "response_size": false,
          "upstream_uri": false,
          "upstream_addr": false,
          "upstream_status": true,
          "upstream_connect_time": false,
          "upstream_header_time": false,
          "upstream_response_time": false,
          "all_upstream_response_time": false,
          "region_id": true,
          "auth_type": false,
          "http_x_forwarded_for": false,
          "http_user_agent": false,
          "error_type": false,
          "access_model2": false,
          "inner_time": false,
          "proxy_protocol_vni": false,
          "proxy_protocol_vpce_id": false,
          "proxy_protocol_addr": false,
          "body_bytes_sent": false,
          "api_name": true,
          "app_name": true,
          "provider_app_id": false,
          "provider_app_name": false,
          "custom_data_log01": false,
          "custom_data_log02": false,
          "custom_data_log03": false,
          "custom_data_log04": false,
          "custom_data_log05": false,
          "custom_data_log06": false,
          "custom_data_log07": false,
          "custom_data_log08": false,
          "custom_data_log09": false,
          "custom_data_log10": false,
          "response_source": false
        },
        "call_data": {
          "log_request_header": false,
          "log_request_query_string": false,
          "log_request_body": false,
          "log_response_header": false,
          "log_response_body": false,
          "request_header_filter": "",
          "request_query_string_filter": "",
          "response_header_filter": "",
          "custom_authorizer": {
            "frontend": [],
            "backend": []
          }
        }
      }
    }
  )
}
`, baseConfig, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccPlugin_kafkaLog_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_apig_plugin" "kafka_log" {
  instance_id = "%[2]s"
  name        = "%[3]s"
  type        = "kafka_log"
  content     = jsonencode(
    {
      "broker_list": [for v in local.connect_addresses: format("%%s:%%d", v, local.connect_port)],
      "topic": "${sbercloud_dms_kafka_topic.test.name}",
      "key": "",
      "max_retry_count": 3,
      "retry_backoff": 10,
      "sasl_config": {
        "security_protocol": "PLAINTEXT",
        "sasl_mechanisms": "PLAIN",
        "sasl_username": "",
        "sasl_password": "",
        "ssl_ca_content": ""
      },
      "meta_config": {
        "system": {
          "start_time": true,
          "request_id": true,
          "client_ip": true,
          "api_id": false,
          "user_name": false,
          "app_id": false,
          "access_model1": false,
          "request_time": true,
          "http_status": true,
          "server_protocol": false,
          "scheme": true,
          "request_method": true,
          "host": false,
          "api_uri_mode": false,
          "uri": false,
          "request_size": false,
          "response_size": false,
          "upstream_uri": false,
          "upstream_addr": true,
          "upstream_status": true,
          "upstream_connect_time": false,
          "upstream_header_time": false,
          "upstream_response_time": true,
          "all_upstream_response_time": false,
          "region_id": false,
          "auth_type": false,
          "http_x_forwarded_for": true,
          "http_user_agent": true,
          "error_type": true,
          "access_model2": false,
          "inner_time": false,
          "proxy_protocol_vni": false,
          "proxy_protocol_vpce_id": false,
          "proxy_protocol_addr": false,
          "body_bytes_sent": false,
          "api_name": false,
          "app_name": false,
          "provider_app_id": false,
          "provider_app_name": false,
          "custom_data_log01": false,
          "custom_data_log02": false,
          "custom_data_log03": false,
          "custom_data_log04": false,
          "custom_data_log05": false,
          "custom_data_log06": false,
          "custom_data_log07": false,
          "custom_data_log08": false,
          "custom_data_log09": false,
          "custom_data_log10": false,
          "response_source": false
        },
        "call_data": {
          "log_request_header": true,
          "log_request_query_string": true,
          "log_request_body": true,
          "log_response_header": true,
          "log_response_body": true,
          "request_header_filter": "X-Custom-Auth-Type",
          "request_query_string_filter": "authId",
          "response_header_filter": "X-Trace-Id",
          "custom_authorizer": {
            "frontend": [
              "user_name",
              "user_age"
            ],
            "backend": [
              "userName",
              "userAge"
            ]
          }
        }
      }
    }
  )
}
`, baseConfig, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}

func TestAccPlugin_breaker(t *testing.T) {
	var (
		plugin plugins.Plugin

		name = acceptance.RandomAccResourceName()

		rNameForBreaker = "sbercloud_apig_plugin.breaker"
		rcForBreaker    = acceptance.InitResourceCheck(rNameForBreaker, &plugin, getPluginFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcForBreaker.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPlugin_breaker_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcForBreaker.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameForBreaker, "type", "breaker"),
					resource.TestCheckResourceAttrSet(rNameForBreaker, "created_at"),
				),
			},
			{
				Config: testAccPlugin_breaker_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcForBreaker.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rNameForBreaker, "updated_at"),
				),
			},
			{
				ResourceName:      rNameForBreaker,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPluginImportStateFunc(rNameForBreaker),
			},
		},
	})
}

func testAccPlugin_breaker_step1(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_apig_plugin" "breaker" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  type        = "breaker"
  content     = jsonencode(
    {
      "breaker_condition": {
        "breaker_type": "timeout",
        "breaker_mode": "percentage",
        "unhealthy_condition": "",
        "unhealthy_threshold": 30,
        "min_call_threshold": 20,
        "unhealthy_percentage": 51,
        "time_window": 15,
        "open_breaker_time": 15
      },
      "downgrade_default": null,
      "downgrade_parameters": null,
      "downgrade_rules": null,
      "scope": "share"
    }
  )
}
`, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccPlugin_breaker_step2(name string) string {
	return fmt.Sprintf(`
resource "sbercloud_apig_plugin" "breaker" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  type        = "breaker"
  content     = jsonencode(
    {
      "breaker_condition": {
        "breaker_type": "condition",
        "breaker_mode": "counter",
        "unhealthy_condition": "[\"OR\",[\"$context.statusCode\",\"in\",\"500,501,504\"],[\"$context.backendResponseTime\",\">\",6000]]",
        "unhealthy_threshold": 30,
        "min_call_threshold": 20,
        "unhealthy_percentage": 51,
        "time_window": 15,
        "open_breaker_time": 15
      },
      "downgrade_default": null,
      "downgrade_parameters": [
        {
          "type": "path",
          "name": "reqPath",
          "value": "reqPath"
        },
        {
          "type": "method",
          "name": "method",
          "value": "method"
        },
        {
          "type": "query",
          "name": "authType",
          "value": "authType"
        }
      ],
      "downgrade_rules": [
        {
          "breaker_condition": {
            "breaker_type": "timeout",
            "breaker_mode": "percentage",
            "unhealthy_condition": "",
            "unhealthy_threshold": 30,
            "min_call_threshold": 20,
            "unhealthy_percentage": 51,
            "time_window": 15,
            "open_breaker_time": 15
          },
          "downgrade_backend": null,
          "rule_name": "rule-qkqe",
          "match_regex": "[\"authType\",\"~=\",\"Token\"]",
          "parameters": [
            "reqPath",
            "method",
            "authType"
          ]
        }
      ],
      "scope": "basic"
    }
  )
}
`, acceptance.SBC_APIG_DEDICATED_INSTANCE_ID, name)
}
