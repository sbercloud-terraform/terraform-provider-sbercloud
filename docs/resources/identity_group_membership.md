---
subcategory: "Identity and Access Management (IAM)"
---

# sbercloud\_identity\_group\_membership

Manages a User Group Membership resource within SberCloud IAM service.

Note: You _must_ have admin privileges in your SberCloud cloud to use
this resource.

## Example Usage

```hcl
resource "sbercloud_identity_group" "group_1" {
  name        = "group1"
  description = "This is a test group"
}

resource "sbercloud_identity_user" "user_1" {
  name     = "user1"
  enabled  = true
  password = "password12345!"
}

resource "sbercloud_identity_user" "user_2" {
  name     = "user2"
  enabled  = true
  password = "password12345!"
}

resource "sbercloud_identity_group_membership" "membership_1" {
  group = "${sbercloud_identity_group.group_1.id}"
  users = ["${sbercloud_identity_user.user_1.id}",
    "${sbercloud_identity_user.user_2.id}"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required) The group ID of this membership. 

* `users` - (Required) A List of user IDs to associate to the group.

## Attributes Reference

The following attributes are exported:

* `group` - See Argument Reference above.

* `users` - See Argument Reference above.

