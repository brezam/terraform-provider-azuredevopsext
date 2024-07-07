package internal

import (
	"context"
	"regexp"
	"slices"
	"strconv"
	"terraform-provider-azuredevopsext/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.ResourceWithConfigure   = (*EnvironmentPermissionResource)(nil)
	_ resource.ResourceWithImportState = (*EnvironmentPermissionResource)(nil)
)

type EnvironmentPermissionResource struct {
	client *client.Client
}

type EnvironmentPermissionResourceModel struct {
	ID            types.String `tfsdk:"id"`
	PipelineId    types.String `tfsdk:"pipeline_id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
}

func NewEnvironmentPermissionResource() resource.Resource {
	return &EnvironmentPermissionResource{}
}

func (r *EnvironmentPermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment_pipeline_permission"
}

func (r *EnvironmentPermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *EnvironmentPermissionResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         environmentPermissionDescription,
		Version:             1,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
			},
			"environment_id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Required:      true,
				Description:   "Environment id where we want to add the pipeline as authorized.",
			},
			"pipeline_id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Required:      true,
				Description:   "Pipeline id to add to environment permissions.",
			},
		},
	}
}

func (r *EnvironmentPermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnvironmentPermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	pipelinePermissions, err := r.client.ListEnvironmentPipelinePermissions(state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read pipeline permissions", err.Error())
		return
	}

	pipeIdIdx := slices.IndexFunc(pipelinePermissions, func(c client.PipelinePermission) bool {
		return strconv.Itoa(c.Id) == state.PipelineId.ValueString()
	})
	if pipeIdIdx == -1 {
		resp.State.RemoveResource(ctx)
	}
}

func (r *EnvironmentPermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EnvironmentPermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.AddEnvironmentPipelinePermission(plan.EnvironmentId.ValueString(), plan.PipelineId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to create environment security", err.Error())
		return
	}
	plan.ID = types.StringValue(plan.PipelineId.ValueString())
	resp.State.Set(ctx, plan)
}

func (r *EnvironmentPermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No update, only delete + create for us
}

func (r *EnvironmentPermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnvironmentPermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.RemoveEnvironmentPipelinePermission(state.EnvironmentId.ValueString(), state.PipelineId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete environment security", err.Error())
	}
}

func (r *EnvironmentPermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	pattern := regexp.MustCompile("^([^-]+)@(.+)$")
	matches := pattern.FindStringSubmatch(req.ID)
	if matches == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: '<pipeline id>@<environment id>'")
		return
	}
	pipelineId := matches[1]
	environmentId := matches[2]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("pipeline_id"), pipelineId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), environmentId)...)
}
