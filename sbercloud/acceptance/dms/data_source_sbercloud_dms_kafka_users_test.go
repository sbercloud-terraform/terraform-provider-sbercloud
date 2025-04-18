package dms

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDmsKafkaUsers_basic(t *testing.T) {
	dataSource := "data.sbercloud_dms_kafka_users.all"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsKafkaUsers_basic(rName),
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

func testDataSourceDataSourceDmsKafkaUsers_basic(name string) string {
	password := acceptance.RandomPassword()
	return fmt.Sprintf(`
%s

data "sbercloud_dms_kafka_users" "all" {
  depends_on = [sbercloud_dms_kafka_user.test]

  instance_id = sbercloud_dms_kafka_instance.test.id
}

data "sbercloud_dms_kafka_users" "test" {
  depends_on = [sbercloud_dms_kafka_user.test]

  instance_id = sbercloud_dms_kafka_instance.test.id
  name        = sbercloud_dms_kafka_user.test.name
  description = sbercloud_dms_kafka_user.test.description
}

locals {
  test_results = data.sbercloud_dms_kafka_users.test
}

output "name_validation" {
  value = alltrue([for v in local.test_results.users[*].name : strcontains(v, sbercloud_dms_kafka_user.test.name)])
}

output "description_validation" {
  value = alltrue([for v in local.test_results.users[*].description : strcontains(v, sbercloud_dms_kafka_user.test.description)])
}
`, testAccDmsKafkaUser_basic(name, password, "test"))
}
