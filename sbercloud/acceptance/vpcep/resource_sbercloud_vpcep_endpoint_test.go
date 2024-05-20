package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/vpcep/v1/endpoints"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccVPCEndpoint_Basic(t *testing.T) {
	var endpoint endpoints.Endpoint

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_vpcep_endpoint.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&endpoint,
		getVpcepEndpointResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpoint_Basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "accepted"),
					resource.TestCheckResourceAttr(resourceName, "enable_dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "interface"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc"),
					resource.TestCheckResourceAttr(resourceName, "enable_whitelist", "true"),
					resource.TestCheckResourceAttr(resourceName, "whitelist.0", "192.168.0.0/24"),
					resource.TestCheckResourceAttrSet(resourceName, "service_name"),
					resource.TestCheckResourceAttrSet(resourceName, "private_domain_name"),
				),
			},
			{
				Config: testAccVPCEndpoint_Update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "accepted"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc-update"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "enable_whitelist", "false"),
					resource.TestCheckResourceAttr(resourceName, "whitelist.#", "0"),
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

func TestAccVPCEndpoint_Public(t *testing.T) {
	var endpoint endpoints.Endpoint
	resourceName := "sbercloud_vpcep_endpoint.myendpoint"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&endpoint,
		getVpcepEndpointResourceFunc,
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpointPublic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "accepted"),
					resource.TestCheckResourceAttr(resourceName, "enable_dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_whitelist", "true"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "interface"),
					resource.TestCheckResourceAttr(resourceName, "whitelist.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "service_name"),
					resource.TestCheckResourceAttrSet(resourceName, "private_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
				),
			},
		},
	})
}

func getVpcepEndpointResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	vpcepClient, err := conf.VPCEPClient(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPCEP client: %s", err)
	}

	return endpoints.Get(vpcepClient, state.Primary.ID).Extract()
}

func testAccVPCEndpoint_Precondition(rName string) string {
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

  network {
    uuid = data.sbercloud_vpc_subnet.test.id
  }
}

resource "sbercloud_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = data.sbercloud_vpc.myvpc.id
  port_id     = sbercloud_compute_instance.ecs.network[0].port
  approval    = false

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
  tags = {
    owner = "tf-acc"
  }
}
`, testAccCompute_data, rName, rName)
}

func testAccVPCEndpoint_Basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_vpcep_endpoint" "test" {
  service_id       = sbercloud_vpcep_service.test.id
  vpc_id           = data.sbercloud_vpc.myvpc.id
  network_id       = data.sbercloud_vpc_subnet.test.id
  enable_dns       = true
  description      = "test description"
  enable_whitelist = true
  whitelist        = ["192.168.0.0/24"]

  tags = {
    owner = "tf-acc"
  }
}
`, testAccVPCEndpoint_Precondition(rName))
}

func testAccVPCEndpoint_Update(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_vpcep_endpoint" "test" {
  service_id       = sbercloud_vpcep_service.test.id
  vpc_id           = data.sbercloud_vpc.myvpc.id
  network_id       = data.sbercloud_vpc_subnet.test.id
  enable_dns       = true
  description      = "test description"
  enable_whitelist = false

  tags = {
    owner = "tf-acc-update"
    foo   = "bar"
  }
}
`, testAccVPCEndpoint_Precondition(rName))
}

var testAccVPCEndpointPublic = `
data "sbercloud_vpc" "myvpc" {
  name = "vpc-default"
}

data "sbercloud_vpc_subnet" "mynet" {
  vpc_id = data.sbercloud_vpc.myvpc.id
  name   = "subnet-default"
}

data "sbercloud_vpcep_public_services" "cloud_service" {
  service_name = "dis"
}

resource "sbercloud_vpcep_endpoint" "myendpoint" {
  service_id       = data.sbercloud_vpcep_public_services.cloud_service.services[0].id
  vpc_id           = data.sbercloud_vpc.myvpc.id
  network_id       = data.sbercloud_vpc_subnet.mynet.id
  enable_dns       = true
  enable_whitelist = true
  whitelist        = ["192.168.0.0/24", "10.10.10.10"]
}
`
