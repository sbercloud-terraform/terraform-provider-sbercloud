package rds

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRdsPgRoles_basic(t *testing.T) {
	dataSource := "data.sbercloud_rds_pg_roles.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRdsPgRoles_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "roles.#"),
					resource.TestCheckResourceAttr(dataSource, "roles.#", "1"),
					resource.TestCheckResourceAttrPair(dataSource, "roles.0",
						"sbercloud_rds_pg_account.test", "name"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRdsPgRoles_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_rds_instance" "test" {
  name               = "%[2]s"
  flavor             = "rds.pg.n1.large.2"
  availability_zone  = [data.sbercloud_availability_zones.test.names[0]]
  security_group_id  = sbercloud_networking_secgroup.test.id
  subnet_id          = data.sbercloud_vpc_subnet.test.id
  vpc_id             = data.sbercloud_vpc.test.id
  time_zone          = "UTC+08:00"

  db {
    type    = "PostgreSQL"
    version = "12"
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, testAccRdsInstance_base(name), name)
}

func testDataSourceDataSourceRdsPgRoles_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_rds_pg_account" "test" {
  instance_id = sbercloud_rds_instance.test.id
  name        = "%[2]s"
  password    = "TestPass1!23!4"
}

data "sbercloud_rds_pg_roles" "test" {
  depends_on = [sbercloud_rds_pg_account.test]

  instance_id = sbercloud_rds_instance.test.id
  account     = "root"
}
`, testDataSourceDataSourceRdsPgRoles_base(name), name)
}
