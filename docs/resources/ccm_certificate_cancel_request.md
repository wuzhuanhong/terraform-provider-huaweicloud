---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_certificate_cancel_request"
description: |-
  Manages a CCM SSL certificate cancel request resource within HuaweiCloud.
---

# huaweicloud_ccm_certificate_cancel_request

Manages a CCM SSL certificate cancel request resource within HuaweiCloud.

-> Destroying this resource will not clear the cancellation status of the CCM SSL certificate, but will only remove
the resource information from the tfstate file.

## Example Usage

```hcl
variable "certificate_id" {}

resource "huaweicloud_ccm_certificate_cancel_request" "test" {
  certificate_id = var.certificate_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `certificate_id` - (Required, String, NonUpdatable) Specifies the CCM SSL certificate ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also the CCM SSL certificate ID).
