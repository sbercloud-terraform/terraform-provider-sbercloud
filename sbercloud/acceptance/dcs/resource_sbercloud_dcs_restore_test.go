package dcs

import (
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDcsRestoreV1_basic(t *testing.T) {
	resourceName := "sbercloud_dcs_restore.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Restore_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "remark", "restore instance"),
				),
			},
		},
	})
}

func testAccDcsV1Restore_basic() string {
	return `
data "sbercloud_dcs_flavors" "single_flavors" {
  cache_mode = "ha"
  capacity   = 1
}

data "sbercloud_vpc" "vpc" {
  name = "vpc-default"
}

data "sbercloud_vpc_subnet" "subnet" {
  id = "c81b93ad-65d7-449c-83ab-600939bfce5a"
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
  vpc_id             = data.sbercloud_vpc.vpc.id
  subnet_id          = data.sbercloud_vpc_subnet.subnet.id
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
