package vpc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccVpcPeeringConnectionAccepter_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcPeeringConnectionAccepterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccVpcPeeringConnectionAccepter_basic(randName), //TODO: Research why normal scenario with peer tenant id is not working in acceptance tests
				ExpectError: regexp.MustCompile(`VPC peering action not permitted: Can not accept/reject peering request not in PENDING_ACCEPTANCE state.`),
			},
		},
	})
}

func testAccCheckVpcPeeringConnectionAccepterDestroy(s *terraform.State) error {
	// We don't destroy the underlying VPC Peering Connection.
	return nil
}

func testAccVpcPeeringConnectionAccepter_basic(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "test1" {
  name = "%s_1"
  cidr = "192.168.0.0/20"
}

resource "sbercloud_vpc" "test2" {
  name = "%s_2"
  cidr = "192.168.128.0/20"
}

resource "sbercloud_vpc_peering_connection" "test" {
  name        = "%s"
  vpc_id      = sbercloud_vpc.test1.id
  peer_vpc_id = sbercloud_vpc.test2.id
}

resource "sbercloud_vpc_peering_connection_accepter" "test" {
  vpc_peering_connection_id = sbercloud_vpc_peering_connection.test.id

  accept = true
}
`, rName, rName, rName)
}
