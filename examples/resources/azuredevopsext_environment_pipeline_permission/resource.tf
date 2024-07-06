resource "azuredevopsext_environment_pipeline_permission" "example" {
  pipeline_id    = "1234"
  environment_id = "1234" # id of azuredevops_environment resource (official azuredevops provider)
}
