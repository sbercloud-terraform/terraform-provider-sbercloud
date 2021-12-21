package sbercloud

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/smn/v2/subscriptions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccSMNV2Subscription_basic(t *testing.T) {
	var subscription1 subscriptions.SubscriptionGet
	var subscription2 subscriptions.SubscriptionGet
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSMNSubscriptionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSMNV2SubscriptionConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSMNV2SubscriptionExists("sbercloud_smn_subscription.subscription_1", &subscription1),
					testAccCheckSMNV2SubscriptionExists("sbercloud_smn_subscription.subscription_2", &subscription2),
					resource.TestCheckResourceAttr(
						"sbercloud_smn_subscription.subscription_1", "endpoint",
						"mailtest@gmail.com"),
					resource.TestCheckResourceAttr(
						"sbercloud_smn_subscription.subscription_2", "endpoint",
						"13600000000"),
				),
			},
		},
	})
}

func testAccCheckSMNSubscriptionV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	smnClient, err := config.SmnV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud smn: %s", err)
	}
	var subscription *subscriptions.SubscriptionGet
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_smn_subscription" {
			continue
		}
		foundList, err := subscriptions.List(smnClient).Extract()
		if err != nil {
			return err
		}
		for _, subObject := range foundList {
			if subObject.SubscriptionUrn == rs.Primary.ID {
				subscription = &subObject
			}
		}
		if subscription != nil {
			return fmt.Errorf("subscription still exists")
		}
	}

	return nil
}

func testAccCheckSMNV2SubscriptionExists(n string, subscription *subscriptions.SubscriptionGet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		smnClient, err := config.SmnV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud smn client: %s", err)
		}

		foundList, err := subscriptions.List(smnClient).Extract()
		if err != nil {
			return err
		}
		for _, subObject := range foundList {
			if subObject.SubscriptionUrn == rs.Primary.ID {
				subscription = &subObject
			}
		}
		if subscription == nil {
			return fmt.Errorf("subscription not found")
		}

		return nil
	}
}

func testAccSMNV2SubscriptionConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_smn_topic" "topic_1" {
  name         = "%s"
  display_name = "The display name of topic_1"
}

resource "sbercloud_smn_subscription" "subscription_1" {
  topic_urn       = "${sbercloud_smn_topic.topic_1.id}"
  endpoint        = "mailtest@gmail.com"
  protocol        = "email"
  remark          = "O&M"
}

resource "sbercloud_smn_subscription" "subscription_2" {
  topic_urn       = "${sbercloud_smn_topic.topic_1.id}"
  endpoint        = "13600000000"
  protocol        = "sms"
  remark          = "O&M"
}
`, rName)
}
