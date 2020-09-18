---
subcategory: "Identity and Access Management (IAM)"
layout: "sbercloud"
page_title: "SberCloud: sbercloud_identity_group_membership_v3"
sidebar_current: "docs-sbercloud-resource-identity-group-membership-v3"
description: |-
  Manages the membership combine User Group resource and User resource  within
  SberCloud IAM service.
---

# sbercloud\_identity\_group_membership_v3

Manages a User Group Membership resource within SberCloud IAM service.

Note: You _must_ have admin privileges in your SberCloud cloud to use
this resource.

## Example Usage

```hcl
resource "sbercloud_identity_group_v3" "group_1" {
  name        = "group1"
  description = "This is a test group"
}

resource "sbercloud_identity_user_v3" "user_1" {
  name     = "user1"
  enabled  = true
  password = "password12345!"
}

resource "sbercloud_identity_user_v3" "user_2" {
  name     = "user2"
  enabled  = true
  password = "password12345!"
}

resource "sbercloud_identity_group_membership_v3" "membership_1" {
  group = "${sbercloud_identity_group_v3.group_1.id}"
  users = ["${sbercloud_identity_user_v3.user_1.id}",
    "${sbercloud_identity_user_v3.user_2.id}"
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

