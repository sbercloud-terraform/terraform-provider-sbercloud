package dms

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kafka/v2/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDmsKafkaConsumerGroupFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.HcDmsV2Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	// Split instance_id and user from resource id
	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<group_name>")
	}
	instanceId := parts[0]
	instanceConsumerGroup := parts[1]

	// List all instance groups
	request := &model.ListInstanceConsumerGroupsRequest{
		InstanceId: instanceId,
	}

	response, err := client.ListInstanceConsumerGroups(request)
	if err != nil {
		return nil, fmt.Errorf("error listing DMS kafka consumer groups in %s, error: %s", instanceId, err)
	}
	if response.Groups != nil && len(*response.Groups) != 0 {
		groups := *response.Groups
		for _, group := range groups {
			if *group.GroupId == instanceConsumerGroup {
				return group, nil
			}
		}
	}

	return nil, fmt.Errorf("can not found DMS kafka consumer group")
}

func TestAccDmsKafkaConsumerGroup_basic(t *testing.T) {
	var consumerGroup model.GroupInfoSimple
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_dms_kafka_consumer_group.test"
	description := "add desсription"
	descriptionUpdate := ""

	rc := acceptance.InitResourceCheck(
		resourceName,
		&consumerGroup,
		getDmsKafkaConsumerGroupFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsKafkaConsumerGroup_basic(rName, description),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccDmsKafkaConsumerGroup_basic(rName, descriptionUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdate),
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

func testAccDmsKafkaConsumerGroup_basic(rName, description string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_dms_kafka_consumer_group" "test" {
  instance_id = sbercloud_dms_kafka_instance.test.id
  name        = "%[2]s"
  description = "%[3]s"
}
`, testAccDmsKafkaInstance_basic(rName), rName, description)
}
