---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "azuredevopsext_environment_pipeline_permission Resource - terraform-provider-azuredevopsext"
subcategory: ""
description: |-
  Environment Pipeline Permission for Azure DevOps Environment.
---

# azuredevopsext_environment_pipeline_permission (Resource)

Environment Pipeline Permission for Azure DevOps Environment.

## Example Usage

```terraform
resource "azuredevopsext_environment_pipeline_permission" "example" {
  pipeline_id    = "1234"
  environment_id = "1234" # id of azuredevops_environment resource (official azuredevops provider)
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) Environment id where we want to add the pipeline as authorized.
- `pipeline_id` (String) Pipeline id to add to environment permissions.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
pipeline_id='1234'    # pipeline id
environment_id='1234' # environment id

terraform import azuredevopsext_environment_security_access "${pipeline_id}@${environment_id}"
```
