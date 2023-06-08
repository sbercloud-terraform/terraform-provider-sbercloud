package sbercloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/services/ces/alarmrule"
	"testing"
)

func TestAccCESAlarmRule_basic(t *testing.T) {
	var ar alarmrule.AlarmRule
	rName := fmt.Sprintf("tf-acc-%s", acctest.RandString(5))
	resourceName := "sbercloud_ces_alarmrule.alarmrule_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCESAlarmRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testCESAlarmRuleExists(resourceName, &ar),
					resource.TestCheckResourceAttr(resourceName, "alarm_name", fmt.Sprintf("rule-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_action_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_level", "2"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.value", "6"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCESAlarmRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.CesV1Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud ces client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_ces_alarmrule" {
			continue
		}

		id := rs.Primary.ID
		_, err := alarmrule.Get(networkingClient, id).Extract()
		if err == nil {
			return fmt.Errorf("Alarm rule still exists")
		}
	}

	return nil
}

func testCESAlarmRuleExists(n string, ar *alarmrule.AlarmRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		networkingClient, err := config.CesV1Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud ces client: %s", err)
		}

		id := rs.Primary.ID
		found, err := alarmrule.Get(networkingClient, id).Extract()
		if err != nil {
			return err
		}

		*ar = *found

		return nil
	}
}

func testCESAlarmRule_base(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_availability_zones" "test" {}

data "sbercloud_compute_flavors" "test" {
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "sbercloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "sbercloud_compute_instance" "vm_1" {
  name              = "ecs-%s"
  image_id          = data.sbercloud_images_image.test.id
  flavor_id         = data.sbercloud_compute_flavors.test.ids[0]
  security_groups   = ["default"]
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  system_disk_type  = "SSD"

  network {
    uuid = data.sbercloud_vpc_subnet.test.id
  }
}

resource "sbercloud_smn_topic" "topic_1" {
  name         = "smn-%s"
  display_name = "The display name of smn topic"
}
`, rName, rName)
}

func testCESAlarmRule_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_ces_alarmrule" "alarmrule_1" {
  alarm_name           = "rule-%s"
  alarm_action_enabled = true

  metric {
    namespace   = "SYS.ECS"
    metric_name = "network_outgoing_bytes_rate_inband"

    dimensions {
      name  = "instance_id"
      value = sbercloud_compute_instance.vm_1.id
    }
  }

  condition  {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6
    unit                = "B/s"
    count               = 1
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      sbercloud_smn_topic.topic_1.topic_urn
    ]
  }
}
`, testCESAlarmRule_base(rName), rName)
}
