---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_logs"
description: |-
  Use this data source to get log list under the specified log stream within HuaweiCloud.
---
# huaweicloud_lts_logs

Use this data source to get log list under the specified log stream within HuaweiCloud.

## Example Usage

### Query all logs of the specified time period

```hcl
variable "log_group_id" {}
variable "log_stream_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_lts_logs" "test" {
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id
  start_time    = var.start_time
  end_time      = var.end_time
}
```

### Query all logs of the specified time period

```hcl
variable "log_group_id" {}
variable "log_stream_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_lts_logs" "test" {
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id
  start_time    = var.start_time
  end_time      = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.
  
* `log_group_id` - (Required, String) Specifies the ID of the log group to which the log to be queried belongs.

* `log_stream_id` - (Required, String) Specifies the ID of the log stream to which the log to be queried belongs.

* `start_time` - (Required, String) Specifies the start time of the log to be queried, in milliseconds.  
  The maximum query time range is `180` days.

* `end_time` - (Required, String) Specifies the end time of the log to be queried, in milliseconds.  
  The maximum query time range is `180` days.

* `__time__` - (Optional, String) Specifies the custom time of the log to be queried.

* `highlight` - (Optional, Bool) Specifies whether to highlight the keywords, defaults to **true**.

* `is_desc` - (Optional, Bool) Specifies whether to sort logs in descending order, defaults to **false**.

* `is_iterative` - (Optional, Bool) Specifies whether the log query is iterative, defaults to **false**.

* `keywords` - (Optional, String) Specifies the keywords of the log to be queried.  
  For more details, please following [document](https://support.huaweicloud.com/intl/en-us/usermanual-lts/lts_05_0111.html).

* `labels` - (Optional, Map) Specifies the key/value pairs of the log labels to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `logs` - The list of logs.  
  The [logs](#data_lts_logs) structure is documented below.

<a name="data_lts_logs"></a>
The `logs` block supports:

* `content` - The content of the log.

* `labels` - The labels of the log.

* `line_num` - The line number of the log.
