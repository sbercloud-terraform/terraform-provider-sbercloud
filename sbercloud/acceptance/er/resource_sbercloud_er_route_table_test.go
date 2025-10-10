package er

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/er/v3/routetables"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getRouteTableResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ErV3Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ER v3 client: %s", err)
	}

	return routetables.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccRouteTable_basic(t *testing.T) {
	var (
		obj routetables.RouteTable

		rName      = "sbercloud_er_route_table.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		baseConfig = testRouteTable_base(name)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRouteTableResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRouteTable_basic_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Create by acc test"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrSet(rName, "is_default_association"),
					resource.TestCheckResourceAttrSet(rName, "is_default_propagation"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testRouteTable_basic_step2(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccRouteTableImportStateFunc(rName),
			},
		},
	})
}

func testAccRouteTableImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, routeTableId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of ER route table is not found in the tfstate", rsName)
		}
		instanceId = rs.Primary.Attributes["instance_id"]
		routeTableId = rs.Primary.ID
		if instanceId == "" || routeTableId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<instance_id>/<id>', but got '%s/%s'",
				instanceId, routeTableId)
		}
		return fmt.Sprintf("%s/%s", instanceId, routeTableId), nil
	}
}

func testRouteTable_base(name string) string {
	bgpAsNum := acctest.RandIntRange(64512, 65534)

	return fmt.Sprintf(`
data "sbercloud_er_availability_zones" "test" {}

%[1]s

resource "sbercloud_er_instance" "test" {
  availability_zones = ["ru-moscow-1a"]
  name               = "%[2]s"
  asn                = %[3]d

  # Enable default routes
  enable_default_propagation = true
  enable_default_association = true
}
`, acceptance.TestVpc(name), name, bgpAsNum)
}

func testRouteTable_basic_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_er_route_table" "test" {
  instance_id = sbercloud_er_instance.test.id
  name        = "%[2]s"
  description = "Create by acc test"

  tags = {
    foo = "bar"
  }
}

resource "sbercloud_er_vpc_attachment" "test" {
  instance_id = sbercloud_er_instance.test.id
  vpc_id      = sbercloud_vpc.test.id
  subnet_id   = sbercloud_vpc_subnet.test.id
  name        = "%[2]s"
}

resource "sbercloud_er_static_route" "test" {
  route_table_id = sbercloud_er_route_table.test.id
  destination    = sbercloud_vpc.test.cidr
  attachment_id  = sbercloud_er_vpc_attachment.test.id
}
`, baseConfig, name)
}

func testRouteTable_basic_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_er_route_table" "test" {
  instance_id = sbercloud_er_instance.test.id
  name        = "%[2]s"

  tags = {
    owner = "terraform"
  }
}

resource "sbercloud_er_vpc_attachment" "test" {
  instance_id = sbercloud_er_instance.test.id
  vpc_id      = sbercloud_vpc.test.id
  subnet_id   = sbercloud_vpc_subnet.test.id
  name        = "%[2]s"
}

resource "sbercloud_er_static_route" "test" {
  route_table_id = sbercloud_er_route_table.test.id
  destination    = sbercloud_vpc.test.cidr
  attachment_id  = sbercloud_er_vpc_attachment.test.id
}
`, baseConfig, name)
}
