package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cci/v1/namespaces"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cci"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func getV2NamespaceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("cci", acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SberCloud CCI v2 client: %s", err)
	}
	return cci.GetNamespaceDetail(client, state.Primary.ID)
}

func TestAccV2Namespace_basic(t *testing.T) {
	var ns namespaces.Namespace
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_cci_namespace.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ns,
		getV2NamespaceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2Namespace_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "api_version", "cci/v2"),
					resource.TestCheckResourceAttr(resourceName, "kind", "Namespace"),
					resource.TestCheckResourceAttr(resourceName, "finalizers.0", "kubernetes"),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
					resource.TestCheckResourceAttrSet(resourceName, "annotations.%"),
					resource.TestCheckResourceAttrSet(resourceName, "labels.%"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_timestamp"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_version"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccV2Namespace_basic(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_cci_namespace" "test" {
  name = "%s"
}
`, rName)
}
