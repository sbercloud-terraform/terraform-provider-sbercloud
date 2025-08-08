package cfw

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cfw"
)

func getEipProtectionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		queryHttpUrl            = "v1/{project_id}/eips/protect"
		getProtectedEipsProduct = "cfw"
	)
	client, err := cfg.NewServiceClient(getProtectedEipsProduct, acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW Client: %s", err)
	}

	resp, err := cfw.QuerySyncedEips(client, queryHttpUrl, state.Primary.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting protected EIPs: %s", err)
	}
	if !cfw.ProtectedEipExist(resp) {
		return nil, golangsdk.ErrDefault404{}
	}
	return resp, nil
}

func TestAccEipProtection_basic(t *testing.T) {
	var (
		obj interface{}

		rName       = "sbercloud_cfw_eip_protection.test"
		basicConfig = testEipProtection_base()

		rc = acceptance.InitResourceCheck(
			rName,
			&obj,
			getEipProtectionResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testEipProtection_basic_step1(basicConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testEipProtection_basic_step2(basicConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testEipProtection_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

resource "sbercloud_vpc_eip" "test" {
  count = 3

  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type  = "PER"
    name        = "%[2]s_${count.index}"
    size        = 10
    charge_mode = "traffic"
  }
}
`, testAccDatasourceFirewalls_basic(), name)
}

func testEipProtection_basic_step1(basicConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cfw_eip_protection" "test" {
  object_id = data.sbercloud_cfw_firewalls.test.records[0].protect_objects[0].object_id

  dynamic "protected_eip" {
    for_each = slice(sbercloud_vpc_eip.test[*], 0, 2)
    content {
      id          = protected_eip.value.id
      public_ipv4 = protected_eip.value.address
    }
  }
}
`, basicConfig)
}

func testEipProtection_basic_step2(basicConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cfw_eip_protection" "test" {
  object_id = data.sbercloud_cfw_firewalls.test.records[0].protect_objects[0].object_id

  dynamic "protected_eip" {
    for_each = slice(sbercloud_vpc_eip.test[*], 1, 3)
    content {
      id          = protected_eip.value.id
      public_ipv4 = protected_eip.value.address
    }
  }
}
`, basicConfig)
}
