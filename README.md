# terraform-provider-azuredevopsext
Terraform provider For Extra Azure Devops resources/datasources 

# Resources/Data Sources
## [Resource] azuredevopsext_environment_security
Adds a person or group as Administrator/User/Reader of an Environment.
In order to add a group you need to use the group's 'originId' attribute.
In order to add a person you need to use the user's 'localId'. Using 'originId' or 'id' does not work (that's how the Azure API works).

# Configuration
In order for the provider to be configured you need to pass the same configuration as [azuredevops's configuration using PAT](https://registry.terraform.io/providers/microsoft/azuredevops/latest/docs/guides/authenticating_using_the_personal_access_token), plus an additional project_id field.

You can configure it like this:
```terraform
provider "azuredevopsext" {
  project_id            = xxxx-xxxx-xxxx                    # or env var AZDO_PROJECT_ID
  org_service_url       = "https://dev.azure.com/<example>" # or env var AZDO_ORG_SERVICE_URL
  personal_access_token = "yyyyyyyyyyyyy"                   # or env var AZDO_PERSONAL_ACCESS_TOKEN
}
```

I recommend using the environment variables option.

# Example
You can check the [examples folder](examples/) for example files, but in short you use it like this:
```terraform
resource "azuredevopsext_environment_security_access" "default" {
  member_id      = # < origin id of a group or local id* of a user >
  environment_id = # < id of an oficial azuredevops_environment.default.id. It's usually a number like 1234 >
  role_name      = # < role: one of "Administrator", "User", "Reader" >
}
```
*: for some reason the API doesn't work with the origin id nor id of the user. You need something called local id.
