package vpcep

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/vpcep/v1/services"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccVPCEPService_Basic(t *testing.T) {
	var service services.Service

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(4)) //acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_vpcep_service.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&service,
		getVpcepServiceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEPService_Basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "approval", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "enable_policy", "false"),
					resource.TestCheckResourceAttr(resourceName, "server_type", "VM"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "interface"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.service_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.terminal_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "organization_permissions.#", "2"),
				),
			},
			{
				Config: testAccVPCEPService_Update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-"+rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "approval", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc-update"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.service_port", "8088"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.terminal_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "permissions.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "organization_permissions.0", "organizations:orgPath::*"),
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

func TestAccVPCEPService_enablePolicy(t *testing.T) {
	var service services.Service

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(4)) //acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_vpcep_service.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&service,
		getVpcepServiceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEPService_enablePolicy(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "approval", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "enable_policy", "true"),
					resource.TestCheckResourceAttr(resourceName, "server_type", "VM"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "interface"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.service_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.terminal_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "2"),
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

func getVpcepServiceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	vpcepClient, err := conf.VPCEPClient(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPCEP client: %s", err)
	}

	return services.Get(vpcepClient, state.Primary.ID).Extract()
}

func testAccVPCEPService_Precondition(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_vpc" "myvpc" {
  name = "vpc-default"
}

resource "sbercloud_compute_instance" "ecs" {
  name               = "%s"
  image_id           = data.sbercloud_images_image.test.id
  flavor_id          = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids = [data.sbercloud_networking_secgroup.test.id]
  availability_zone  = data.sbercloud_availability_zones.test.names[0]
  system_disk_type = "SSD"

  network {
    uuid = data.sbercloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName)
}

func testAccVPCEPService_Basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_vpcep_service" "test" {
  name                     = "%s"
  server_type              = "VM"
  vpc_id                   = data.sbercloud_vpc.myvpc.id
  port_id                  = sbercloud_compute_instance.ecs.network[0].port
  approval                 = false
  description              = "test description"
  permissions              = ["iam:domain::*"]
  organization_permissions = ["organizations:orgPath::*"]

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
  tags = {
    owner = "tf-acc"
  }
}
`, testAccVPCEPService_Precondition(rName), rName)
}

func testAccVPCEPService_Update(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_vpcep_service" "test" {
  name                     = "tf-%s"
  server_type              = "VM"
  vpc_id                   = data.sbercloud_vpc.myvpc.id
  port_id                  = sbercloud_compute_instance.ecs.network[0].port
  approval                 = true
  description              = "test description update"
  permissions              = ["*"]
  organization_permissions = ["organizations:orgPath::*"]

  port_mapping {
    service_port  = 8088
    terminal_port = 80
  }
  tags = {
    owner = "tf-acc-update"
  }
}
`, testAccVPCEPService_Precondition(rName), rName)
}

func testAccVPCEPService_enablePolicy(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_vpcep_service" "test" {
  name          = "%s"
  server_type   = "VM"
  vpc_id        = data.sbercloud_vpc.myvpc.id
  port_id       = sbercloud_compute_instance.ecs.network[0].port
  approval      = false
  description   = "test description"
  enable_policy = true
  permissions   = ["iam:domain::6e9dfd5d1124e8d8498dce894923a0dd", "iam:domain::6e9dfd5d1124e8d8498dce894923a0de"]

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}
`, testAccVPCEPService_Precondition(rName), rName)
}

const testAccCompute_data = `
data "sbercloud_availability_zones" "test" {}

data "sbercloud_compute_flavors" "test" {
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "sbercloud_images_image" "test" {
  name        = "Ubuntu 20.04 server 64bit"
  most_recent = true
}

data "sbercloud_networking_secgroup" "test" {
  name = "default"
}
`
