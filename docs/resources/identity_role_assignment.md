---
subcategory: "Identity and Access Management (IAM)"
---

# sbercloud\_identity\_role\_assignment

Manages a V3 Role assignment within group on SberCloud IAM Service.

Note: You _must_ have admin privileges in your SberCloud cloud to use
this resource. 

## Example Usage: Assign Role On Project Level

```hcl
resource "sbercloud_identity_group" "group_1" {
  name = "group_1"
}

data "sbercloud_identity_role" "role_1" {
  name = "system_all_4" #ECS admin
}

resource "sbercloud_identity_role_assignment" "role_assignment_1" {
  group_id   = "${sbercloud_identity_group.group_1.id}"
  project_id = "${var.project_id}"
  role_id    = "${data.sbercloud_identity_role.role_1.id}"
}
```

## Example Usage: Assign Role On Domain Level

```hcl

variable "domain_id" {
  default     = "01aafcf63744d988ebef2b1e04c5c34"
  description = "this is the domain id"
}

resource "sbercloud_identity_group" "group_1" {
  name = "group_1"
}

data "sbercloud_identity_role" "role_1" {
  name = "secu_admin" #security admin
}

resource "sbercloud_identity_role_assignment" "role_assignment_1" {
  group_id  = "${sbercloud_identity_group.group_1.id}"
  domain_id = "${var.domain_id}"
  role_id   = "${data.sbercloud_identity_role.role_1.id}"
}

```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Optional; Required if `project_id` is empty) The domain to assign the role in.

* `group_id` - (Optional; Required if `user_id` is empty) The group to assign the role to.

* `project_id` - (Optional; Required if `domain_id` is empty) The project to assign the role in.

* `role_id` - (Required) The role to assign.

## Attributes Reference

The following attributes are exported:

* `domain_id` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `group_id` - See Argument Reference above.
* `role_id` - See Argument Reference above.
