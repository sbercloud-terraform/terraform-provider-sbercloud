---
subcategory: "Identity and Access Management (IAM)"
---

# sbercloud_identity_projects

Use this data source to query the project list within SberCloud.

~> You *must* have IAM read privileges to use this data source.

## Example Usage

### Obtain project information by name

```hcl
data "sbercloud_identity_projects" "test" {
  name = "ru-moscow-1"
}
```

### Obtain special project information by name

```hcl
data "sbercloud_identity_projects" "test" {
  name = "MOS" // The project for OBS Billing
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the project name to query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `projects` - The details of the query projects. The structure is documented below.

The `projects` block supports:

* `id` - The project ID.

* `name` - The project name.

* `enabled` - Whether project is enabled.