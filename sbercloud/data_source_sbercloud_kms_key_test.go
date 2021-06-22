package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccKmsKeyDataSource_Basic(t *testing.T) {
	var keyAlias = fmt.Sprintf("key_alias_%s", acctest.RandString(5))
	var datasourceName = "data.sbercloud_kms_key.key_1"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyDataSource_Basic(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyDataSourceID(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(datasourceName, "region", SBC_REGION_NAME),
				),
			},
		},
	})
}

func TestAccKmsKeyDataSource_WithTags(t *testing.T) {
	var keyAlias = fmt.Sprintf("key_alias_%s", acctest.RandString(5))
	var datasourceName = "data.sbercloud_kms_key.key_1"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyDataSource_WithTags(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyDataSourceID(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(datasourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(datasourceName, "tags.key", "value"),
				),
			},
		},
	})
}

func TestAccKmsKeyDataSource_WithEpsId(t *testing.T) {
	var keyAlias = fmt.Sprintf("key_alias_%s", acctest.RandString(5))
	var datasourceName = "data.sbercloud_kms_key.key_1"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckEpsID(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyDataSource_epsId(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyDataSourceID(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(datasourceName, "enterprise_project_id", SBC_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckKmsKeyDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Kms key data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Kms key data source ID not set")
		}

		return nil
	}
}

func testAccKmsKeyDataSource_Basic(keyAlias string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_kms_key" "key_1" {
  key_alias = sbercloud_kms_key.key_1.key_alias
  key_id    = sbercloud_kms_key.key_1.id
  key_state = "2"
}
`, testAccKmsKey_Basic(keyAlias))
}

func testAccKmsKeyDataSource_WithTags(keyAlias string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_kms_key" "key_1" {
  key_alias = sbercloud_kms_key.key_1.key_alias
  key_id    = sbercloud_kms_key.key_1.id
  key_state = "2"
}
`, testAccKmsKey_WithTags(keyAlias))
}

func testAccKmsKeyDataSource_epsId(keyAlias string) string {
	return fmt.Sprintf(`
resource "sbercloud_kms_key" "key_1" {
  key_alias       = "%s"
  key_description = "test description"
  pending_days    = "7"
  is_enabled      = true
  enterprise_project_id = "%s"
}

data "sbercloud_kms_key" "key_1" {
  key_alias       = sbercloud_kms_key.key_1.key_alias
  key_id          = sbercloud_kms_key.key_1.id
  key_description = "test description"
  key_state       = "2"
  enterprise_project_id = sbercloud_kms_key.key_1.enterprise_project_id
}
`, keyAlias, SBC_ENTERPRISE_PROJECT_ID_TEST)
}
