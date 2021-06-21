---
subcategory: "Simple Message Notification (SMN)"
---

# sbercloud\_smn\_subscription

Manages SMN subscription resource within SberCloud.

## Example Usage

```hcl
resource "sbercloud_smn_topic" "topic_1" {
  name         = "topic_1"
  display_name = "The display name of topic_1"
}

resource "sbercloud_smn_subscription" "subscription_1" {
  topic_urn = sbercloud_smn_topic.topic_1.id
  endpoint  = "mailtest@gmail.com"
  protocol  = "email"
  remark    = "O&M"
}

resource "sbercloud_smn_subscription" "subscription_2" {
  topic_urn = sbercloud_smn_topic.topic_1.id
  endpoint  = "13600000000"
  protocol  = "sms"
  remark    = "O&M"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SMN subscription resource. If omitted, the provider-level region will be used. Changing this creates a new SMN subscription resource.

* `topic_urn` - (Required, String, ForceNew) Resource identifier of a topic, which is unique.

* `endpoint` - (Required, String, ForceNew) Message endpoint.
     For an HTTP subscription, the endpoint starts with http\://.
     For an HTTPS subscription, the endpoint starts with https\://.
     For an email subscription, the endpoint is a mail address.
     For an SMS message subscription, the endpoint is a phone number.

* `protocol` - (Required, String, ForceNew) Protocol of the message endpoint. Currently, email,
     sms, http, and https are supported.

* `remark` - (Optional, String, ForceNew) Remark information. The remarks must be a UTF-8-coded
     character string containing 128 bytes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `subscription_urn` - Resource identifier of a subscription, which is unique.

* `owner` - Project ID of the topic creator.

* `status` - Subscription status.
     0 indicates that the subscription is not confirmed.
     1 indicates that the subscription is confirmed.
     3 indicates that the subscription is canceled.
