package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/peerings"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func getPeeringConnectionResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.NetworkingV2Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SberCloud Network client: %s", err)
	}
	return peerings.Get(c, state.Primary.ID).Extract()
}

func TestAccVpcPeeringConnection_basic(t *testing.T) {
	var peering peerings.Peering

	randName := acceptance.RandomAccResourceName()
	updateName := randName + "_update"
	resourceName := "sbercloud_vpc_peering_connection.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&peering,
		getPeeringConnectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnection_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "vpc_id",
						"${sbercloud_vpc.test1.id}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "peer_vpc_id",
						"${sbercloud_vpc.test2.id}"),
				),
			},
			{
				Config: testAccVpcPeeringConnection_basic(updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "vpc_id",
						"${sbercloud_vpc.test1.id}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "peer_vpc_id",
						"${sbercloud_vpc.test2.id}"),
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

func testAccVpcPeeringConnection_basic(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "test1" {
  name = "%s_1"
  cidr = "172.16.0.0/20"
}

resource "sbercloud_vpc" "test2" {
  name = "%s_2"
  cidr = "172.16.128.0/20"
}

resource "sbercloud_vpc_peering_connection" "test" {
  name        = "%s"
  vpc_id      = sbercloud_vpc.test1.id
  peer_vpc_id = sbercloud_vpc.test2.id
}
`, rName, rName, rName)
}
