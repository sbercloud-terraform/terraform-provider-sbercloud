package dcs

import (
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDcsRestoreV1_basic(t *testing.T) {
	resourceName := "sbercloud_dcs_restore.test"

	projectId, instanceId, backupId, err := CreateInstanceAndBackup()
	if err != nil {
		t.Errorf("instance and backup creating error")
	}
	//projectId := "0f5181caba0024e72f89c0045e707b91"
	//instanceId := "578655e4-5846-4f1b-bfe4-4938ebc7e19e"
	//backupId := "ed466175-3a5d-42d0-90b0-bb3ec29e1465"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Restore_basic(projectId, instanceId, backupId),
				Check: resource.ComposeTestCheckFunc(
					//resource.TestCheckResourceAttr(resourceName, "project_id", projectId),
					resource.TestCheckResourceAttr(resourceName, "instance_id", instanceId),
					resource.TestCheckResourceAttr(resourceName, "backup_id", backupId),
					resource.TestCheckResourceAttr(resourceName, "remark", "restore instance"),
				),
			},
		},
	})
}

func testAccDcsV1Restore_basic(projectId, instanceId, backupId string) string {
	return `
data "sbercloud_dcs_flavors" "single_flavors" {
  cache_mode = "ha"
  capacity   = 1
}


resource "sbercloud_vpc" "vpc" {
  name = "vpc_name"
  cidr = "192.168.0.0/16"
}

resource "sbercloud_vpc_subnet" "subnet" {
  name       = "subnet_name"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = sbercloud_vpc.vpc.id
}

data "sbercloud_identity_projects" "test" {
  name = "ru-moscow-1"
}

resource "sbercloud_dcs_instance" "instance_1" {
  name               = "redis_single_instance"
  engine             = "Redis"
  engine_version     = "5.0"
  capacity           = data.sbercloud_dcs_flavors.single_flavors.capacity
  flavor             = data.sbercloud_dcs_flavors.single_flavors.flavors[0].name
  availability_zones = ["ru-moscow-1a", "ru-moscow-1b"]
  password           = "YourPassword@123"
  vpc_id             = sbercloud_vpc.vpc.id
  subnet_id          = sbercloud_vpc_subnet.subnet.id
}

resource "sbercloud_dcs_backup" "test1"{
  instance_id = sbercloud_dcs_instance.instance_1.id
  description   = "test DCS backup remark"
  backup_format = "rdb"
  depends_on = [
    sbercloud_dcs_instance.instance_1
  ]
}

resource "sbercloud_dcs_restore" "test" {
  project_id  = data.sbercloud_identity_projects.test.projects[0].id
  instance_id = sbercloud_dcs_instance.instance_1.id
  backup_id   = replace(replace(sbercloud_dcs_backup.test1.id, sbercloud_dcs_instance.instance_1.id, ""), "/", "")
  remark      = "restore instance"
  
  depends_on = [
    sbercloud_dcs_backup.test1
  ]
}
	`
}

func CreateInstanceAndBackup() (string, string, string, error) {

	return "", "", "", nil
}
