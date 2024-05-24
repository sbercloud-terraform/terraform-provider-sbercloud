package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dli/v1/databases"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func getDatabaseResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.DliV1Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DLI v1 client: %s", err)
	}

	return dli.GetDliSQLDatabaseByName(c, state.Primary.Attributes["name"])
}

func TestAccDliDatabase_basic(t *testing.T) {
	var database databases.Database

	rName := acceptance.RandomAccResourceName()
	resourceName := "sbercloud_dli_database.test"

	rc := acceptance.InitResourceCheck(
		resourceName,

		&database,
		getDatabaseResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliDatabase_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "For terraform acc test"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.SBC_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(resourceName, "owner"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDatabaseImportStateFunc(resourceName),
			},
		},
	})
}

func testAccDatabaseImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		name := rs.Primary.Attributes["name"]
		if name == "" {
			return "", fmt.Errorf("the database name is incorrect, got '%s'", name)
		}
		return name, nil
	}
}

func testAccDliDatabase_basic(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_dli_database" "test" {
  name                  = "%s"
  description           = "For terraform acc test"
  enterprise_project_id = "%s"
}
`, rName, acceptance.SBC_ENTERPRISE_PROJECT_ID_TEST)
}
