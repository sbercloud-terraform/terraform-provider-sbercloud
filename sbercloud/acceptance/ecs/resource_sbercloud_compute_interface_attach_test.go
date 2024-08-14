package ecs

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPrivateCAResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	computeClient, err := conf.NewServiceClient("ecs", state.Primary.Attributes["region"])
	if err != nil {
		return nil, fmt.Errorf("error creating compute client: %s", err)
	}

	listNicsHttpUrl := "v1/{project_id}/cloudservers/{server_id}/os-interface"
	listNicsPath := computeClient.Endpoint + listNicsHttpUrl
	listNicsPath = strings.ReplaceAll(listNicsPath, "{project_id}", computeClient.ProjectID)
	listNicsPath = strings.ReplaceAll(listNicsPath, "{server_id}", state.Primary.Attributes["instance_id"])
	listNicsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listNicsResp, err := computeClient.Request("GET", listNicsPath, &listNicsOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ECS NICs: %s", err)
	}
	listNicsRespBody, err := utils.FlattenResponse(listNicsResp)
	if err != nil {
		return nil, fmt.Errorf("error prasing ECS NICs: %s", err)
	}

	jsonPaths := fmt.Sprintf("interfaceAttachments[?port_id=='%s']|[0]", state.Primary.ID)
	nic, err := jmespath.Search(jsonPaths, listNicsRespBody)
	if err != nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return nic, nil
}

func TestAccComputeInterfaceAttach_Basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_compute_interface_attach.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPrivateCAResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInterfaceAttach_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "source_dest_check", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_ids.0",
						"sbercloud_networking_secgroup.test", "id"),
				),
			},
			{
				Config: testAccComputeInterfaceAttach_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "source_dest_check", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_ids.0",
						"sbercloud_networking_secgroup.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testComputeInterfaceAttachImportState(resourceName),
			},
		},
	})
}

func testAccComputeInterfaceAttachBase(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_availability_zones" "test" {}

resource "sbercloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_subnet" "test" {
  vpc_id     = sbercloud_vpc.test.id
  name       = "%[1]s"
  cidr       = cidrsubnet(sbercloud_vpc.test.cidr, 4, 0)
  gateway_ip = cidrhost(cidrsubnet(sbercloud_vpc.test.cidr, 4, 0), 1)
}

resource "sbercloud_networking_secgroup" "test" {
  name = "%[1]s"
}

data "sbercloud_compute_flavors" "test" {
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "sbercloud_images_images" "test" {
  flavor_id = data.sbercloud_compute_flavors.test.ids[0]

  os         = "Ubuntu"
  visibility = "public"
}

resource "sbercloud_compute_instance" "test" {
  name               = "%[1]s"
  image_id           = data.sbercloud_images_images.test.images[0].id
  flavor_id          = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids = [sbercloud_networking_secgroup.test.id]
  availability_zone  = data.sbercloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = sbercloud_vpc_subnet.test.id
  }
}
`, rName)
}

func testAccComputeInterfaceAttach_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_interface_attach" "test" {
  instance_id        = sbercloud_compute_instance.test.id
  network_id         = sbercloud_vpc_subnet.test.id
  security_group_ids = [sbercloud_networking_secgroup.test.id]
  source_dest_check  = false
}
`, testAccComputeInterfaceAttachBase(rName))
}

func testAccComputeInterfaceAttach_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_interface_attach" "test" {
  instance_id        = sbercloud_compute_instance.test.id
  network_id         = sbercloud_vpc_subnet.test.id
  security_group_ids = [sbercloud_networking_secgroup.test.id]
}
`, testAccComputeInterfaceAttachBase(rName))
}

// testComputeInterfaceAttachImportState use to return an id with format <instance_id>/<port_id>
func testComputeInterfaceAttachImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		serverID := rs.Primary.Attributes["instance_id"]
		if serverID == "" {
			return "", fmt.Errorf("attribute `instance_id` of the resource (%s) not found", name)
		}

		return serverID + "/" + rs.Primary.ID, nil
	}
}
