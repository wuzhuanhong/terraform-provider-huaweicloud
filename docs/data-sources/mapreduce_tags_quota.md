---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_tags_quota"
description: |-
  Use this data source to query the tags quota information of cluster within HuaweiCloud.
---

# huaweicloud_mapreduce_tags_quota

Use this data source to query the tags quota information of cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_mapreduce_tags_quota" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the tags quota is located.  
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_quota` - The total quota of tags.

* `available_quota` - The available quota of tags.
