package rocketmq

import (
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDmsRocketMQUserResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.SBC_REGION_NAME
	// getRocketmqUser: query DMS rocketmq user
	var (
		getRocketmqUserHttpUrl = "v2/{project_id}/instances/{instance_id}/users/{user_name}"
		getRocketmqUserProduct = "dms"
	)
	getRocketmqUserClient, err := cfg.NewServiceClient(getRocketmqUserProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DmsRocketMQUser Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<access_key>")
	}
	instanceID := parts[0]
	user := parts[1]
	getRocketmqUserPath := getRocketmqUserClient.Endpoint + getRocketmqUserHttpUrl
	getRocketmqUserPath = strings.ReplaceAll(getRocketmqUserPath, "{project_id}", getRocketmqUserClient.ProjectID)
	getRocketmqUserPath = strings.ReplaceAll(getRocketmqUserPath, "{instance_id}", instanceID)
	getRocketmqUserPath = strings.ReplaceAll(getRocketmqUserPath, "{user_name}", user)

	getRocketmqUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqUserResp, err := getRocketmqUserClient.Request("GET", getRocketmqUserPath, &getRocketmqUserOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DmsRocketMQUser: %s", err)
	}
	return utils.FlattenResponse(getRocketmqUserResp)
}

func TestAccDmsRocketMQUser_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	resourceName := "sbercloud_dms_rocketmq_user.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQUserResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQUser_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "access_key", "testRocketmqAK"),
					resource.TestCheckResourceAttr(resourceName, "secret_key", "testRocketmqSK123"),
				),
			},
			{
				Config: testDmsRocketMQUser_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "access_key", "testRocketmqAK"),
					resource.TestCheckResourceAttr(resourceName, "secret_key", "testRocketmqSK123"),
					resource.TestCheckResourceAttr(resourceName, "white_remote_address", "10.10.10.10"),
					resource.TestCheckResourceAttr(resourceName, "admin", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_topic_perm", "PUB"),
					resource.TestCheckResourceAttr(resourceName, "default_group_perm", "SUB"),
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

func testAccDmsRocketmqUser_base(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_availability_zones" "test" {}

resource "sbercloud_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = "4.8.0"
  storage_space     = 600
  vpc_id            = sbercloud_vpc.test.id
  subnet_id         = sbercloud_vpc_subnet.test.id
  security_group_id = sbercloud_networking_secgroup.test.id

  availability_zones = [
    data.sbercloud_availability_zones.test.names[0]
  ]

  flavor_id         = "c6.4u8g.cluster"
  storage_spec_code = "dms.physical.storage.high.v2"
  broker_num        = 1
  enable_acl        = true
}
`, acceptance.TestBaseNetwork(rName), rName)
}

func testDmsRocketMQUser_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_dms_rocketmq_user" "test" {
  instance_id          = sbercloud_dms_rocketmq_instance.test.id
  access_key           = "testRocketmqAK"
  secret_key           = "testRocketmqSK123"
  white_remote_address = "10.10.10.20"
  default_topic_perm   = "SUB"
  default_group_perm   = "PUB"
}
`, testAccDmsRocketmqUser_base(name))
}

func testDmsRocketMQUser_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_dms_rocketmq_user" "test" {
  instance_id          = sbercloud_dms_rocketmq_instance.test.id
  access_key           = "testRocketmqAK"
  secret_key           = "testRocketmqSK123"
  white_remote_address = "10.10.10.10"
  admin                = "false"
  default_topic_perm   = "PUB"
  default_group_perm   = "SUB"
}
`, testAccDmsRocketmqUser_base(name))
}
