package cbh

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getAssetAgencyAuthorizationResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.SBC_REGION_NAME
		product = "cbh"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CBH client: %s", err)
	}

	basePath := client.Endpoint + "v2/{project_id}/cbs/agency/authorization"
	basePath = strings.ReplaceAll(basePath, "{project_id}", client.ProjectID)
	baseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", basePath, &baseOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CBH asset agency authorization: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccAssetAgencyAuthorization_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "sbercloud_cbh_asset_agency_authorization.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAssetAgencyAuthorizationResourceFunc,
	)

	// Avoid CheckDestroy, because there is nothing in the resource destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAssetAgencyAuthorization_basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "csms", "true"),
					resource.TestCheckResourceAttr(rName, "kms", "true"),
				),
			},
			{
				Config: testAccAssetAgencyAuthorization_update1,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "csms", "true"),
					resource.TestCheckResourceAttr(rName, "kms", "false")),
			},
			{
				Config: testAccAssetAgencyAuthorization_update2,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "csms", "false"),
					resource.TestCheckResourceAttr(rName, "kms", "false")),
			},
		},
	})
}

const testAccAssetAgencyAuthorization_basic = `
resource "sbercloud_cbh_asset_agency_authorization" "test" {
  csms = true
  kms  = true
}
`

const testAccAssetAgencyAuthorization_update1 = `
resource "sbercloud_cbh_asset_agency_authorization" "test" {
  csms = true
  kms  = false
}
`

const testAccAssetAgencyAuthorization_update2 = `
resource "sbercloud_cbh_asset_agency_authorization" "test" {
  csms = false
  kms  = false
}
`
