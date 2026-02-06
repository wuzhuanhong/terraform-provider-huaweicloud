---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_access_keys"
description: |-
  Use this data source to get permanent access key list under the specified IAM user within HuaweiCloud.
---

# huaweicloud_identityv5_access_keys

Use this data source to get permanent access key list under the specified IAM user within HuaweiCloud.

## Example Usage

```hcl
variable "user_id" {}

data "huaweicloud_identityv5_access_keys" "test" {
  user_id = var.user_id
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, String) Specifies the ID of the IAM user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `access_keys` - The list of permanent access keys.  
  The [access_keys](#v5_access_keys) structure is documented below.

<a name="v5_access_keys"></a>
The `access_keys` block supports:

* `access_key_id` - The ID of the permanent access key (AK).

* `user_id` - The ID of the IAM user.

* `created_at` - The creation time of the access key.

* `status` - The status of the access key.
  + **active**
  + **inactive**
