package rds

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccRdsEngineVersionsV3DataSource_basic(t *testing.T) {
	dataSourceName := "data.sbercloud_rds_engine_versions.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsEngineVersionsV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "type", "MySQL"),
					resource.TestMatchResourceAttr(dataSourceName, "versions.#", regexp.MustCompile("\\d+")),
				),
			},
		},
	})
}

const testAccRdsEngineVersionsV3DataSource_basic string = "data \"sbercloud_rds_engine_versions\" \"test\" {}"
