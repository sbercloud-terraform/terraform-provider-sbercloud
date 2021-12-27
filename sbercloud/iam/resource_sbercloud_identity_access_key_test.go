package iam

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/credentials"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getIdentityAccessKeyResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	iamClient, err := c.IAMV3Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("Error creating SberCloud identity client: %s", err)
	}

	found, err := credentials.Get(iamClient, state.Primary.ID).Extract()
	if err != nil {
		return nil, err
	}

	if found.AccessKey != state.Primary.ID {
		return nil, fmtp.Errorf("Access Key not found")
	}
	return found, nil
}

func TestAccIdentityAccessKey_basic(t *testing.T) {
	var cred credentials.Credential
	var userName = acceptance.RandomAccResourceName()
	resourceName := "sbercloud_identity_access_key.key_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&cred,
		getIdentityAccessKeyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityAccessKey_basic(userName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "description", "access key by terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				Config: testAccIdentityAccessKey_update(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "inactive"),
					resource.TestCheckResourceAttr(resourceName, "description", "access key by terraform updated"),
				),
			},
		},
	})
}

func testAccIdentityAccessKey_basic(userName string) string {
	return fmt.Sprintf(`
resource "sbercloud_identity_user" "user_1" {
  name        = "%s"
  password    = "password123@!"
  enabled     = true
  description = "tested by terraform"
}

resource "sbercloud_identity_access_key" "key_1" {
  user_id     = sbercloud_identity_user.user_1.id
  description = "access key by terraform"
  secret_file = "./credentials.csv"
}
`, userName)
}

func testAccIdentityAccessKey_update(userName string) string {
	return fmt.Sprintf(`
resource "sbercloud_identity_user" "user_1" {
  name        = "%s"
  password    = "password123@!"
  enabled     = true
  description = "tested by terraform"
}

resource "sbercloud_identity_access_key" "key_1" {
  user_id     = sbercloud_identity_user.user_1.id
  description = "access key by terraform updated"
  secret_file = "./credentials.csv"
  status      = "inactive"
}
`, userName)
}
