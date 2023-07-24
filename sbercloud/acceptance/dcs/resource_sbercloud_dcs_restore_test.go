package dcs

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dcs/v1/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDcsRestoreV1_basic(t *testing.T) {
	var instance instances.Instance
	resourceName := "sbercloud_dcs_restore.test"

	projectId := "0f5181caba0024e72f89c0045e707b91"
	instanceId := "578655e4-5846-4f1b-bfe4-4938ebc7e19e"
	backupId := "ed466175-3a5d-42d0-90b0-bb3ec29e1465"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDcsRestoreV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Restore_basic(projectId, instanceId, backupId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsRestoreV1Exists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectId),
					resource.TestCheckResourceAttr(resourceName, "instance_id", instanceId),
					resource.TestCheckResourceAttr(resourceName, "backup_id", backupId),
					resource.TestCheckResourceAttr(resourceName, "remark", "restore instance"),
				),
			},
		},
	})
}

func testAccCheckDcsRestoreV1Destroy(s *terraform.State) error {
	//config := acceptance.TestAccProvider.Meta().(*config.Config)
	//dcsClient, err := config.DcsV1Client(acceptance.SBC_REGION_NAME)
	//if err != nil {
	//	return fmt.Errorf("Error creating SberCloud instance client: %s", err)
	//}
	//
	//for _, rs := range s.RootModule().Resources {
	//	if rs.Type != "sbercloud_dcs_instance" {
	//		continue
	//	}
	//
	//	_, err := instances.Get(dcsClient, rs.Primary.ID).Extract()
	//	if err == nil {
	//		return fmt.Errorf("the DCS instance still exists")
	//	}
	//}
	return nil
}

func testAccCheckDcsRestoreV1Exists(n string, instance instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		//rs, ok := s.RootModule().Resources[n]
		//if !ok {
		//	return fmt.Errorf("Not found: %s", n)
		//}
		//
		//if rs.Primary.ID == "" {
		//	return fmt.Errorf("No ID is set")
		//}
		//
		//config := acceptance.TestAccProvider.Meta().(*config.Config)
		//dcsClient, err := config.DcsV1Client(acceptance.SBC_REGION_NAME)
		//if err != nil {
		//	return fmt.Errorf("Error creating SberCloud instance client: %s", err)
		//}
		//
		//v, err := instances.Get(dcsClient, rs.Primary.ID).Extract()
		//if err != nil {
		//	return fmt.Errorf("Error getting SberCloud instance: %s, err: %s", rs.Primary.ID, err)
		//}
		//
		//if v.InstanceID != rs.Primary.ID {
		//	return fmt.Errorf("the DCS instance not found")
		//}
		//instance = *v
		return nil
	}
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
