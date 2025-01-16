package rds

import (
	"encoding/json"
	"fmt"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getSQLServerAccountResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.SBC_REGION_NAME
	// getSQLServerAccount: query RDS SQLServer account
	var (
		getSQLServerAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user/detail?page=1&limit=100"
		getSQLServerAccountProduct = "rds"
	)
	getSQLServerAccountClient, err := cfg.NewServiceClient(getSQLServerAccountProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and user from resource id
	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	accountName := parts[1]

	getSQLServerAccountPath := getSQLServerAccountClient.Endpoint + getSQLServerAccountHttpUrl
	getSQLServerAccountPath = strings.ReplaceAll(getSQLServerAccountPath, "{project_id}",
		getSQLServerAccountClient.ProjectID)
	getSQLServerAccountPath = strings.ReplaceAll(getSQLServerAccountPath, "{instance_id}", instanceId)

	getSQLServerAccountResp, err := pagination.ListAllItems(
		getSQLServerAccountClient,
		"page",
		getSQLServerAccountPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS SQLServer account: %s", err)
	}

	getSQLServerAccountRespJson, err := json.Marshal(getSQLServerAccountResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS SQLServer account: %s", err)
	}
	var getSQLServerAccountRespBody interface{}
	err = json.Unmarshal(getSQLServerAccountRespJson, &getSQLServerAccountRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS SQLServer account: %s", err)
	}

	account := utils.PathSearch(fmt.Sprintf("users[?name=='%s']|[0]", accountName), getSQLServerAccountRespBody, nil)

	if account != nil {
		return account, nil
	}

	return nil, fmt.Errorf("error get RDS SQLServer account by instanceID %s and account %s", instanceId, accountName)
}

func TestAccSQLServerAccount_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "sbercloud_rds_sqlserver_account.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSQLServerAccountResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSQLServerAccount_basic(name, "Test@12345678"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"sbercloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "state"),
				),
			},
			{
				Config: testSQLServerAccount_basic(name, "Test@1234567890"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"sbercloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testSQLServerAccount_basic(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "sbercloud_rds_sqlserver_account" "test" {
  instance_id = sbercloud_rds_instance.test.id
  name        = "%[2]s"
  password    = "%[3]s"
}
`, testAccRdsInstance_sqlserver(name), name, password)
}
