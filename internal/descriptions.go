package internal

const (
	// provider
	providerDescription         = "Azure DevOps Extra: Unofficial provider for additional azure devops resources."
	providerMarkDownDescription = "Azure DevOps Extra: Unofficial provider for additional azure devops resources.\n" +
		"It can be configured like this:\n" +
		"```hcl\n" +
		"provider \"azuredevopsext\" {\n" +
		"    project_id            = \"xxxx-xxxx-xxxx\"                  # or env var AZDO_PROJECT_ID\n" +
		"    org_service_url       = \"https://dev.azure.com/<example>\" # or env var AZDO_ORG_SERVICE_URL\n" +
		"    personal_access_token = \"yyyyyyyyyyyyy\"                   # or env var AZDO_PERSONAL_ACCESS_TOKEN\n" +
		"}\n" +
		"```"

	// resource: environment_security
	environmentSecurityDescription        = "Environment Security Access for Azure DevOps Environment."
	environmentSecurityMarkdowDescription = "Environment Security Access.\n" +
		"This resource allows you to add an user or a group as Administrator/User/Reader of an environment.\n" +
		"It can be used like this:\n" +
		"```hcl\n" +
		"resource \"azuredevopsext_environment_security_access\" \"example\" {\n" +
		"    member_id      = \"xxxx-xxxx-xxxx\"\n" +
		"    role_name      = \"Administrator\" \n" +
		"    environment_id = \"1234\" # it should be sourced directly from an azuredevops_environment resource (from the official provider)\n" +
		"}\n" +
		"```"
)
