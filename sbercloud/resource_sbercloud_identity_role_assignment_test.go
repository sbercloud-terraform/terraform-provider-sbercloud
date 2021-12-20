package sbercloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/groups"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/roles"
	"github.com/huaweicloud/golangsdk/pagination"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func extractRoleAssignmentID(roleAssignmentID string) (string, string, string, string) {
	split := strings.Split(roleAssignmentID, "/")
	return split[0], split[1], split[2], split[3]
}

func TestAccIdentityV3RoleAssignment_basic(t *testing.T) {
	var role roles.Role
	var group groups.Group
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "sbercloud_identity_role_assignment.role_assignment_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3RoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3RoleAssignment_project(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3RoleAssignmentExists(resourceName, &role, &group),
					resource.TestCheckResourceAttrPtr(resourceName, "group_id", &group.ID),
					resource.TestCheckResourceAttrPtr(resourceName, "role_id", &role.ID),
					resource.TestCheckResourceAttr(resourceName, "project_id", SBC_PROJECT_ID),
				),
			},
			{
				Config: testAccIdentityV3RoleAssignment_domain(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3RoleAssignmentExists(resourceName, &role, &group),
					resource.TestCheckResourceAttrPtr(resourceName, "group_id", &group.ID),
					resource.TestCheckResourceAttrPtr(resourceName, "role_id", &role.ID),
					resource.TestCheckResourceAttr(resourceName, "domain_id", SBC_DOMAIN_ID),
				),
			},
		},
	})
}

func testAccCheckIdentityV3RoleAssignmentDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	identityClient, err := config.IdentityV3Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_identity_role_assignment" {
			continue
		}

		_, err := roles.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Role assignment still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3RoleAssignmentExists(n string, role *roles.Role, group *groups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		identityClient, err := config.IdentityV3Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud identity client: %s", err)
		}

		domainID, projectID, groupID, roleID := extractRoleAssignmentID(rs.Primary.ID)

		var opts roles.ListAssignmentsOpts
		opts = roles.ListAssignmentsOpts{
			GroupID:        groupID,
			ScopeDomainID:  domainID,
			ScopeProjectID: projectID,
		}

		pager := roles.ListAssignments(identityClient, opts)
		var assignment roles.RoleAssignment

		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			assignmentList, err := roles.ExtractRoleAssignments(page)
			if err != nil {
				return false, err
			}

			for _, a := range assignmentList {
				if a.ID == roleID {
					assignment = a
					return false, nil
				}
			}

			return true, nil
		})
		if err != nil {
			return err
		}

		g, err := groups.Get(identityClient, groupID).Extract()
		if err != nil {
			return fmt.Errorf("Group not found")
		}
		*group = *g
		r, err := roles.Get(identityClient, assignment.ID).Extract()
		if err != nil {
			return fmt.Errorf("Role not found")
		}
		*role = *r

		return nil
	}
}

func testAccIdentityV3RoleAssignment_project(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_identity_role" "role_1" {
  name = "rds_adm"
}

resource "sbercloud_identity_group" "group_1" {
  name = "%s"
}

resource "sbercloud_identity_role_assignment" "role_assignment_1" {
  role_id    = data.sbercloud_identity_role.role_1.id
  group_id   = sbercloud_identity_group.group_1.id
  project_id = "%s"
}
`, rName, SBC_PROJECT_ID)
}

func testAccIdentityV3RoleAssignment_domain(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_identity_role" "role_1" {
  name = "secu_admin"
}

resource "sbercloud_identity_group" "group_1" {
  name = "%s"
}

resource "sbercloud_identity_role_assignment" "role_assignment_1" {
  role_id    = data.sbercloud_identity_role.role_1.id
  group_id   = sbercloud_identity_group.group_1.id
  domain_id = "%s"
}
`, rName, SBC_DOMAIN_ID)
}
