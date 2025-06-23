---
subcategory: "Distributed Message Service (DMS)"
layout: "sbercloud"
page_title: "SberCloud: sbercloud_dms_kafka_consumer_groups"
description: |-
  Use this data source to get the list of Kafka instance consumer groups.
---

# sbercloud_dms_kafka_consumer_groups

Use this data source to get the list of Kafka instance consumer groups.

## Example Usage

### Get all groups for an instance

```hcl
variable "instance_id" {}

data "sbercloud_dms_kafka_consumer_groups" "test" {
  instance_id = var.instance_id
}
```

### Get specific group for an instance

```hcl
variable "instance_id" {}

data "sbercloud_dms_kafka_consumer_groups" "test" {
  instance_id = var.instance_id
  name        = "test_group"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `name` - (Optional, String) Specifies the group name.

* `description` - (Optional, String) Specifies the group description.

* `lag` - (Optional, Int) Specifies the number of accumulated messages.

* `coordinator_id` - (Optional, Int) Specifies the coordinator ID.

* `state` - (Optional, String) Specifies the consumer group status.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - Indicates the groups list.

  The [groups](#groups_struct) structure is documented below.

<a name="groups_struct"></a>
The `groups` block supports:

* `name` - Indicates the consumer group name.

* `description` - Indicates the consumer group description.

* `lag` - Indicates the number of accumulated messages.

* `coordinator_id` - Indicates the coordinator ID.

* `state` - Indicates the consumer group status.

* `created_at` - Indicates the create time.

* `assignment_strategy` - Indicates the partition assignment strategy.

* `members` - Indicates the consumer group members

  The [members](#members_struct) structure is documented below.

* `group_message_offsets` - Indicates the group message offsets.

  The [group_message_offsets](#group_message_offsets_struct) structure is documented below.

<a name="members_struct"></a>

The `members` block supports:

* `host` - Indicates the consumer address.

* `member_id` - Indicates the member ID.

* `client_id` - Indicates the client ID.

* `assignment` - Indicates the details about the partition assigned to the consumer.

  The [assignment](#assignment_struct) structure is documented below.

<a name="group_message_offsets_struct"></a>

The `group_message_offsets` block supports:

* `partition` - Indicates the partition.

* `lag` - Indicates the number of accumulated messages.

* `topic` - Indicates the topic name.

* `message_current_offset` - Indicates the message current offset.

* `message_log_end_offset` - Indicates the message log end offset.

<a name="assignment_struct"></a>

The `assignment` block supports:

* `topic` - Indicates the topic name.

* `partitions` - Indicates the partitions.
