---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_password_policy"
description: |-
  Manages the account password policy within HuaweiCloud.
---

# huaweicloud_identity_password_policy

Manages the account password policy within HuaweiCloud.

-> You *must* have admin privileges to use this resource.  
   This resource overwrites an existing configuration, make sure one resource per account.  
   During action `terraform destroy` it sets values the same as defaults for this resource.

## Example Usage

```hcl
resource "huaweicloud_identity_password_policy" "test" {
  password_char_combination             = 4
  minimum_password_length               = 12
  number_of_recent_passwords_disallowed = 2
  password_validity_period              = 180
}
```

## Argument Reference

The following arguments are supported:

* `maximum_consecutive_identical_chars` - (Optional, Int) Specifies the maximum number of times that a character is
  allowed to consecutively present in a password.  
  The valid value is range from `0` to `32` and defaults to `0` which indicates that consecutive identical characters
  are allowed in a password.  
  For example, value `2` indicates that two consecutive identical characters are not allowed in a password.

* `minimum_password_age` - (Optional, Int) Specifies the minimum period (minutes) after which users are allowed to make
  a password change.  
  The valid value is range from `0` to `1,440` and defaults to `0`.

* `minimum_password_length` - (Optional, Int) Specifies the minimum number of characters that a password must contain.
  The valid value is range from `6` to `32` and defaults to `8`.

* `number_of_recent_passwords_disallowed` - (Optional, Int) Specifies the member of previously used passwords that are
  not allowed.  
  The valid value is range from `0` to `10` and defaults to `1`. For example, value `3` indicates that the user
  cannot set the last three passwords that the user has previously used when setting a new password.

* `password_not_username_or_invert` - (Optional, Bool) Specifies whether the password can be the username or the
  username spelled backwards.  
  Defaults to `true`, which indicates that the username or the inversion of username is not allowed to be used as a
  password.

* `password_validity_period` - (Optional, Int) Specifies the password validity period (days).
  The value ranges from `0` to `180` and defaults to `0` which indicates that this requirement does not apply.

* `password_char_combination` - (Optional, Int) Specifies the minimum number of character types that a password must
  contain.  
  The valid value is range from `2` to `4` and defaults to `2` which indicates that a password must contain at least
  two of the following: uppercase letters, lowercase letters, digits, and special characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of account password policy, which is the same as the account ID.

* `maximum_password_length` - The maximum number of characters that a password can contain.

## Import

Identity password policy can be imported using the account ID or domain ID, e.g.

```bash
$ terraform import huaweicloud_identity_password_policy.example <domain_id>
```
