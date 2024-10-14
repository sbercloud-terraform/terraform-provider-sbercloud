package rds

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPgPluginUpdate_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testPgPluginUpdate_basic(),
			},
		},
	})
}

func testPgPluginUpdate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_pg_plugin" "test" {
  instance_id   = "%[1]s"
  name          = "pgl_ddl_deploy"
  database_name = huaweicloud_rds_pg_database.test.name
}

resource "huaweicloud_rds_pg_plugin_update" "test" {
  depends_on = [huaweicloud_rds_pg_plugin.test]

  instance_id    = "%[1]s"
  database_name  = huaweicloud_rds_pg_database.test.name
  extension_name = huaweicloud_rds_pg_plugin.test.name
}
`, acceptance.SBC_RDS_INSTANCE_ID)
}
