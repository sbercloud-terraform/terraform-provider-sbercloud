package sbercloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/volumeattach"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccComputeV2VolumeAttach_basic(t *testing.T) {
	var va volumeattach.VolumeAttachment
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2VolumeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2VolumeAttach_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2VolumeAttachExists("sbercloud_compute_volume_attach.va_1", &va),
				),
			},
			{
				ResourceName:      "sbercloud_compute_volume_attach.va_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeV2VolumeAttach_device(t *testing.T) {
	var va volumeattach.VolumeAttachment
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2VolumeAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2VolumeAttach_device(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2VolumeAttachExists("sbercloud_compute_volume_attach.va_1", &va),
				),
			},
		},
	})
}

func testAccCheckComputeV2VolumeAttachDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	computeClient, err := config.ComputeV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Sbercloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_compute_volume_attach" {
			continue
		}

		instanceId, volumeId, err := parseComputeVolumeAttachmentId(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = volumeattach.Get(computeClient, instanceId, volumeId).Extract()
		if err == nil {
			return fmt.Errorf("Volume attachment still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2VolumeAttachExists(n string, va *volumeattach.VolumeAttachment) resource.TestCheckFunc {
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

		instanceId, volumeId, err := parseComputeVolumeAttachmentId(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := volumeattach.Get(computeClient, instanceId, volumeId).Extract()
		if err != nil {
			return err
		}

		if found.ServerID != instanceId || found.VolumeID != volumeId {
			return fmt.Errorf("VolumeAttach not found")
		}

		*va = *found

		return nil
	}
}

func parseComputeVolumeAttachmentId(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("Unable to determine volume attachment ID")
	}

	instanceId := idParts[0]
	attachmentId := idParts[1]

	return instanceId, attachmentId, nil
}

func testAccCheckComputeV2VolumeAttachDevice(
	va *volumeattach.VolumeAttachment, device string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if va.Device != device {
			return fmt.Errorf("Requested device of volume attachment (%s) does not match: %s",
				device, va.Device)
		}

		return nil
	}
}

func testAccComputeV2VolumeAttach_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_evs_volume" "test" {
  name = "%s"
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  volume_type = "SAS"
  size = 10
}

resource "sbercloud_compute_instance" "instance_1" {
  name = "%s"
  image_id = data.sbercloud_images_image.test.id
  flavor_id = data.sbercloud_compute_flavors.test.ids[0]
  security_groups = ["default"]
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  system_disk_type  = "SSD"
  network {
    uuid = data.sbercloud_vpc_subnet.test.id
  }
}

resource "sbercloud_compute_volume_attach" "va_1" {
  instance_id = sbercloud_compute_instance.instance_1.id
  volume_id = sbercloud_evs_volume.test.id
}
`, testAccCompute_data, rName, rName)
}

func testAccComputeV2VolumeAttach_device(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_evs_volume" "test" {
  name = "%s"
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  volume_type = "SAS"
  size = 10
}

resource "sbercloud_compute_instance" "instance_1" {
  name = "%s"
  image_id = data.sbercloud_images_image.test.id
  flavor_id = data.sbercloud_compute_flavors.test.ids[0]
  security_groups = ["default"]
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  system_disk_type  = "SSD"
  network {
    uuid = data.sbercloud_vpc_subnet.test.id
  }
}

resource "sbercloud_compute_volume_attach" "va_1" {
  instance_id = sbercloud_compute_instance.instance_1.id
  volume_id = sbercloud_evs_volume.test.id
  device = "/dev/vdb"
}
`, testAccCompute_data, rName, rName)
}
