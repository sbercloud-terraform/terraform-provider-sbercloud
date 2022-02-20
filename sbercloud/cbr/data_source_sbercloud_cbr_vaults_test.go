package cbr

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/cbr/v3/vaults"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func TestAccCbrVaultsV3_BasicServer(t *testing.T) {
	randName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_cbr_vault.test"
	dataSourceName := "data.sbercloud_cbr_vaults.test"

	var vault vaults.Vault

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCBRV3Vault_serverBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccCbrVaultsV3_serverBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.consistent_level", "app_consistent"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "backup"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "200"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.resources.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.enterprise_project_id",
						acceptance.SBC_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

func TestAccCbrVaultsV3_BasicVolume(t *testing.T) {
	randName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_cbr_vault.test"
	dataSourceName := "data.sbercloud_cbr_vaults.test"

	var vault vaults.Vault

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCBRV3Vault_volumeBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccCbrVaultsV3_volumeBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeDisk),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "backup"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "50"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.resources.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.enterprise_project_id",
						acceptance.SBC_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

func TestAccCbrVaultsV3_BasicTurbo(t *testing.T) {
	randName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "sbercloud_cbr_vault.test"
	dataSourceName := "data.sbercloud_cbr_vaults.test"

	var vault vaults.Vault

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCBRV3Vault_turboBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccCbrVaultsV3_turboBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.protection_type", "backup"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.size", "800"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.resources.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vaults.0.enterprise_project_id",
						acceptance.SBC_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

func testAccCbrVaultsV3_serverBasic(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_cbr_vaults" "test" {
  name = sbercloud_cbr_vault.test.name
}
`, testAccCBRV3Vault_serverBasic(rName))
}

func testAccCbrVaultsV3_volumeBasic(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_cbr_vaults" "test" {
  name = sbercloud_cbr_vault.test.name
}
`, testAccCBRV3Vault_volumeBasic(rName))
}

func testAccCbrVaultsV3_turboBasic(rName string) string {
	return fmt.Sprintf(`
%s

data "sbercloud_cbr_vaults" "test" {
  name = sbercloud_cbr_vault.test.name
}
`, testAccCBRV3Vault_turboBasic(rName))
}
