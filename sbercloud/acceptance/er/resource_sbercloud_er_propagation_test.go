package er

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/er/v3/propagations"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/er"
)

func getPropagationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ErV3Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ER v3 client: %s", err)
	}

	return er.QueryPropagationById(client, state.Primary.Attributes["instance_id"],
		state.Primary.Attributes["route_table_id"], state.Primary.ID)
}

func TestAccPropagation_basic(t *testing.T) {
	var (
		obj propagations.Propagation

		rName = "sbercloud_er_propagation.test"
		name  = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPropagationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPropagation_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"sbercloud_er_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "route_table_id",
						"sbercloud_er_route_table.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "attachment_id",
						"sbercloud_er_vpc_attachment.test", "id"),
					resource.TestCheckResourceAttr(rName, "attachment_type", "vpc"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPropagationImportStateFunc(rName),
			},
		},
	})
}

func testAccPropagationImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, routeTableId, propagationId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of ER propagation is not found in the tfstate", rsName)
		}
		instanceId = rs.Primary.Attributes["instance_id"]
		routeTableId = rs.Primary.Attributes["route_table_id"]
		propagationId = rs.Primary.ID
		if instanceId == "" || routeTableId == "" || propagationId == "" {
			return "", fmt.Errorf("some import IDs are missing, want "+
				"'<instance_id>/<route_table_id>/<id>', but got '%s/%s/%s'",
				instanceId, routeTableId, propagationId)
		}
		return fmt.Sprintf("%s/%s/%s", instanceId, routeTableId, propagationId), nil
	}
}

func testAccPropagation_base(name string) string {
	bgpAsNum := acctest.RandIntRange(64512, 65534)

	return fmt.Sprintf(`
data "sbercloud_er_availability_zones" "test" {}

resource "sbercloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_subnet" "test" {
  vpc_id = sbercloud_vpc.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(sbercloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(sbercloud_vpc.test.cidr, 4, 1), 1)
}

resource "sbercloud_er_instance" "test" {
  availability_zones = slice(data.sbercloud_er_availability_zones.test.names, 0, 1)

  name = "%[1]s"
  asn  = %[2]d
}

resource "sbercloud_er_vpc_attachment" "test" {
  instance_id = sbercloud_er_instance.test.id
  vpc_id      = sbercloud_vpc.test.id
  subnet_id   = sbercloud_vpc_subnet.test.id

  name                   = "%[1]s"
  auto_create_vpc_routes = true
}

resource "sbercloud_er_route_table" "test" {
  instance_id = sbercloud_er_instance.test.id

  name = "%[1]s"
}
`, name, bgpAsNum)
}

func testAccPropagation_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_er_propagation" "test" {
  instance_id    = sbercloud_er_instance.test.id
  route_table_id = sbercloud_er_route_table.test.id
  attachment_id  = sbercloud_er_vpc_attachment.test.id
}
`, testAccPropagation_base(name))
}
