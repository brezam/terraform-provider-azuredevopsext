resource "azuredevopsext_environment_security_access" "example" {
  member_id      = "xxxx-xxxx-xxxx"
  environment_id = "1234" # id of azuredevops_environment resource (official azuredevops provider)
  role_name      = "Administrator"
}
