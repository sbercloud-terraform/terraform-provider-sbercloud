package lb

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/elb/v2/loadbalancers"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccELBV2LoadbalancerDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName1 := "data.sbercloud_lb_loadbalancer.test_by_name"
	dc1 := acceptance.InitDataSourceCheck(dataSourceName1)
	dataSourceName2 := "data.sbercloud_lb_loadbalancer.test_by_description"
	dc2 := acceptance.InitDataSourceCheck(dataSourceName2)

	var lb loadbalancers.LoadBalancer
	resourceName := "sbercloud_lb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getLoadBalancerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccELBV2LoadbalancerDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName1, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName2, "name", rName),
				),
			},
		},
	})
}

func testAccELBV2LoadbalancerDataSource_basic(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "sbercloud_lb_loadbalancer" "test" {
  name          = "%s"
  vip_subnet_id = data.sbercloud_vpc_subnet.test.ipv4_subnet_id
  description   = "test for load balancer data source"
}

data "sbercloud_lb_loadbalancer" "test_by_name" {
  name = sbercloud_lb_loadbalancer.test.name

  depends_on = [sbercloud_lb_loadbalancer.test]
}

data "sbercloud_lb_loadbalancer" "test_by_description" {
  description = sbercloud_lb_loadbalancer.test.description

  depends_on = [sbercloud_lb_loadbalancer.test]
}
`, rName)
}
