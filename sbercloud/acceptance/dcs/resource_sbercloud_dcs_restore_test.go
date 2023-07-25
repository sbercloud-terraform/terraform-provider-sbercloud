package dcs

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDcsRestoreV1_basic(t *testing.T) {
	resourceName := "sbercloud_dcs_restore.test"

	projectId := "0f5181caba0024e72f89c0045e707b91"
	instanceId := "578655e4-5846-4f1b-bfe4-4938ebc7e19e"
	backupId := "ed466175-3a5d-42d0-90b0-bb3ec29e1465"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Restore_basic(projectId, instanceId, backupId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectId),
					resource.TestCheckResourceAttr(resourceName, "instance_id", instanceId),
					resource.TestCheckResourceAttr(resourceName, "backup_id", backupId),
					resource.TestCheckResourceAttr(resourceName, "remark", "restore instance"),
				),
			},
		},
	})
}

func testAccDcsV1Restore_basic(projectId, instanceId, backupId string) string {
	return fmt.Sprintf(`
resource "sbercloud_dcs_restore" "test" {

	project_id  = %q

  	instance_id = %q

  	backup_id   = %q

  	remark      = "restore instance"

}
	`, projectId, instanceId, backupId)
}
