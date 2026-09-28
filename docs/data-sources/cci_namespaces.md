---
subcategory: "Cloud Container Instance (CCI)"
layout: "sbercloud"
page_title: "Sbercloud: sbercloud_cci_namespaces"
description: |-
  Use this data source to get the list of CCI v2 namespaces within Sbercloud.
---

# sbercloud_cci_namespaces

Use this data source to get the list of CCI v2 namespaces within Sbercloud.

## Example Usage

```hcl
variable "namespace_name" {}

data "sbercloud_cci_namespaces" "test" {
  name = var.namespace_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the CCI namespaces.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the CCI namespace.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `namespaces` - The list of CCI namespaces.
  The [namespaces](#cci_namespaces) structure is documented below.

<a name="cci_namespaces"></a>
The `namespaces` block supports:

* `name` - The name of the namespace.

* `api_version` - The API version of the namespace.

* `kind` - The kind of the namespace.

* `annotations` - The annotations of the namespace.

* `labels` - The labels of the namespace.

* `creation_timestamp` - The creation timestamp of the namespace.

* `finalizers` - The finalizers of the namespace.

* `resource_version` - The resource version of the namespace.

* `uid` - The uid of the namespace.

* `status` - The status of the namespace.
