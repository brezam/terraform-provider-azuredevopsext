terraform {
  required_version = "1.2.9"

  backend "local" {}

  required_providers {
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = "1.1.1"
    }
    azuredevopsext = {
      source  = "registry.terraform.io/brezam/azuredevopsext"
      version = "0.0.1"
    }
  }
}

provider "azuredevops" {}

provider "azuredevopsext" {
  project_id = local.project_id # Also available as environment variable AZDO_PROJECT_ID
}


locals {
  project_id = "xxxx-xxxx-xxxx"
  member_id  = "yyyy-yyyy-yyyy"
}

resource "azuredevops_environment" "default" {
  project_id = local.project_id
  name       = "testing-123"
}

/*
This will add the member_id (either the localId of an user, or the originId of an azure devops group) as 
an Administrator to the environment
*/
resource "azuredevopsext_environment_security_access" "default" {
  member_id      = local.member_id
  environment_id = azuredevops_environment.default.id
  role_name      = "Administrator"
}
