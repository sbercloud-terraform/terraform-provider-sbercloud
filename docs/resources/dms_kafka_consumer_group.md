---
subcategory: "Distributed Message Service (DMS)"
layout: "sbercloud"
page_title: "SberCloud: sbercloud_dms_kafka_consumer_group"
description: ""
---

# sbercloud_dms_kafka_consumer_group

Manages a DMS kafka consumer group resource within SberCloud.

## Example Usage

```hcl
variable "kafka_instance_id" {}

resource "sbercloud_dms_kafka_consumer_group" "group1" {
  instance_id = var.kafka_instance_id
  name        = "group1"
  description = "Group description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DMS kafka consumer group resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the DMS kafka instance to which the consumer group belongs.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the consumer group. Changing this creates a new resource.

* `description` - (Optional, String) Specifies the description of the consumer group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is formatted `<instance_id>/<consumer_group_name>`.

* `state` - Indicates the state of the consumer group.

* `coordinator_id` - Indicates the coordinator id of the consumer group.

* `lag` - Indicates the lag number of the consumer group.

* `created_at` - Indicates the create time.

## Import

DMS kafka consumer groups can be imported using the kafka instance ID and consumer group name separated by a slash, e.g.

```bash
terraform import sbercloud_dms_kafka_user.user c8057fe5-23a8-46ef-ad83-c0055b4e0c5c/group1
```
