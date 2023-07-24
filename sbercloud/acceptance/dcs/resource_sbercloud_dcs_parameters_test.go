package dcs

import (
	"fmt"
	"github.com/chnsz/golangsdk/openstack/dcs/v1/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"
)

// make testacc TEST='./sbercloud/acceptance/dcs' TESTARGS='-run TestAccDCSParameters_basic'
func TestAccDCSParameters_basic(t *testing.T) {
	var instanceName = fmt.Sprintf("testacc_dcs_instance_%s", acctest.RandString(5))
	var instance instances.Instance
	resourceInstanceName := "sbercloud_dcs_instance.instance_1"
	resourceParamsName := "sbercloud_dcs_parameters.test_new"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDCSParameters_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceInstanceName, instance),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.timeout", "1000"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.maxmemory-policy", "allkeys-lru"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.hash-max-ziplist-entries", "1024"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.hash-max-ziplist-value", "128"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.set-max-intset-entries", "1024"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.zset-max-ziplist-entries", "256"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.zset-max-ziplist-value", "128"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.latency-monitor-threshold", "1"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.maxclients", "2100"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.repl-backlog-size", "1049600"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.repl-backlog-ttl", "4000"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.appendfsync", "always"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.appendonly", "no"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.slowlog-log-slower-than", "15000"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.slowlog-max-len", "256"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.lua-time-limit", "1000"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.repl-timeout", "120"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.proto-max-bulk-len", "1048576"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.master-read-only", "yes"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.client-output-buffer-slave-soft-limit", "0"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.client-output-buffer-slave-hard-limit", "0"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.client-output-buffer-limit-slave-soft-seconds", "0"),
					resource.TestCheckResourceAttr(resourceParamsName, "parameters.active-expire-num", "100"),
				),
			},
		}})
}

func testAccDCSParameters_basic(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc" "vpc_1" {
  name = "vpc-default"
}
data "sbercloud_vpc_subnet" "subnet_1" {
  vpc_id = data.sbercloud_vpc.vpc_1.id
}
data "sbercloud_availability_zones" "test" {}

resource "sbercloud_dcs_instance" "instance_1" {
  name = "%s"
  engine             = "Redis"
  engine_version    = "5.0"
  flavor             = "redis.ha.xu1.large.r2.4"
  capacity          = 4
  vpc_id            = data.sbercloud_vpc.vpc_1.id
  subnet_id         = data.sbercloud_vpc_subnet.subnet_1.id
  availability_zones = [data.sbercloud_availability_zones.test.names[0]]
}

resource "sbercloud_dcs_parameters" "test_new" {
  instance_id = sbercloud_dcs_instance.instance_1.id
  project_id  = "%s"

  parameters = {
    timeout = "1000"
    maxmemory-policy = "allkeys-lru"
    hash-max-ziplist-entries = "1024"
    hash-max-ziplist-value = "128"
    set-max-intset-entries = "1024"
    zset-max-ziplist-entries = "256" 
    zset-max-ziplist-value = "128" 
    latency-monitor-threshold ="1"
    maxclients = "2100"
    repl-backlog-size = "1049600"
    repl-backlog-ttl = "4000"
    appendfsync = "always"
    appendonly = "no"
    slowlog-log-slower-than = "15000" 
    slowlog-max-len = "256"  
    lua-time-limit = "1000"
    repl-timeout = "120" 
    proto-max-bulk-len = "1048576" 
    master-read-only = "yes" 
    client-output-buffer-slave-soft-limit = "0" 
    client-output-buffer-slave-hard-limit = "0"
    client-output-buffer-limit-slave-soft-seconds = "0"
    active-expire-num = "100"
  }
}
`, rName, acceptance.SBC_PROJECT_ID)
}
