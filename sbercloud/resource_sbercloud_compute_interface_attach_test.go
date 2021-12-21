package sbercloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccComputeV2InterfaceAttach_Basic(t *testing.T) {
	var ai attachinterfaces.Interface
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InterfaceAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InterfaceAttach_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InterfaceAttachExists("sbercloud_compute_interface_attach.ai_1", &ai),
					testAccCheckComputeV2InterfaceAttachIP(&ai, "192.168.0.199"),
				),
			},
		},
	})
}

func testAccCheckComputeV2InterfaceAttachDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	computeClient, err := config.ComputeV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Sbercloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_compute_interface_attach" {
			continue
		}

		instanceId, portId, err := computeInterfaceAttachV2ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = attachinterfaces.Get(computeClient, instanceId, portId).Extract()
		if err == nil {
			return fmt.Errorf("Volume attachment still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2InterfaceAttachExists(n string, ai *attachinterfaces.Interface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		computeClient, err := config.ComputeV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Sbercloud compute client: %s", err)
		}

		instanceId, portId, err := computeInterfaceAttachV2ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := attachinterfaces.Get(computeClient, instanceId, portId).Extract()
		if err != nil {
			return err
		}

		//if found.instanceID != instanceID || found.PortID != portId {
		if found.PortID != portId {
			return fmt.Errorf("InterfaceAttach not found")
		}

		*ai = *found

		return nil
	}
}

func computeInterfaceAttachV2ParseID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("Unable to determine huaweicloud_compute_interface_attach_v2 %s ID", id)
	}

	instanceId := idParts[0]
	attachmentId := idParts[1]

	return instanceId, attachmentId, nil
}

func testAccCheckComputeV2InterfaceAttachIP(
	ai *attachinterfaces.Interface, ip string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range ai.FixedIPs {
			if i.IPAddress == ip {
				return nil
			}
		}
		return fmt.Errorf("Requested ip (%s) does not exist on port", ip)

	}
}

func testAccComputeV2InterfaceAttach_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_instance" "instance_1" {
  name = "%s"
  image_id          = data.sbercloud_images_image.test.id
  flavor_id         = data.sbercloud_compute_flavors.test.ids[0]
  security_groups = ["default"]
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  system_disk_type  = "SSD"
  network {
    uuid = data.sbercloud_vpc_subnet.test.id
  }
}

resource "sbercloud_compute_interface_attach" "ai_1" {
  instance_id = sbercloud_compute_instance.instance_1.id
  network_id =  data.sbercloud_vpc_subnet.test.id
  fixed_ip = "192.168.0.199"
}
`, testAccCompute_data, rName)
}
