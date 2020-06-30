---
layout: "sbercloud"
page_title: "SberCloud: sbercloud_identity_project_v3"
sidebar_current: "docs-sbercloud-resource-identity-project-v3"
description: |-
  Manages a Project resource.
---

# sbercloud\_identity\_project_v3

Manages a Project resource within SberCloud Identity And Access 
Management service.

Note: You _must_ have security admin privileges in your SberCloud 
cloud to use this resource. please refer to [User Management Model](
https://docs.otc.t-systems.com/en-us/usermanual/iam/iam_01_0034.html)

## Example Usage

```hcl
resource "sbercloud_identity_project_v3" "project_1" {
  name        = "eu-de_project1"
  description = "This is a test project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the project. it must start with the
    name of an existing region_ and be less than or equal to 64 characters.
    Example: ru-moscow-1_project1.

* `description` - (Optional) A description of the project.

* `domain_id` - (Optional) The domain this project belongs to. Changing this
    creates a new Project.

* `parent_id` - (Optional) The parent of this project. Changing this creates
    a new Project.

## Attributes Reference

The following attributes are exported:

* `domain_id` - See Argument Reference above.
* `parent_id` - See Argument Reference above.

## Import

Projects can be imported using the `id`, e.g.

```
$ terraform import sbercloud_identity_project_v3.project_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
