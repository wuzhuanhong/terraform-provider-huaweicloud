---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_policy"
description: |-
  Manages an identity policy resource within HuaweiCloud.
---
# huaweicloud_identityv5_policy

Manages an identity policy resource within HuaweiCloud.

-> You **must** have admin privileges to use this resource.

## Example Usage

```hcl
variable "policy_name" {}

resource "huaweicloud_identityv5_policy" "test" {
  name            = var.policy_name
  description     = "created by terraform"
  policy_document = jsonencode(
    {
      Statement = [
        {
          Action = ["*"]
          Effect = "Allow"
        },
      ]
      Version = "5.0"
    }
  )
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, NonUpdatable) Specifies the name of identity policy.  
  The valid length is limited from `1` to `128`.  
  Only English letters, digits, underscores (_), plus (+), equals (=), dots (.), ats (@) and hyphens (-) are allowed.

* `policy_document` - (Required, String) Specifies the policy document of the identity policy, in JSON format.  
  If updated, a new version of policy will be created and set to default version.  
  At most `5` versions of each policy are allowed, if the version limit has been reached when creating a new version,
  the oldest version will be replaced.

* `path` - (Optional, String, NonUpdatable) Specifies the resource path of the identity policy.  
  It's a part of the uniform resource name and made of several strings, each containing one or more English letters,
  digits, underscores (_), plus (+), equals (=), comma (,), dots (.), at (@) and hyphens (-), and must be ended with
  slash (/).  
  Such as **foo/bar/**.  Defaults to empty path.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the identity policy.

* `version_to_delete` - (Optional, String) Specifies the ID the policy version to be deleted, for example, **v3**.
  If specified, this version will be deleted instead of the earliest one when updating the `policy_document`.
  The value must be an existing version and can not be the default version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The identity policy ID.

* `urn` - The URN (Uniform Resource Name) of the identity policy.  
  The format is `iam::$accountID:policy:$path$policyName` where `$accountID` is IAM account ID, `$path` is the value of
  parameter `path`, `$policyName` is the value of parameter `name`.

* `type` - The type of the identity policy.

* `default_version_id` - The default version ID of the identity policy.

* `version_ids` - The version ID list of the identity policy.

  -> The version ID with smaller index in the version ID list, means it created relatively late.

* `attachment_count` - The attachment count of the identity policy.

* `created_at` - The time when the identity policy was created.

* `updated_at` - The latest update time of the identity policy.

## Import

Identity policies can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_identityv5_policy.test <id>
```
