---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "e2e_vpc Resource - terraform-provider-e2e"
subcategory: ""
description: |-
  
---

# e2e_vpc (Resource)

Provides an e2e vpc resource.
This resource allows you to manage vpc on your e2e clusters. When applied, a new vpc is created. When destroyed, this vpc is removed.

<!-- schema generated by tfplugindocs -->
## Example Usage
```hcl
 resource "e2e_vpc" "vpc1" {
	name              = "vpc_name"
    region            = "Delhi"
 }
```
## Schema

### Required

- `region` (String) Region should specified
- `vpc_name` (String)

### Optional

- `network_size` (Number)

### Read-Only

- `created_at` (String)
- `gateway_ip` (String)
- `id` (String) The ID of this resource.
- `ipv4_cidr` (String)
- `is_active` (Boolean)
- `network_id` (Number) The id of network
- `pool_size` (Number)
- `state` (String)