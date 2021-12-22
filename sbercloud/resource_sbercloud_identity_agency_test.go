package sbercloud

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/identity/v3/agency"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccIdentityAgency_basic(t *testing.T) {
	var agency agency.Agency

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_identity_agency.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityAgencyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityAgency_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityAgencyExists(resourceName, &agency),
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
					testAccCheckIdentityAgencyExists(resourceName, &agency),
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

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_identity_agency.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityAgencyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityAgency_domain(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityAgencyExists(resourceName, &agency),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a test agency"),
					resource.TestCheckResourceAttr(resourceName, "delegated_domain_name", SBC_DOMAIN_NAME),
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
					testAccCheckIdentityAgencyExists(resourceName, &agency),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a updated test agency"),
					resource.TestCheckResourceAttr(resourceName, "delegated_domain_name", SBC_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "duration", "FOREVER"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "2"),
				),
			},
		},
	})
}

func testAccCheckIdentityAgencyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.IAMV3Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud IAM client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_identity_agency" {
			continue
		}

		v, err := agency.Get(client, rs.Primary.ID).Extract()
		if err == nil && v.ID == rs.Primary.ID {
			return fmt.Errorf("Identity Agency <%s> still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckIdentityAgencyExists(n string, ag *agency.Agency) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		client, err := config.IAMV3Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud Identity Agency: %s", err)
		}

		found, err := agency.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Identity Agency <%s> not found", rs.Primary.ID)
		}
		ag = found

		return nil
	}
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
`, rName, SBC_DOMAIN_NAME)
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
`, rName, SBC_DOMAIN_NAME)
}
