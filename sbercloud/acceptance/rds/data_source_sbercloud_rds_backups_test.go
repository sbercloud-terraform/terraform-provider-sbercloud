package rds

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceBackups_basic(t *testing.T) {
	rName := "data.sbercloud_rds_backups.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceBackups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "backups.0.id", "sbercloud_rds_backup.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "backups.0.name", "sbercloud_rds_backup.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "backups.0.instance_id",
						"sbercloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "backups.0.type", "manual"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.size"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.status"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.begin_time"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.end_time"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.associated_with_ddm"),
					resource.TestCheckResourceAttr(rName, "backups.0.datastore.#", "1"),
					resource.TestCheckResourceAttr(rName, "backups.0.databases.#", "0"),
				),
			},
		},
	})
}

func testAccDatasourceBackups_basic() string {
	backupConfig := testBackup_pg_basic(acceptance.RandomAccResourceName())
	return fmt.Sprintf(`
%s 

data "sbercloud_rds_backups" "test" {
  instance_id = sbercloud_rds_instance.test.id
  backup_type = "manual"

  depends_on = [
    sbercloud_rds_backup.test
  ]
}
`, backupConfig)
}
