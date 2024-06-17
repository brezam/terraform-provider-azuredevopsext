package internal

import (
	"context"
	"os"
	"strings"
	"terraform-provider-azuredevopsext/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	patEnvName           = "AZDO_PERSONAL_ACCESS_TOKEN"
	orgServiceUrlEnvName = "AZDO_ORG_SERVICE_URL"
	projectIdEnvName     = "AZDO_PROJECT_ID"
)

var (
	_ provider.Provider = (*azureDevopsExtProvider)(nil)
)

func NewProvider() provider.Provider {
	return &azureDevopsExtProvider{}
}

type azureDevopsExtProvider struct{}

type azureDevopsExtProviderModel struct {
	Pat           types.String `tfsdk:"personal_access_token"`
	OrgServiceUrl types.String `tfsdk:"org_service_url"`
	ProjectId     types.String `tfsdk:"project_id"`
}

// Metadata returns the provider type name.
func (p *azureDevopsExtProvider) Metadata(ctx context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "azuredevopsext"
}

// Schema defines the provider-level schema for configuration data.
func (p *azureDevopsExtProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: providerMarkDownDescription,
		Description:         providerDescription,
		Attributes: map[string]schema.Attribute{
			"personal_access_token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Azure Devops Personal Access Token.",
			},
			"org_service_url": schema.StringAttribute{
				Optional:    true,
				Description: "Azure Devops Organization Service Url.",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "Azure Devops Project Id.",
			},
		},
	}
}

// Configure defines the provider configuration and what is passed onto resource and datasources.
func (p *azureDevopsExtProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Azure Devops Ext client")

	var config azureDevopsExtProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pat := os.Getenv(patEnvName)
	orgServiceUrl := os.Getenv(orgServiceUrlEnvName)
	projectId := os.Getenv(projectIdEnvName)

	if !config.Pat.IsNull() {
		pat = config.Pat.ValueString()
	}
	if !config.OrgServiceUrl.IsNull() {
		orgServiceUrl = config.OrgServiceUrl.ValueString()
	}
	if !config.ProjectId.IsNull() {
		projectId = config.ProjectId.ValueString()
	}

	if pat == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("pat"),
			"Missing Personal Access Token", "Provider requires Azure Devops Personal Access Token",
		)
	}
	if orgServiceUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("orgServiceUrl"),
			"Missing Organization Service URL", "Provider requires Azure Devops Organization Service Url",
		)
	} else if !strings.HasPrefix(orgServiceUrl, "https://") {
		resp.Diagnostics.AddAttributeError(
			path.Root("orgServiceUrl"),
			"Invalid Organization Service URL", "Organization Service URL must start with https://",
		)
	}
	if projectId == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("projectId"),
			"Missing Project Id", "Provider requires Azure Devops Project Id",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "client", map[string]any{"pat": pat[:3] + "...", "orgServiceUrl": orgServiceUrl, "projectId": projectId})

	client_ := client.NewClient(pat, orgServiceUrl, projectId)
	resp.DataSourceData = client_
	resp.ResourceData = client_
}

// DataSources defines the data sources implemented in the provider.
func (p *azureDevopsExtProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *azureDevopsExtProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewEnvironmentSecurityResource,
	}
}
