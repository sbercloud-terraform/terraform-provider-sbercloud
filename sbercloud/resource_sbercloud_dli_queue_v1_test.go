package sbercloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccDliQueueV1_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDliQueueV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDliQueueV1_basic(acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDliQueueV1Exists(),
				),
			},
		},
	})
}

func testAccDliQueueV1_basic(val string) string {
	return fmt.Sprintf(`
resource "sbercloud_dli_queue" "queue" {
  name = "terraform_dli_queue_test_%s"
  cu_count = 16
}
	`, val)
}

func testAccCheckDliQueueV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.DliV1Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_dli_queue" {
			continue
		}

		_, err = fetchDliQueueV1ByListOnTest(rs, client)
		if err != nil {
			if strings.Index(err.Error(), "Error finding the resource by list api") != -1 {
				return nil
			}
			return err
		}
		return fmt.Errorf("sbercloud_dli_queue still exists")
	}

	return nil
}

func testAccCheckDliQueueV1Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.DliV1Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["sbercloud_dli_queue.queue"]
		if !ok {
			return fmt.Errorf("Error checking sbercloud_dli_queue.queue exist, err=not found this resource")
		}

		_, err = fetchDliQueueV1ByListOnTest(rs, client)
		if err != nil {
			if strings.Index(err.Error(), "Error finding the resource by list api") != -1 {
				return fmt.Errorf("sbercloud_dli_queue is not exist")
			}
			return fmt.Errorf("Error checking sbercloud_dli_queue.queue exist, err=%s", err)
		}
		return nil
	}
}

func fetchDliQueueV1ByListOnTest(rs *terraform.ResourceState,
	client *golangsdk.ServiceClient) (interface{}, error) {
	link := client.ServiceURL("queues")

	return findDliQueueV1ByList(client, link, rs.Primary.ID)
}

func findDliQueueV1ByList(client *golangsdk.ServiceClient, link, resourceID string) (interface{}, error) {
	r, err := sendDliQueueV1ListRequest(client, link)
	if err != nil {
		return nil, err
	}
	for _, item := range r.([]interface{}) {
		val, ok := item.(map[string]interface{})["queue_name"]
		if ok && resourceID == convertToStr(val) {
			return item, nil
		}
	}

	return nil, fmtp.Errorf("Error finding the resource by list api")
}

func sendDliQueueV1ListRequest(client *golangsdk.ServiceClient, url string) (interface{}, error) {
	r := golangsdk.Result{}
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/json"}})
	if r.Err != nil {
		return nil, fmtp.Errorf("Error running api(list) for resource(DliQueueV1), err=%s", r.Err)
	}

	v, err := navigateValue(r.Body, []string{"queues"}, nil)
	if err != nil {
		return nil, err
	}
	return v, nil
}
