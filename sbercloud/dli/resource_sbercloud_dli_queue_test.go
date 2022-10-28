package dli

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dli/v1/queues"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	act "github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDliQueueResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.DliV1Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating Dli v1 client, err=%s", err)
	}

	result := queues.Get(client, state.Primary.ID)
	return result.Body, result.Err
}

func TestAccDliQueue_basic(t *testing.T) {
	rName := act.RandomAccResourceName()
	resourceName := "sbercloud_dli_queue.test"

	var obj queues.CreateOpts
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliQueueResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { act.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliQueue_basic(rName, dli.CU_16),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "queue_type", dli.QUEUE_TYPE_SQL),
					resource.TestCheckResourceAttr(resourceName, "cu_count", fmt.Sprintf("%d", dli.CU_16)),
					resource.TestCheckResourceAttrSet(resourceName, "resource_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"tags",
				},
			},
		},
	})
}

func testAccDliQueue_basic(rName string, cuCount int) string {
	return fmt.Sprintf(`
resource "sbercloud_dli_queue" "test" {
  name     = "%s"
  cu_count = %d

  tags = {
    foo = "bar"
  }
}
`, rName, cuCount)
}
