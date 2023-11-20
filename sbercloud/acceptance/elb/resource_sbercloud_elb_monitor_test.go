package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/monitors"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func getELBMonitorResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ElbV3Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}
	return monitors.Get(client, state.Primary.ID).Extract()
}

func TestAccElbV3Monitor_basic(t *testing.T) {
	var monitor monitors.Monitor
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_elb_monitor.monitor_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&monitor,
		getELBMonitorResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3MonitorConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "interval", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "5"),
					resource.TestCheckResourceAttr(resourceName, "url_path", "/aa"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "www.aa.com"),
					resource.TestCheckResourceAttr(resourceName, "port", "8000"),
					resource.TestCheckResourceAttr(resourceName, "status_code", "200,401-500,502"),
				),
			},
			{
				Config: testAccElbV3MonitorConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "interval", "30"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "20"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "8"),
					resource.TestCheckResourceAttr(resourceName, "url_path", "/bb"),
					resource.TestCheckResourceAttr(resourceName, "domain_name", "www.bb.com"),
					resource.TestCheckResourceAttr(resourceName, "port", "8888"),
					resource.TestCheckResourceAttr(resourceName, "status_code", "200,301,404-500,504"),
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

func testAccElbV3MonitorConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_elb_monitor" "monitor_1" {
  pool_id     = sbercloud_elb_pool.test.id
  protocol    = "HTTP"
  interval    = 20
  timeout     = 10
  max_retries = 5
  url_path    = "/aa"
  domain_name = "www.aa.com"
  port        = "8000"
  status_code = "200,401-500,502"
}
`, testAccElbV3PoolConfig_basic(rName))
}

func testAccElbV3MonitorConfig_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_elb_monitor" "monitor_1" {
  pool_id     = sbercloud_elb_pool.test.id
  protocol    = "HTTPS"
  interval    = 30
  timeout     = 20
  max_retries = 8
  url_path    = "/bb"
  domain_name = "www.bb.com"
  port        = 8888
  status_code = "200,301,404-500,504"
}
`, testAccElbV3PoolConfig_basic(rName))
}
