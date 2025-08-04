package rds

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

func getRdsInstanceEipAssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.SBC_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances?id={instance_id}"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	publicIP := utils.PathSearch("instances|[0].public_ips[0]", getRespBody, nil)
	if publicIP == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccRdsInstanceEipAssociate_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "sbercloud_rds_instance_eip_associate.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRdsInstanceEipAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceEipAssociate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"sbercloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip",
						"sbercloud_vpc_eip.test", "address"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip_id",
						"sbercloud_vpc_eip.test", "id"),
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

func testAccRdsInstanceEipAssociate_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "sbercloud_availability_zones" "test" {}

data "sbercloud_rds_flavors" "test" {
  db_type       = "PostgreSQL"
  db_version    = "17"
  instance_mode = "single"
  //group_type    = "dedicated"
  vcpus         = 4
}

resource "sbercloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.sbercloud_rds_flavors.test.flavors[0].name
  security_group_id = sbercloud_networking_secgroup.test.id
  subnet_id         = sbercloud_vpc_subnet.test.id
  vpc_id            = sbercloud_vpc.test.id
  availability_zone = [data.sbercloud_availability_zones.test.names[0]]

  db {
    type    = "PostgreSQL"
    version = "17"
    //port    = 3306
  }

  volume {
    type = "ESSD"
    size = 40
  }
}

resource "sbercloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    share_type  = "PER"
    size        = 5
    charge_mode = "traffic"
  }
}

resource "sbercloud_rds_instance_eip_associate" "test" { 
  instance_id  = sbercloud_rds_instance.test.id
  public_ip    = sbercloud_vpc_eip.test.address
  public_ip_id = sbercloud_vpc_eip.test.id
}`, acceptance.TestBaseNetwork(rName), rName)
}
