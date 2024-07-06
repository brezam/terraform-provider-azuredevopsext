provider "azuredevopsext" {
  project_id            = "xxxx-xxxx-xxxx"                  # Also available as environment variable AZDO_PROJECT_ID
  org_service_url       = "https://dev.azure.com/<example>" # Also available as environment variable AZDO_ORG_SERVICE_URL
  personal_access_token = "xxxx-xxxx-xxxx"                  # Also available as environment variable AZDO_PERSONAL_ACCESS_TOKEN
}
