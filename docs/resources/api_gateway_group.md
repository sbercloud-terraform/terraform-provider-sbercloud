---
subcategory: "API Gateway (APIG)"
---

# sbercloud\_api\_gateway\_group

Provides an API gateway group resource.

## Example Usage

```hcl
resource "sbercloud_api_gateway_group" "apigw_group" {
  group_id    = "id"
  name        = "apigw_group"
  description = "your descpiption"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region where the dedicated instance is located.  
  If omitted, the provider-level region will be used. Changing this creates a new gateway group resource.

* `instance_id` - (Required, String, ForceNew) The ID of the dedicated instance to which the group belongs.

* `name` - (Required, String) The group name.

* `description` - (Optional, String) The group description.

* `environment` - (Optional, Set) The array of one or more environments of the associated group. 

Each `environment` block supports the following:

* `environment_id` - (Required, String) The ID of the environment to which the variables belongs.
* `variable` - (Required, Set) The array of one or more environment variables. Each `variable` block supports the following:
    * `name` - (Required, String) The variable name.
    * `value` - (Required, String) The variable value.
    * `id` - (Computed, String) The ID of the variable that the group has.
    * `variable_id` - (Computed, String, Deprecated) The ID of the variable that the group has.  
      _Deprecated: Use `id` instead._

* `url_domains` - (Optional, Set, MaxItems: 5) The associated domain information of the group. Each `url_domains` block supports the following:
    * `name` - (Required, String) The associated domain name.
    * `min_ssl_version` - (Optional, String) The minimum SSL protocol version.
    * `is_http_redirect_to_https` - (Optional, Bool) Whether to enable redirection from HTTP to HTTPS.

* `domain_access_enabled` - (Optional, Bool) Specifies whether to use the debugging domain name to access the APIs within the group.

* `force_destroy` - (Optional, Bool) Whether to delete all sub-resources (for API) from this group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the API group.
* `status` - Status of the API group.
* `created_at` - The creation time of the group, in RFC3339 format.
* `updated_at` - The latest update time of the group, in RFC3339 format.
* `update_time` - The latest update time of the group.


## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 10 minute.

