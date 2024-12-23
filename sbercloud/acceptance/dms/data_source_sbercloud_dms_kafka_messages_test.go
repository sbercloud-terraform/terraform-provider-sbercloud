package dms

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDmsKafkaMessages_basic(t *testing.T) {
	dataSource := "data.sbercloud_dms_kafka_messages.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsKafkaMessages_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "messages.#"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.timestamp"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.huge_message"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.message_offset"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.partition"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.size"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDmsKafkaMessages_basic(name string) string {
	startTime := strconv.FormatInt(time.Now().UnixMilli(), 10)
	endTime := strconv.FormatInt(time.Now().Add(1*time.Hour).UnixMilli(), 10)
	return fmt.Sprintf(`
%[1]s

data "sbercloud_dms_kafka_messages" "test" {
  depends_on = [sbercloud_dms_kafka_message_produce.test]

  instance_id = sbercloud_dms_kafka_instance.test.id
  topic       = sbercloud_dms_kafka_topic.topic.name
  start_time  = "%[2]s"
  end_time    = "%[3]s"
  download    = false
  partition   = 1
}

data "sbercloud_dms_kafka_messages" "by_keyword" {
  depends_on = [sbercloud_dms_kafka_message_produce.test]

  instance_id = sbercloud_dms_kafka_instance.test.id
  topic       = sbercloud_dms_kafka_topic.topic.name
  start_time  = "%[2]s"
  end_time    = "%[3]s"
  download    = false
  keyword     = sbercloud_dms_kafka_message_produce.test.body
}

output "by_keyword_validation" {
  value = length(data.sbercloud_dms_kafka_messages.by_keyword.messages) == 1
}

data "sbercloud_dms_kafka_messages" "by_offset" {
  depends_on = [sbercloud_dms_kafka_message_produce.test]
  
  instance_id    = sbercloud_dms_kafka_instance.test.id
  topic          = sbercloud_dms_kafka_topic.topic.name
  partition      = 1
  message_offset = 0
}

output "by_offset_validation" {
  value = length(data.sbercloud_dms_kafka_messages.by_offset.messages) == 1 && alltrue(
    [for v in data.sbercloud_dms_kafka_messages.by_offset.messages[*].message_offset : v == 0]
  )
}
`, testAccKafkaMessageProduce_basic(name), startTime, endTime)
}