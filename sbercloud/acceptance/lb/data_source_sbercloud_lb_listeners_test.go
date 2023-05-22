package lb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccDatasourceListeners_basic(t *testing.T) {
	var (
		rName            = acceptance.RandomAccResourceNameWithDash()
		dcByName         = acceptance.InitDataSourceCheck("data.cloud_lb_listeners.by_name")
		dcByProtocol     = acceptance.InitDataSourceCheck("data.sbercloud_lb_listeners.by_protocol")
		dcByProtocolPort = acceptance.InitDataSourceCheck("data.sbercloud_lb_listeners.by_protocol_port")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceListeners_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_query_result_validation", "true"),
					resource.TestCheckResourceAttrSet("data.sbercloud_lb_listeners.by_name",
						"listeners.0.name"),
					resource.TestCheckResourceAttrSet("data.sbercloud_lb_listeners.by_name",
						"listeners.0.protocol"),
					resource.TestCheckResourceAttrSet("data.sbercloud_lb_listeners.by_name",
						"listeners.0.protocol_port"),
					resource.TestCheckResourceAttrSet("data.sbercloud_lb_listeners.by_name",
						"listeners.0.connection_limit"),
					resource.TestCheckResourceAttrSet("data.sbercloud_lb_listeners.by_name",
						"listeners.0.http2_enable"),
					resource.TestCheckResourceAttr("data.sbercloud_lb_listeners.by_name",
						"listeners.0.loadbalancers.#", "1"),
					dcByProtocol.CheckResourceExists(),
					resource.TestCheckOutput("protocol_query_result_validation", "true"),
					dcByProtocolPort.CheckResourceExists(),
					resource.TestCheckOutput("protocol_port_query_result_validation", "true"),
				),
			},
		},
	})
}

func testAccDatasourceListeners_base(rName string) string {
	rCidr := acceptance.RandomCidr()
	return fmt.Sprintf(`
variable "listener_configuration" {
  type = list(object({
    protocol_port = number
    protocol      = string
  }))
  default = [
    {protocol_port = 306, protocol = "TCP"},
    {protocol_port = 406, protocol = "UDP"},
    {protocol_port = 506, protocol = "HTTP"},
  ]
}

resource "sbercloud_vpc" "test" {
  name = "%[1]s"
  cidr = "%[2]s"
}

resource "sbercloud_vpc_subnet" "test" {
  vpc_id = sbercloud_vpc.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(sbercloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(sbercloud_vpc.test.cidr, 4, 1), 1)
}

resource "sbercloud_lb_loadbalancer" "test" {
  name          = "%[1]s"
  vip_subnet_id = sbercloud_vpc_subnet.test.ipv4_subnet_id
}

resource "sbercloud_lb_listener" "test" {
  count = length(var.listener_configuration)

  loadbalancer_id = sbercloud_lb_loadbalancer.test.id

  name          = "%[1]s-${count.index}"
  protocol      = var.listener_configuration[count.index]["protocol"]
  protocol_port = var.listener_configuration[count.index]["protocol_port"]
}
`, rName, rCidr)
}

func testAccDatasourceListeners_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_lb_listeners" "by_name" {
  depends_on = [sbercloud_lb_listener.test]

  name = sbercloud_lb_listener.test[0].name
}

data "sbercloud_lb_listeners" "by_protocol" {
  depends_on = [sbercloud_lb_listener.test]

  protocol = sbercloud_lb_listener.test[1].protocol
}

data "sbercloud_lb_listeners" "by_protocol_port" {
  depends_on = [sbercloud_lb_listener.test]

  protocol_port = sbercloud_lb_listener.test[2].protocol_port
}

output "name_query_result_validation" {
  value = contains(data.sbercloud_lb_listeners.by_name.listeners[*].id,
  sbercloud_lb_listener.test[0].id) && !contains(data.sbercloud_lb_listeners.by_name.listeners[*].id,
  sbercloud_lb_listener.test[1].id) && !contains(data.sbercloud_lb_listeners.by_name.listeners[*].id,
  sbercloud_lb_listener.test[2].id)
}

output "protocol_query_result_validation" {
  value = contains(data.sbercloud_lb_listeners.by_protocol.listeners[*].id,
  sbercloud_lb_listener.test[1].id) && !contains(data.sbercloud_lb_listeners.by_protocol.listeners[*].id,
  sbercloud_lb_listener.test[0].id) && !contains(data.sbercloud_lb_listeners.by_protocol.listeners[*].id,
  sbercloud_lb_listener.test[2].id)
}

output "protocol_port_query_result_validation" {
  value = contains(data.sbercloud_lb_listeners.by_protocol_port.listeners[*].id,
  sbercloud_lb_listener.test[2].id) && !contains(data.sbercloud_lb_listeners.by_protocol_port.listeners[*].id,
  sbercloud_lb_listener.test[0].id) && !contains(data.sbercloud_lb_listeners.by_protocol_port.listeners[*].id,
  sbercloud_lb_listener.test[1].id)
}
`, testAccDatasourceListeners_base(rName))
}
