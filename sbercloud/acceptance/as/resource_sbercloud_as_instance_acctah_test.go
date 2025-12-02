package as

import (
	"fmt"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/instances"
)

func getASInstanceAttachResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.SBC_REGION_NAME
	client, err := cfg.AutoscalingV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating autoscaling client: %s", err)
	}

	groupID := state.Primary.Attributes["scaling_group_id"]
	instanceID := state.Primary.Attributes["instance_id"]
	page, err := instances.List(client, groupID, nil).AllPages()
	if err != nil {
		return nil, err
	}

	allInstances, err := page.(instances.InstancePage).Extract()
	if err != nil {
		return nil, fmt.Errorf("failed to fetching instances in AS group %s: %s", groupID, err)
	}

	for _, ins := range allInstances {
		if ins.ID == instanceID {
			return &ins, nil
		}
	}

	return nil, fmt.Errorf("can not find the instance %s in AS group %s", instanceID, groupID)
}

// The current test case does not have the test field `append_instance`, because testing this field may cause the
// automatically added ECS instance resources to remain.
func TestAccASInstanceAttach_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "sbercloud_as_instance_attach.test0"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getASInstanceAttachResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testASInstanceAttach_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "scaling_group_id", "sbercloud_as_group.acc_as_group", "id"),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "sbercloud_compute_instance.test.0", "id"),
					resource.TestCheckResourceAttr(rName, "protected", "true"),
					resource.TestCheckResourceAttr(rName, "standby", "true"),
					resource.TestCheckResourceAttr(rName, "status", "STANDBY"),
				),
			},
			{
				Config: testASInstanceAttach_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "protected", "true"),
					resource.TestCheckResourceAttr(rName, "standby", "false"),
					resource.TestCheckResourceAttr(rName, "status", "INSERVICE"),
				),
			},
			{
				Config: testASInstanceAttach_step3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "protected", "false"),
					resource.TestCheckResourceAttr(rName, "standby", "false"),
					resource.TestCheckResourceAttr(rName, "status", "INSERVICE"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"append_instance"},
			},
		},
	})
}

func testASInstanceAttach_base(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_compute_instance" "test" {
  count              = 2
  name               = "%s-${count.index}"
  description        = "instance for AS attach"
  image_id           = data.sbercloud_images_image.test.id
  flavor_id          = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids = [sbercloud_networking_secgroup.test.id]

  system_disk_type = "SAS"

  network {
    uuid = sbercloud_vpc_subnet.test.id
  }
}
`, testASGroup_basic(name), name)
}

func testASInstanceAttach_step1(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_as_instance_attach" "test0" {
  scaling_group_id = sbercloud_as_group.acc_as_group.id
  instance_id      = sbercloud_compute_instance.test[0].id
  protected        = true
  standby          = true
}

resource "sbercloud_as_instance_attach" "test1" {
  scaling_group_id = sbercloud_as_group.acc_as_group.id
  instance_id      = sbercloud_compute_instance.test[1].id
  protected        = false
  standby          = false
}
`, testASInstanceAttach_base(name))
}

func testASInstanceAttach_step2(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_as_instance_attach" "test0" {
  scaling_group_id = sbercloud_as_group.acc_as_group.id
  instance_id      = sbercloud_compute_instance.test[0].id
  protected        = true
  standby          = false
}

resource "sbercloud_as_instance_attach" "test1" {
  scaling_group_id = sbercloud_as_group.acc_as_group.id
  instance_id      = sbercloud_compute_instance.test[1].id
  protected        = false
  standby          = true
}
`, testASInstanceAttach_base(name))
}

func testASInstanceAttach_step3(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_as_instance_attach" "test0" {
  scaling_group_id = sbercloud_as_group.acc_as_group.id
  instance_id      = sbercloud_compute_instance.test[0].id
  protected        = false
  standby          = false
}

# When the instance status is standby, closing protected is not supported.
resource "sbercloud_as_instance_attach" "test1" {
  scaling_group_id = sbercloud_as_group.acc_as_group.id
  instance_id      = sbercloud_compute_instance.test[1].id
  protected        = false
  standby          = true
}
`, testASInstanceAttach_base(name))
}
