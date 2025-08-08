---
subcategory: "FunctionGraph"
layout: "sbercloud"
page_title: "SberCloud: sbercloud_fgs_function_triggers"
description: |-
  Use this data source to get the list of function triggers of FunctionGraph within SberCloud.
---

# sbercloud_fgs_function_triggers

Use this data source to get the list of function triggers of FunctionGraph within SberCloud.

## Example Usage

```hcl
variable function_urn {}

data "sbercloud_fgs_function_triggers" "test" {
  function_urn = var.function_urn
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `function_urn` - (Required, String) Specifies the function URN to which the trigger belongs.

* `trigger_id` - (Optional, String) Specifies the ID of the function trigger.

* `type` - (Optional, String) Specifies type of the function trigger.
  The valid values are as follows:
  + **TIMER**
  + **APIG**
  + **CTS**
  + **DDS**
  + **DMS**
  + **DIS**
  + **LTS**
  + **OBS**
  + **SMN**
  + **KAFKA**
  + **RABBITMQ**
  + **DEDICATEDGATEWAY**
  + **OPENSOURCEKAFKA**
  + **APIC**
  + **GAUSSMONGO**
  + **EVENTGRID**
  + **IOTDA**

* `status` - (Optional, String) Specifies status of the function trigger.
  The valid values are as follows:
  + **ACTIVE**
  + **DISABLED**

* `start_time` - (Optional, String) Specifies start time of creation time of the function trigger.
  The format is `YYYY-MM-DDThh:mm:ss{timezone}`.

* `end_time` - (Optional, String) Specifies end time of creation time of the function trigger.
  The format is `YYYY-MM-DDThh:mm:ss{timezone}`.

  -> The `status`, `start_time` and `end_time` parameters does not take effect for some triggers, e.g. `SMN`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `triggers` - All triggers that match the filter parameters.
  The [triggers](#function_triggers) structure is documented below.

<a name="function_triggers"></a>
The `triggers` block supports:

* `id` - The ID of the function trigger.

* `type` - The type of the function trigger.

* `event_data` - The detailed configuration of the function trigger.

* `status` - The current status of the function trigger.

* `created_at` - The creation time of the function trigger, in RFC3339 format.

* `updated_at` - The latest update time of the function trigger, in RFC3339 format.
