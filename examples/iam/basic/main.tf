
data "sbercloud_identity_role" "auth_admin" {
  name = "system_all_0"
}

resource "sbercloud_identity_user" "user_A" {
  name     = var.iden_user_name
  password = var.password
}

resource "sbercloud_identity_group" "group" {
  count = length(var.iden_group)

  name        = lookup(var.iden_group[count.index], "name", null)
  description = lookup(var.iden_group[count.index], "description", null)
}

resource "sbercloud_identity_group" "default_group" {
  name        = "default_group"
  description = "This is a default identity group."
}

resource "sbercloud_identity_group_membership" "membership_1" {
  group = length(sbercloud_identity_group.group) >= 2 ? sbercloud_identity_group.group[1].id : sbercloud_identity_group.default_group.id
  users = [sbercloud_identity_user.user_A.id]
}

resource "sbercloud_identity_role_assignment" "role_assignment_1" {
  group_id  = length(sbercloud_identity_group.group) >= 2 ? sbercloud_identity_group.group[1].id : sbercloud_identity_group.default_group.id
  domain_id = var.domain_id
  role_id   = data.sbercloud_identity_role.auth_admin.id
}
