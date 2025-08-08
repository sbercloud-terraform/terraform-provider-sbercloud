package cfw

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

func getLtsLogFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.SBC_REGION_NAME

	var (
		httpUrl = "v1/{project_id}/cfw/logs/configuration"
		product = "cfw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW client: %s", err)
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path += fmt.Sprintf("?fw_instance_id=%s", state.Primary.ID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", path, &opt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CFW lts log configuration: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	lts_enable := utils.PathSearch("data.lts_enable", respBody, float64(0)).(float64)
	if lts_enable == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func TestAccLtsLog_basic(t *testing.T) {
	var obj interface{}

	rName := "sbercloud_cfw_lts_log.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLtsLogFunc,
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
				Config: testAccLtsLog_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "fw_instance_id", acceptance.SBC_CFW_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "lts_attack_log_stream_enable", "1"),
					resource.TestCheckResourceAttr(rName, "lts_access_log_stream_enable", "1"),
					resource.TestCheckResourceAttr(rName, "lts_flow_log_stream_enable", "1"),
					resource.TestCheckResourceAttrSet(rName, "lts_attack_log_stream_id"),
					resource.TestCheckResourceAttrSet(rName, "lts_log_group_id"),
					resource.TestCheckResourceAttrSet(rName, "lts_access_log_stream_id"),
					resource.TestCheckResourceAttrSet(rName, "lts_flow_log_stream_id"),
				),
			},
			{
				Config: testAccLtsLog_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "lts_attack_log_stream_enable", "0"),
					resource.TestCheckResourceAttr(rName, "lts_access_log_stream_enable", "0"),
					resource.TestCheckResourceAttr(rName, "lts_flow_log_stream_enable", "1"),
					resource.TestCheckResourceAttrSet(rName, "lts_log_group_id"),
					resource.TestCheckResourceAttrSet(rName, "lts_flow_log_stream_id"),
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

func testAccLtsLog_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cfw_lts_log" test {
  fw_instance_id               = "%[2]s"
  lts_log_group_id             = sbercloud_lts_group.g1.id
  lts_attack_log_stream_enable = 1
  lts_access_log_stream_enable = 1
  lts_flow_log_stream_enable   = 1
  lts_attack_log_stream_id     = sbercloud_lts_stream.s1.id
  lts_access_log_stream_id     = sbercloud_lts_stream.s2.id
  lts_flow_log_stream_id       = sbercloud_lts_stream.s3.id
}
`, testAccLtsLog_base(), acceptance.SBC_CFW_INSTANCE_ID)
}

func testAccLtsLog_basic_update() string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_cfw_lts_log" test {
  fw_instance_id               = "%[2]s"
  lts_log_group_id             = sbercloud_lts_group.g2.id
  lts_attack_log_stream_enable = 0
  lts_access_log_stream_enable = 0
  lts_flow_log_stream_enable   = 1
  lts_flow_log_stream_id       = sbercloud_lts_stream.s4.id
}
`, testAccLtsLog_base(), acceptance.SBC_CFW_INSTANCE_ID)
}

func testAccLtsLog_base() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "sbercloud_lts_group" "g1" {
  group_name  = "%[1]s"
  ttl_in_days = 1
}

resource "sbercloud_lts_group" "g2" {
  group_name  = "%[1]s_2"
  ttl_in_days = 1
}
	  
resource "sbercloud_lts_stream" "s1" {
  group_id    = sbercloud_lts_group.g1.id
  stream_name = "%[1]s_s1"
}
	  
resource "sbercloud_lts_stream" "s2" {
  group_id    = sbercloud_lts_group.g1.id
  stream_name = "%[1]s_s2"
}
	  
resource "sbercloud_lts_stream" "s3" {
  group_id    = sbercloud_lts_group.g1.id
  stream_name = "%[1]s_s3"
}

resource "sbercloud_lts_stream" "s4" {
  group_id    = sbercloud_lts_group.g2.id
  stream_name = "%[1]s_s4"
}
`, name)
}
