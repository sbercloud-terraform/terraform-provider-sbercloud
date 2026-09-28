package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cci"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func getPersistentVolumeClaimResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("cci", acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCI client: %s", err)
	}
	return cci.GetV2PersistentVolumeClaimDetail(client, state.Primary.Attributes["namespace"], state.Primary.Attributes["name"])
}

func TestAccV2PersistentVolumeClaim_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_cci_pvc.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPersistentVolumeClaimResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2PersistentVolumeClaim_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "namespace", "sbercloud_cci_namespace.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "storage_class_name", "csi-obs"),
					resource.TestCheckResourceAttr(resourceName, "access_modes.0", "ReadWriteMany"),
					resource.TestCheckResourceAttr(resourceName, "volume_mode", "Filesystem"),
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

func testAccV2PersistentVolumeClaim_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cci_pvc" "test" {
  name               = "%[2]s"
  namespace          = sbercloud_cci_namespace.test.name
  access_modes       = ["ReadWriteMany"]
  storage_class_name = "csi-obs"
  volume_mode        = "Filesystem"

  annotations = {
    "everest.io/obs-volume-type"       = "STANDARD"
    "csi.storage.k8s.io/fstype"        = "s3fs"
    "everest.io/enterprise-project-id" = "0"
  }

  resources {
    requests = {
      storage = "1Gi"
    }
  }
}
`, testAccV2Namespace_basic(rName), rName)
}
