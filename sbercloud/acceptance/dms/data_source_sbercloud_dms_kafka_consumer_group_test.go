package dms

import (
	"fmt"
	"testing"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDmsKafkaConsumerGroup_basic(t *testing.T) {
	dataSource := "data.sbercloud_dms_kafka_consumer_groups.all"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsKafkaConsumerGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "users.#"),
					resource.TestCheckOutput("name_validation", "true"),
					resource.TestCheckOutput("description_validation", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDmsKafkaConsumerGroup_basic(name string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_dms_kafka_consumer_groups" "all" {
  depends_on = [sbercloud_dms_kafka_consumer_group.test]

  instance_id = sbercloud_dms_kafka_consumer_group.test.instance_id
}

data "sbercloud_dms_kafka_consumer_groups" "test" {
  depends_on = [sbercloud_dms_kafka_consumer_group.test]

  instance_id = sbercloud_dms_kafka_instance.test.id
  name        = sbercloud_dms_kafka_consumer_group.test.name
  description = sbercloud_dms_kafka_consumer_group.test.description
}

locals {
  test_results = data.sbercloud_dms_kafka_consumer_groups.test
}

output "name_validation" {
  value = alltrue([for v in local.test_results.groups[*].name : strcontains(v, sbercloud_dms_kafka_consumer_group.test.name)])
}

output "description_validation" {
  value = alltrue([for v in local.test_results.groups[*].description : strcontains(v, sbercloud_dms_kafka_consumer_group.test.description)])
}
`, testAccDmsKafkaConsumerGroup_basic(name, "test"))
}
