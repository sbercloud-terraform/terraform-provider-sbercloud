package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/kms/v1/keys"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccKmsKey_Basic(t *testing.T) {
	var key keys.Key
	var keyAlias = fmt.Sprintf("kms_%s", acctest.RandString(5))
	var keyAliasUpdate = fmt.Sprintf("kms_updated_%s", acctest.RandString(5))
	var resourceName = "sbercloud_kms_key.key_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKmsKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_Basic(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyExists(resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "region", SBC_REGION_NAME),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"key_usage",
					"pending_days",
				},
			},
			{
				Config: testAccKmsKeyUpdate(keyAliasUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyExists(resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAliasUpdate),
					resource.TestCheckResourceAttr(resourceName, "key_description", "key update description"),
					resource.TestCheckResourceAttr(resourceName, "region", SBC_REGION_NAME),
				),
			},
		},
	})
}

func TestAccKmsKey_Enable(t *testing.T) {
	var key1, key2, key3 keys.Key
	rName := fmt.Sprintf("kms_%s", acctest.RandString(5))
	var resourceName = "sbercloud_kms_key.key_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKmsKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_enabled(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyExists(resourceName, &key1),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
					testAccCheckKmsKeyIsEnabled(&key1, true),
				),
			},
			{
				Config: testAccKmsKey_disabled(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyExists(resourceName, &key2),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
					testAccCheckKmsKeyIsEnabled(&key2, false),
				),
			},
			{
				Config: testAccKmsKey_enabled(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyExists(resourceName, &key3),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "true"),
					testAccCheckKmsKeyIsEnabled(&key3, true),
				),
			},
		},
	})
}

func TestAccKmsKey_WithTags(t *testing.T) {
	var key keys.Key
	var keyAlias = fmt.Sprintf("kms_%s", acctest.RandString(5))
	var resourceName = "sbercloud_kms_key.key_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKmsKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_WithTags(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyExists(resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
		},
	})
}

func TestAccKmsKey_WithEpsId(t *testing.T) {
	var key keys.Key
	var keyAlias = fmt.Sprintf("kms_%s", acctest.RandString(5))
	var resourceName = "sbercloud_kms_key.key_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEpsID(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKmsKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey_epsId(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyExists(resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", SBC_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckKmsKeyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	kmsClient, err := config.KmsKeyV1Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud kms client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_kms_key" {
			continue
		}
		v, err := keys.Get(kmsClient, rs.Primary.ID).ExtractKeyInfo()
		if err != nil {
			return err
		}
		if v.KeyState != "4" {
			return fmt.Errorf("key still exists")
		}
	}
	return nil
}

func testAccCheckKmsKeyExists(n string, key *keys.Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		kmsClient, err := config.KmsKeyV1Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud kms client: %s", err)
		}
		found, err := keys.Get(kmsClient, rs.Primary.ID).ExtractKeyInfo()
		if err != nil {
			return err
		}
		if found.KeyID != rs.Primary.ID {
			return fmt.Errorf("key not found")
		}

		*key = *found
		return nil
	}
}

func testAccCheckKmsKeyIsEnabled(key *keys.Key, isEnabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if (key.KeyState == huaweicloud.EnabledState) != isEnabled {
			return fmt.Errorf("Expected key %s to have is_enabled=%t, given %s",
				key.KeyID, isEnabled, key.KeyState)
		}

		return nil
	}
}

func testAccKmsKey_Basic(keyAlias string) string {
	return fmt.Sprintf(`
resource "sbercloud_kms_key" "key_1" {
  key_alias    = "%s"
  pending_days = "7"
  region       = "%s"
}
`, keyAlias, SBC_REGION_NAME)
}

func testAccKmsKey_WithTags(keyAlias string) string {
	return fmt.Sprintf(`
resource "sbercloud_kms_key" "key_1" {
  key_alias    = "%s"
  pending_days = "7"
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, keyAlias)
}

func testAccKmsKey_epsId(keyAlias string) string {
	return fmt.Sprintf(`
resource "sbercloud_kms_key" "key_1" {
  key_alias    = "%s"
  pending_days = "7"
  enterprise_project_id = "%s"
}
`, keyAlias, SBC_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccKmsKeyUpdate(keyAliasUpdate string) string {
	return fmt.Sprintf(`
resource "sbercloud_kms_key" "key_1" {
  key_alias       = "%s"
  key_description = "key update description"
  pending_days    = "7"
}
`, keyAliasUpdate)
}

func testAccKmsKey_enabled(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_kms_key" "key_1" {
  key_description = "Terraform acc test is_enabled %s"
  pending_days    = "7"
  key_alias       = "tf-acc-test-kms-key-%s"
}`, rName, rName)
}

func testAccKmsKey_disabled(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_kms_key" "key_1" {
  key_description = "Terraform acc test is_enabled %s"
  pending_days    = "7"
  key_alias       = "tf-acc-test-kms-key-%s"
  is_enabled      = false
}`, rName, rName)
}
