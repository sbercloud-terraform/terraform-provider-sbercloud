package rocketmq

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceDmsRocketMQBroker_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.sbercloud_dms_rocketmq_broker.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDmsRocketMQBroker_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "brokers.0", "broker-0"),
				),
			},
		},
	})
}

func testAccDatasourceDmsRocketMQBroker_base(name string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = "4.8.0"
  storage_space     = 600
  vpc_id            = sbercloud_vpc.test.id
  subnet_id         = sbercloud_vpc_subnet.test.id
  security_group_id = sbercloud_networking_secgroup.test.id

  availability_zones = [
    data.sbercloud_availability_zones.test.names[0]
  ]

  flavor_id         = "c6.4u8g.cluster"
  storage_spec_code = "dms.physical.storage.high.v2"
  broker_num        = 1
}
`, acceptance.TestBaseNetwork(name), name)
}

func testAccDatasourceDmsRocketMQBroker_basic(name string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_dms_rocketmq_broker" "test" {
  instance_id = sbercloud_dms_rocketmq_instance.test.id
}
`, testAccDatasourceDmsRocketMQBroker_base(name))
}
