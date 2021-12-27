---
subcategory: "Identity and Access Management (IAM)"
---

# sbercloud_identity_custom_role

Use this data source to get the ID of an IAM **custom policy**.

## Example Usage

```hcl
data "sbercloud_identity_custom_role" "role" {
  name = "custom_role"
}
```

## Argument Reference

* `name` - (Optional, String) Name of the custom policy.

* `id` - (Optional, String) ID of the custom policy.

* `domain_id` - (Optional, String) The domain the policy belongs to.

* `description` - (Optional, String) Description of the custom policy.

* `type` - (Optional, String) Display mode. Valid options are _AX_: Account level and _XA_: Project level.

* `references` - (Optional, Int) The number of citations for the custom policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `policy` - Document of the custom policy.

* `catalog` - The catalog of the custom policy.
