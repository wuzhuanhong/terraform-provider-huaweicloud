---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_mysql_database_privilege"
description: ""
---

# huaweicloud_rds_mysql_database_privilege

Manages RDS Mysql database privilege resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "db_name" {}
variable "user_name_1" {}
variable "user_name_2" {}

resource "huaweicloud_rds_mysql_database_privilege" "test" {
  instance_id = var.instance_id
  db_name     = var.db_name

  users {
    name     = var.user_name_1
    readonly = true
  }

  users {
    name     = var.user_name_2
    readonly = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region where the database and users (accounts) are located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the MySQL instance.

* `db_name` - (Required, String, NonUpdatable) Specifies the database name to which the users (accounts) are privileged.

* `users` - (Required, List) Specifies the user (account) permissions with the database.  
  The [users](#rds_mysql_database_privilege_users) structure is documented below.

<a name="rds_mysql_database_privilege_users"></a>
The `users` block supports:

* `name` - (Required, String) Specifies the username of the database account.

* `readonly` - (Optional, Bool) Specifies whether the user has read-only permission.  
  The valid values are as follows:
  + **true**: The database grants the current user **read-only** permission.
  + **false**: The database grants the current user **read-and-write** permission.

  The default value is **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of database privilege which is formatted `<instance_id>/<db_name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

RDS database privilege can be imported using the `instance id` and `db_name`, e.g.

```bash
$ terraform import huaweicloud_rds_mysql_database_privilege.test <instance_id>/<db_name>
```

~> During the import process, all privileges under the database and managed remotely will be synchronized to the
   `terraform.tfstate` file.
