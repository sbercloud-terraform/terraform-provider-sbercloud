package dis

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDisStreamV2_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDisStreamV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDisStreamV2_basic(acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDisStreamV2Exists(),
				),
			},
		},
	})
}

func testAccCheckDisStreamV2Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.DisV2Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_dis_stream" {
			continue
		}

		url, err := acceptance.ReplaceVarsForTest(rs, "streams/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err == nil {
			return fmt.Errorf("sbercloud_dis_stream still exists at %s", url)
		}
	}

	return nil
}

func testAccCheckDisStreamV2Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.DisV2Client(acceptance.SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["sbercloud_dis_stream.stream"]
		if !ok {
			return fmt.Errorf("Error checking sbercloud_dis_stream.stream exist, err=not found this resource")
		}

		url, err := acceptance.ReplaceVarsForTest(rs, "streams/{id}")
		if err != nil {
			return fmt.Errorf("Error checking sbercloud_dis_stream.stream exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmt.Errorf("sbercloud_dis_stream.stream is not exist")
			}
			return fmt.Errorf("Error checking sbercloud_dis_stream.stream exist, err=send request failed: %s", err)
		}
		return nil
	}
}

func testAccDisStreamV2_basic(val string) string {
	return fmt.Sprintf(`
resource "sbercloud_dis_stream" "stream" {
 stream_name = "terraform_test_dis_stream_%s"
 partition_count = 1
}
	`, val)
}
