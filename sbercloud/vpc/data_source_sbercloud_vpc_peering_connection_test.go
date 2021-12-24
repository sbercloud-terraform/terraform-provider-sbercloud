package vpc

import (
	"fmt"
	"testing"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcPeeringConnectionDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.sbercloud_vpc_peering_connection.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionDataSource_base(randName),
			},
			{
				Config: testAccVpcPeeringConnectionDataSource_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func TestAccVpcPeeringConnectionDataSource_byVpcId(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.sbercloud_vpc_peering_connection.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionDataSource_base(randName),
			},
			{
				Config: testAccVpcPeeringConnectionDataSource_byVpcId(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func TestAccVpcPeeringConnectionDataSource_byPeerVpcId(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.sbercloud_vpc_peering_connection.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionDataSource_base(randName),
			},
			{
				Config: testAccVpcPeeringConnectionDataSource_byPeerVpcId(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func TestAccVpcPeeringConnectionDataSource_byVpcIds(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.sbercloud_vpc_peering_connection.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConnectionDataSource_base(randName),
			},
			{
				Config: testAccVpcPeeringConnectionDataSource_byVpcIds(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccVpcPeeringConnectionDataSource_base(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc" "vpc_1" {
  name = "%s_1"
  cidr = "172.16.0.0/20"
}

resource "sbercloud_vpc" "vpc_2" {
  name = "%s_2"
  cidr = "172.16.128.0/20"
}

resource "sbercloud_vpc_peering_connection" "test" {
  name        = "%s"
  vpc_id      = sbercloud_vpc.vpc_1.id
  peer_vpc_id = sbercloud_vpc.vpc_2.id
}
`, rName, rName, rName)
}

func testAccVpcPeeringConnectionDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_peering_connection" "test" {
  id = sbercloud_vpc_peering_connection.test.id
}
`, testAccVpcPeeringConnectionDataSource_base(rName))
}

func testAccVpcPeeringConnectionDataSource_byVpcId(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_peering_connection" "test" {
	vpc_id = sbercloud_vpc_peering_connection.test.vpc_id
}
`, testAccVpcPeeringConnectionDataSource_base(rName))
}

func testAccVpcPeeringConnectionDataSource_byPeerVpcId(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_peering_connection" "test" {
	peer_vpc_id = sbercloud_vpc_peering_connection.test.peer_vpc_id
}
`, testAccVpcPeeringConnectionDataSource_base(rName))
}

func testAccVpcPeeringConnectionDataSource_byVpcIds(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc_peering_connection" "test" {
  vpc_id      = sbercloud_vpc_peering_connection.test.vpc_id
  peer_vpc_id = sbercloud_vpc_peering_connection.test.peer_vpc_id
}
`, testAccVpcPeeringConnectionDataSource_base(rName))
}
