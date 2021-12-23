package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func getBandwidthResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.NetworkingV1Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SberCloud Network client: %s", err)
	}
	return bandwidths.Get(c, state.Primary.ID).Extract()
}

func TestAccVpcBandWidth_basic(t *testing.T) {
	var bandwidth bandwidths.BandWidth

	randName := acceptance.RandomAccResourceName()
	resourceName := "sbercloud_vpc_bandwidth.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&bandwidth,
		getBandwidthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcBandWidth_basic(randName, 5),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "size", "5"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "WHOLE"),
					resource.TestCheckResourceAttr(resourceName, "status", "NORMAL"),
				),
			},
			{
				Config: testAccVpcBandWidth_basic(randName+"_update", 6),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"_update"),
					resource.TestCheckResourceAttr(resourceName, "size", "6"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "WHOLE"),
					resource.TestCheckResourceAttr(resourceName, "status", "NORMAL"),
				),
			},
		},
	})
}

func TestAccVpcBandWidth_WithEpsId(t *testing.T) {
	var bandwidth bandwidths.BandWidth

	randName := acceptance.RandomAccResourceName()
	resourceName := "sbercloud_vpc_bandwidth.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&bandwidth,
		getBandwidthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcBandWidth_epsId(randName, 5),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.SBC_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

func testAccVpcBandWidth_basic(rName string, size int) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc_bandwidth" "test" {
  name = "%s"
  size = "%d"
}
`, rName, size)
}

func testAccVpcBandWidth_epsId(rName string, size int) string {
	return fmt.Sprintf(`
resource "sbercloud_vpc_bandwidth" "test" {
  name                  = "%s"
  size                  = "%d"
  enterprise_project_id = "%s"
}
`, rName, size, acceptance.SBC_ENTERPRISE_PROJECT_ID)
}
