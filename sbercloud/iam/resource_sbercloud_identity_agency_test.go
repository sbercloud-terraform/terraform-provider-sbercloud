package iam

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"

	"github.com/chnsz/golangsdk/openstack/identity/v3/agency"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getIdentityAgencyResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMV3Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("Error creating SberCloud IAM client: %s", err)
	}
	return agency.Get(client, state.Primary.ID).Extract()
}

func TestAccIdentityAgency_basic(t *testing.T) {
	var agency agency.Agency
	rName := acceptance.RandomAccResourceName()
	resourceName := "sbercloud_identity_agency.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&agency,
		getIdentityAgencyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityAgency_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a test agency"),
					resource.TestCheckResourceAttr(resourceName, "delegated_service_name", "op_svc_evs"),
					resource.TestCheckResourceAttr(resourceName, "duration", "FOREVER"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityAgency_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a updated test agency"),
					resource.TestCheckResourceAttr(resourceName, "delegated_service_name", "op_svc_evs"),
					resource.TestCheckResourceAttr(resourceName, "duration", "FOREVER"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "2"),
				),
			},
		},
	})
}

func TestAccIdentityAgency_domain(t *testing.T) {
	var agency agency.Agency
	rName := acceptance.RandomAccResourceName()
	resourceName := "sbercloud_identity_agency.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&agency,
		getIdentityAgencyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityAgency_domain(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a test agency"),
					resource.TestCheckResourceAttr(resourceName, "delegated_domain_name", acceptance.SBC_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "duration", "FOREVER"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityAgency_domainUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a updated test agency"),
					resource.TestCheckResourceAttr(resourceName, "delegated_domain_name", acceptance.SBC_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "duration", "FOREVER"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "2"),
				),
			},
		},
	})
}

func testAccIdentityAgency_basic(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_identity_agency" "test" {
  name                   = "%s"
  description            = "This is a test agency"
  delegated_service_name = "op_svc_evs"

  domain_roles = [
    "Tenant Administrator",
  ]
}
`, rName)
}

func testAccIdentityAgency_update(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_identity_agency" "test" {
  name                   = "%s"
  description            = "This is a updated test agency"
  delegated_service_name = "op_svc_evs"

  domain_roles = [
    "Tenant Administrator", "KMS Administrator",
  ]
}
`, rName)
}

func testAccIdentityAgency_domain(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_identity_agency" "test" {
  name                  = "%s"
  description           = "This is a test agency"
  delegated_domain_name = "%s"

  domain_roles = [
    "DAYU Administrator",
  ]
}
`, rName, acceptance.SBC_DOMAIN_NAME)
}

func testAccIdentityAgency_domainUpdate(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_identity_agency" "test" {
  name                  = "%s"
  description           = "This is a updated test agency"
  delegated_domain_name = "%s"

  domain_roles = [
    "DAYU Administrator",
    "VPC Administrator",
  ]
}
`, rName, acceptance.SBC_DOMAIN_NAME)
}
