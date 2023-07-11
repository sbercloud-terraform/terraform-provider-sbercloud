package iam

import (
	"fmt"
	"testing"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIdentityCustomRoleDataSource_basic(t *testing.T) {
	resourceName := "data.sbercloud_identity_custom_role.role_1"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCustomRoleDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(
						"data.sbercloud_identity_custom_role.role_1", "name", rName),
				),
			},
		},
	})
}

func testAccIdentityCustomRoleDataSource_basic(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_identity_role" test {
  name        = "%s"
  description = "created by terraform"
  type        = "AX"
  policy      = <<EOF
{
  "Version": "1.1",
  "Statement": [
    {
      "Action": [
        "obs:bucket:GetBucketAcl"
      ],
      "Effect": "Allow",
      "Resource": [
        "obs:*:*:bucket:*"
      ]
    }
  ]
}
EOF
}

data "sbercloud_identity_custom_role" "role_1" {
  name = sbercloud_identity_role.test.name

  depends_on = [
	sbercloud_identity_role.test
  ]
}
`, rName)
}
