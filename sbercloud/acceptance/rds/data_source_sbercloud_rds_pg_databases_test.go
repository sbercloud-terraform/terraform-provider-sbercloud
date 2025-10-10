package rds

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRdsPgDatabases_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.sbercloud_rds_pg_databases.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePgDatabases_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "databases.#"),
					resource.TestCheckResourceAttrSet(rName, "databases.0.name"),
					resource.TestCheckResourceAttrSet(rName, "databases.0.owner"),
					resource.TestCheckResourceAttrSet(rName, "databases.0.character_set"),
					resource.TestCheckResourceAttrSet(rName, "databases.0.lc_collate"),
					resource.TestCheckResourceAttrSet(rName, "databases.0.size"),
					resource.TestCheckResourceAttrSet(rName, "databases.0.description"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("owner_filter_is_useful", "true"),
					resource.TestCheckOutput("character_set_filter_is_useful", "true"),
					resource.TestCheckOutput("lc_collate_filter_is_useful", "true"),
					resource.TestCheckOutput("size_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourcePgDatabases_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.sbercloud_availability_zones.test.names[0]]
  security_group_id = sbercloud_networking_secgroup.test.id
  subnet_id         = data.sbercloud_vpc_subnet.test.id
  vpc_id            = data.sbercloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    type    = "PostgreSQL"
    version = "14"
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}

resource "sbercloud_rds_pg_database" "test" {
  instance_id   = sbercloud_rds_instance.test.id
  name          = "%[2]s"
  owner         = "root"
  character_set = "UTF8"
  template      = "template1"
  lc_collate    = "en_US.UTF-8"
  lc_ctype      = "en_US.UTF-8"
}
`, testAccRdsInstance_base(name), name)
}

func testAccDatasourcePgDatabases_basic(name string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_rds_pg_databases" "test" {
  depends_on  = [sbercloud_rds_pg_database.test]
  instance_id = sbercloud_rds_pg_database.test.instance_id
}

data "sbercloud_rds_pg_databases" "name_filter" {
  depends_on  = [sbercloud_rds_pg_database.test]
  instance_id = sbercloud_rds_pg_database.test.instance_id
  name        = sbercloud_rds_pg_database.test.name
}

locals {
  name = sbercloud_rds_pg_database.test.name
}

output "name_filter_is_useful" {
  value = length(data.sbercloud_rds_pg_databases.name_filter.databases) > 0 && alltrue(
    [for v in data.sbercloud_rds_pg_databases.name_filter.databases[*].name : v == local.name]
  )
}

data "sbercloud_rds_pg_databases" "owner_filter" {
  depends_on  = [sbercloud_rds_pg_database.test]
  instance_id = sbercloud_rds_pg_database.test.instance_id
  owner       = sbercloud_rds_pg_database.test.owner
}

locals {
  owner = sbercloud_rds_pg_database.test.owner
}

output "owner_filter_is_useful" {
  value = length(data.sbercloud_rds_pg_databases.owner_filter.databases) > 0 && alltrue(
    [for v in data.sbercloud_rds_pg_databases.owner_filter.databases[*].owner : v == local.owner]
  )
}

data "sbercloud_rds_pg_databases" "character_set_filter" {
  depends_on    = [sbercloud_rds_pg_database.test]
  instance_id   = sbercloud_rds_pg_database.test.instance_id
  character_set = sbercloud_rds_pg_database.test.character_set
}

locals {
  character_set = sbercloud_rds_pg_database.test.character_set
}

output "character_set_filter_is_useful" {
  value = length(data.sbercloud_rds_pg_databases.character_set_filter.databases) > 0 && alltrue(
    [for v in data.sbercloud_rds_pg_databases.character_set_filter.databases[*].character_set : v == local.character_set]
  )
}

data "sbercloud_rds_pg_databases" "lc_collate_filter" {
  depends_on  = [sbercloud_rds_pg_database.test]
  instance_id = sbercloud_rds_pg_database.test.instance_id
  lc_collate  = sbercloud_rds_pg_database.test.lc_collate
}

locals {
  lc_collate = sbercloud_rds_pg_database.test.lc_collate
}

output "lc_collate_filter_is_useful" {
  value = length(data.sbercloud_rds_pg_databases.lc_collate_filter.databases) > 0 && alltrue(
    [for v in data.sbercloud_rds_pg_databases.lc_collate_filter.databases[*].lc_collate : v == local.lc_collate]
  )
}

data "sbercloud_rds_pg_databases" "size_filter" {
  depends_on  = [sbercloud_rds_pg_database.test]
  instance_id = sbercloud_rds_pg_database.test.instance_id
  size        = sbercloud_rds_pg_database.test.size
}

locals {
  size = sbercloud_rds_pg_database.test.size
}

output "size_filter_is_useful" {
  value = length(data.sbercloud_rds_pg_databases.size_filter.databases) > 0 && alltrue(
    [for v in data.sbercloud_rds_pg_databases.size_filter.databases[*].size : v == local.size]
  )
}

`, testAccDatasourcePgDatabases_base(name))
}
