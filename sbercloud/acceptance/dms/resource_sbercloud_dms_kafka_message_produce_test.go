package dms

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccKafkaMessageProduce_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaMessageProduce_basic(rName),
			},
		},
	})
}

func testAccKafkaMessageProduce_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_dms_kafka_message_produce" "test" {
  depends_on = [sbercloud_dms_kafka_topic.topic]

  instance_id = sbercloud_dms_kafka_instance.test.id
  topic       = sbercloud_dms_kafka_topic.topic.name
  body        = "test"

  property_list {
    name  = "KEY"
    value = "testKey"
  }

  property_list {
    name  = "PARTITION"
    value = "1"
  }
}`, testAccDmsKafkaTopic_basic(rName))
}
