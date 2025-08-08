package cfw

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cfw"
)

func getResourceAntiVirusFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.SBC_REGION_NAME
	product := "cfw"

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW client: %s", err)
	}

	return cfw.GetAntiVirusConfigs(client, state.Primary.Attributes["object_id"])
}

func TestAccResourceAntiVirus_basic(t *testing.T) {
	var obj interface{}

	resourceName := "sbercloud_cfw_anti_virus.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceAntiVirusFunc,
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
				Config: testResourceAntiVirus_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "scan_protocol_configs.#", "3"),
				),
			},
			{
				Config: testResourceAntiVirus_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "scan_protocol_configs.#", "2"),
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

func testResourceAntiVirus_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cfw_anti_virus" "test" {
  object_id = data.sbercloud_cfw_firewalls.test.records[0].protect_objects[0].object_id

  scan_protocol_configs {
    protocol_type =  1
    action        =  1
  }

  scan_protocol_configs {
    protocol_type =  2
    action        =  1
  }

  scan_protocol_configs {
    protocol_type =  3
    action        =  1
  }
}
`, testAccDatasourceFirewalls_basic())
}

func testResourceAntiVirus_update() string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cfw_anti_virus" "test" {
  object_id = data.sbercloud_cfw_firewalls.test.records[0].protect_objects[0].object_id

  scan_protocol_configs {
    protocol_type =  3
    action        =  0
  }

  scan_protocol_configs {
    protocol_type =  4
    action        =  1
  }
}
`, testAccDatasourceFirewalls_basic())
}
