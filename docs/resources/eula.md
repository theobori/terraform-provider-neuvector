---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "neuvector_eula Resource - terraform-provider-neuvector"
subcategory: ""
description: |-
  
---

# neuvector_eula (Resource)



## Example Usage

```terraform
resource "neuvector_eula" "test" {
    accepted = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `accepted` (Boolean) Accept the EULA.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import neuvector_eula.name 0
```
