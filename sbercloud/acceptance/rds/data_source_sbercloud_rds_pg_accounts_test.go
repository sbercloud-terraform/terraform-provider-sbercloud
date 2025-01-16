package rds

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRdsPgAccounts_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.sbercloud_rds_pg_accounts.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePgAccounts_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "users.#"),
					resource.TestCheckResourceAttrSet(rName, "users.0.name"),
					resource.TestCheckResourceAttrSet(rName, "users.0.attributes.0.rolsuper"),
					resource.TestCheckResourceAttrSet(rName, "users.0.attributes.0.rolinherit"),
					resource.TestCheckResourceAttrSet(rName, "users.0.attributes.0.rolcreaterole"),
					resource.TestCheckResourceAttrSet(rName, "users.0.attributes.0.rolcreatedb"),
					resource.TestCheckResourceAttrSet(rName, "users.0.attributes.0.rolcanlogin"),
					resource.TestCheckResourceAttrSet(rName, "users.0.attributes.0.rolconnlimit"),
					resource.TestCheckResourceAttrSet(rName, "users.0.attributes.0.rolreplication"),
					resource.TestCheckResourceAttrSet(rName, "users.0.attributes.0.rolbypassrls"),
					resource.TestCheckResourceAttrSet(rName, "users.0.description"),
					resource.TestCheckOutput("user_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourcePgAccounts_base(name string) string {
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
    version = "12"
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
}

resource "sbercloud_rds_pg_account" "test" {
  instance_id = sbercloud_rds_instance.test.id
  name        = "%[2]s"
  password    = "TestPass1!23!4"
  description = "test_description"
}
`, testAccRdsInstance_base(name), name)
}

func testAccDatasourcePgAccounts_basic(name string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_rds_pg_accounts" "test" {
  depends_on  = [sbercloud_rds_pg_account.test]
  instance_id = sbercloud_rds_pg_account.test.instance_id
  user_name   = "%s"
}

data "sbercloud_rds_pg_accounts" "user_name_filter" {
  depends_on  = [sbercloud_rds_pg_account.test]
  instance_id = sbercloud_rds_pg_account.test.instance_id
  user_name   = sbercloud_rds_pg_account.test.name
}

locals {
  user_name = sbercloud_rds_pg_account.test.name
}

output "user_name_filter_is_useful" {
  value = length(data.sbercloud_rds_pg_accounts.user_name_filter.users) > 0 && alltrue(
    [for v in data.sbercloud_rds_pg_accounts.user_name_filter.users[*].name : v == local.user_name]
  )
}

`, testAccDatasourcePgAccounts_base(name), name)
}
