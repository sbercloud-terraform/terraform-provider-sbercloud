package fgs

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud"
	"testing"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/trigger"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func getTriggerResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.FgsV2Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SberCloud FunctionGraph v2 client: %s", err)
	}
	return trigger.Get(c, state.Primary.Attributes["function_urn"], state.Primary.Attributes["type"],
		state.Primary.ID).Extract()
}

func TestAccFunctionGraphTrigger_basic(t *testing.T) {
	var (
		timeTrigger  trigger.Trigger
		randName     = acceptance.RandomAccResourceName()
		resourceName = "sbercloud_fgs_trigger.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&timeTrigger,
		getTriggerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphTimingTrigger_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "TIMER"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule_type", "Rate"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule", "3d"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "function_urn",
						"${sbercloud_fgs_function.test.urn}"),
				),
			},
			{
				Config: testAccFunctionGraphTimingTrigger_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "TIMER"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule_type", "Rate"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule", "3d"),
					resource.TestCheckResourceAttr(resourceName, "status", "DISABLED"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "function_urn",
						"${sbercloud_fgs_function.test.urn}"),
				),
			},
		},
	})
}

func TestAccFunctionGraphTrigger_cronTimer(t *testing.T) {
	var (
		randName     = acceptance.RandomAccResourceName()
		resourceName = "sbercloud_fgs_trigger.test"
		timeTrigger  trigger.Trigger
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&timeTrigger,
		getTriggerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphTimingTrigger_cron(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "TIMER"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule_type", "Cron"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule", "@every 1h30m"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "function_urn",
						"${sbercloud_fgs_function.test.urn}"),
				),
			},
			{
				Config: testAccFunctionGraphTimingTrigger_cronUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "TIMER"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule_type", "Cron"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule", "@every 1h30m"),
					resource.TestCheckResourceAttr(resourceName, "status", "DISABLED"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "function_urn",
						"${sbercloud_fgs_function.test.urn}"),
				),
			},
		},
	})
}

func TestAccFunctionGraphTrigger_smn(t *testing.T) {
	var (
		randName     = acceptance.RandomAccResourceName()
		resourceName = "sbercloud_fgs_trigger.test"
		timeTrigger  trigger.Trigger
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&timeTrigger,
		getTriggerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphSmnTrigger_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "SMN"),
					resource.TestCheckResourceAttrSet(resourceName, "smn.0.topic_urn"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "function_urn",
						"${sbercloud_fgs_function.test.urn}"),
				),
			},
		},
	})
}

func TestAccFunctionGraphTrigger_lts(t *testing.T) {
	var (
		randName     = acceptance.RandomAccResourceName()
		resourceName = "sbercloud_fgs_trigger.test"
		ltsTrigger   trigger.Trigger
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ltsTrigger,
		getTriggerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckFgsTrigger(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphLtsTrigger_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "LTS"),
					resource.TestCheckResourceAttrSet(resourceName, "lts.0.log_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "lts.0.log_topic_id"),
				),
			},
		},
	})
}

func testAccFunctionGraphTimingTrigger_base(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_fgs_function" "test" {
  name        = "%s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 10
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}`, rName)
}

func testAccFunctionGraphTimingTrigger_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_fgs_trigger" "test" {
  function_urn = sbercloud_fgs_function.test.urn
  type         = "TIMER"

  timer {
    name          = "%s"
    schedule_type = "Rate"
    schedule      = "3d"
  }
}
`, testAccFunctionGraphTimingTrigger_base(rName), rName)
}

func testAccFunctionGraphTimingTrigger_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_fgs_trigger" "test" {
  function_urn = sbercloud_fgs_function.test.urn
  type         = "TIMER"
  status       = "DISABLED"

  timer {
	name          = "%s"
	schedule_type = "Rate"
	schedule      = "3d"
  }
}
`, testAccFunctionGraphTimingTrigger_base(rName), rName)
}

func testAccFunctionGraphTimingTrigger_cron(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_fgs_trigger" "test" {
  function_urn = sbercloud_fgs_function.test.urn
  type         = "TIMER"

  timer {
    name          = "%s"
    schedule_type = "Cron"
    schedule      = "@every 1h30m"
  }
}
`, testAccFunctionGraphTimingTrigger_base(rName), rName)
}

func testAccFunctionGraphTimingTrigger_cronUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_fgs_trigger" "test" {
  function_urn = sbercloud_fgs_function.test.urn
  type         = "TIMER"
  status       = "DISABLED"

  timer {
	name          = "%s"
	schedule_type = "Cron"
	schedule      = "@every 1h30m"
  }
}
`, testAccFunctionGraphTimingTrigger_base(rName), rName)
}

func testAccFunctionGraphSmnTrigger_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_smn_topic" "test" {
  name = "%s"
}

resource "sbercloud_fgs_trigger" "test" {
  function_urn = sbercloud_fgs_function.test.urn
  type         = "SMN"

  smn {
    topic_urn = sbercloud_smn_topic.test.topic_urn
  }
}`, testAccFunctionGraphTimingTrigger_base(rName), rName)
}

func testAccFunctionGraphLtsTrigger_basic(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 1
}

resource "sbercloud_lts_stream" "test" {
  group_id    = sbercloud_lts_group.test.id
  stream_name = "%[1]s"
}

resource "sbercloud_identity_agency" "test" {
  name = "%[1]s"
  delegated_service_name = "%[3]s"

  project_role {
    project = "%[2]s"
    roles = ["LTS FullAccess"]
  }
}

resource "sbercloud_fgs_function" "test" {
  name        = "%[1]s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 10
  runtime     = "Python2.7"
  code_type   = "inline"
  agency      = sbercloud_identity_agency.test.name
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}

resource "sbercloud_fgs_trigger" "test" {
  function_urn = sbercloud_fgs_function.test.urn
  type         = "LTS"

  lts {
    log_group_id = sbercloud_lts_group.test.id
    log_topic_id = sbercloud_lts_stream.test.id
  }
}`, rName, acceptance.SBC_REGION_NAME, acceptance.SBC_FGS_TRIGGER_LTS_AGENCY)
}

func testAccNetwork_config(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_networking_secgroup_rule" "test" {
  security_group_id = sbercloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 9092
  port_range_max    = 9092
  remote_ip_prefix  = "0.0.0.0/0"
}`, sbercloud.TestBaseNetwork(rName))
}

func testAccDmsKafka_config(rName, password string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_dms_az" "test" {}

data "sbercloud_dms_product" "test" {
  engine            = "kafka"
  version           = "1.1.0"
  instance_type     = "cluster"
  partition_num     = 300
  storage           = 600
  storage_spec_code = "dms.physical.storage.high"
}

resource "sbercloud_dms_kafka_instance" "test" {
  name              = "%s"
  vpc_id            = sbercloud_vpc.test.id
  network_id        = sbercloud_vpc_subnet.test.id
  security_group_id = sbercloud_networking_secgroup.test.id
  available_zones   = [data.sbercloud_dms_az.test.id]
  product_id        = data.sbercloud_dms_product.test.id
  engine_version    = data.sbercloud_dms_product.test.version
  bandwidth         = data.sbercloud_dms_product.test.bandwidth
  storage_space     = data.sbercloud_dms_product.test.storage
  storage_spec_code = data.sbercloud_dms_product.test.storage_spec_code
  manager_user      = "%s"
  manager_password  = "%s"
}

resource "sbercloud_dms_kafka_topic" "test" {
  instance_id = sbercloud_dms_kafka_instance.test.id
  name        = "%s"
  partitions  = 20
}`, testAccNetwork_config(rName), rName, rName, password, rName)
}
