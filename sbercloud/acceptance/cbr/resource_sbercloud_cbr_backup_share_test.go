package cbr

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cbr/v3/members"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getBackupShareResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CbrV3Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CBR v3 client: %s", err)
	}
	var (
		backupId = state.Primary.ID
		opts     = members.ListOpts{
			BackupId: backupId,
		}
	)
	memberList, err := members.List(client, opts)
	if len(memberList) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}
	return memberList, err
}

func TestAccBackupShare_basic(t *testing.T) {
	var (
		memberList   []members.Member
		name         = acceptance.RandomAccResourceName()
		resourceName = "sbercloud_cbr_backup_share.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&memberList,
		getBackupShareResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDestProjectIds(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBackupShare_basic(name, acceptance.SBC_DEST_PROJECT_ID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "backup_id"),
					resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "members.0.dest_project_id", acceptance.SBC_DEST_PROJECT_ID),
				),
			},
			{
				Config: testAccBackupShare_basic(name, acceptance.SBC_DEST_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "backup_id"),
					resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "members.0.dest_project_id", acceptance.SBC_DEST_PROJECT_ID_TEST),
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

func testAccBackupShare_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.sbercloud_images_image.test.id
  flavor_id          = data.sbercloud_compute_flavors.test.ids[0]
  security_group_ids = [sbercloud_networking_secgroup.test.id]
  availability_zone  = data.sbercloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = sbercloud_vpc_subnet.test.id
  }
}

resource "sbercloud_cbr_vault" "test" {
  name             = "%[2]s"
  type             = "server"
  consistent_level = "crash_consistent"
  protection_type  = "backup"
  size             = 10

  resources {
    server_id = sbercloud_compute_instance.test.id
  }
}

resource "sbercloud_cbr_checkpoint" "test" {
  vault_id = sbercloud_cbr_vault.test.id
  name     = "%[2]s"

  backups {
    type        = "OS::Nova::Server"
    resource_id = sbercloud_compute_instance.test.id
  }
}
`, acceptance.TestBaseComputeResources(name), name)
}

func testAccBackupShare_basic(name, destProjectId string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cbr_backup_share" "test" {
  backup_id = try(tolist(sbercloud_cbr_checkpoint.test.backups)[0].id, "")

  members {
    # Different user (ID) in the same region.
    dest_project_id = "%[2]s"
  }
}
`, testAccBackupShare_base(name), destProjectId)
}
