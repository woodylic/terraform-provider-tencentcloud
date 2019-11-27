---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_nat_gateway"
sidebar_current: "docs-tencentcloud-resource-nat_gateway"
description: |-
  Provides a resource to create a NAT gateway.
---

# tencentcloud_nat_gateway

Provides a resource to create a NAT gateway.

## Example Usage

```hcl
resource "tencentcloud_nat_gateway" "foo" {
  name             = "test_nat_gateway"
  vpc_id           = "vpc-4xxr2cy7"
  bandwidth        = 100
  max_connection   = 1000000
  assigned_eip_set = ["1.1.1.1"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the NAT gateway.
* `vpc_id` - (Required, ForceNew) ID of the vpc.
* `assigned_eip_set` - (Optional) EIP IP address set bound to the gateway. The value of at least 1 and at most 10.
* `bandwidth` - (Optional) The maximum public network output bandwidth of NAT gateway (unit: Mbps), the available values include: 20,50,100,200,500,1000,2000,5000. Default is 100.
* `max_concurrent` - (Optional) The upper limit of concurrent connection of NAT gateway, the available values include: 1000000,3000000,10000000. Default is 1000000.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_time` - Create time of the NAT gateway.


## Import

NAT gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_nat_gateway.foo nat-1asg3t63
```

