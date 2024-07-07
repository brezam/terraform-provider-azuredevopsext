package internal

import (
	"context"
	"regexp"
	"slices"
	"terraform-provider-azuredevopsext/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.ResourceWithConfigure   = (*EnvironmentSecurityResource)(nil)
	_ resource.ResourceWithImportState = (*EnvironmentSecurityResource)(nil)
)

type EnvironmentSecurityResource struct {
	client *client.Client
}

type EnvironmentSecurityResourceModel struct {
	ID            types.String `tfsdk:"id"`
	MemberId      types.String `tfsdk:"member_id"`
	RoleName      types.String `tfsdk:"role_name"`
	EnvironmentId types.String `tfsdk:"environment_id"`
}

func NewEnvironmentSecurityResource() resource.Resource {
	return &EnvironmentSecurityResource{}
}

func (r *EnvironmentSecurityResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment_security_access"
}

func (r *EnvironmentSecurityResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *EnvironmentSecurityResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: environmentSecurityDescription,
		Version:     1,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
			},
			"member_id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Required:      true,
				Description:   "Member id of the user or group to add to the environment security.",
			},
			"environment_id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Required:      true,
				Description:   "Environment id to add the user or group to.",
			},
			"role_name": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Required:      true,
				Validators: []validator.String{
					stringvalidator.OneOf("Administrator", "User", "Reader"),
				},
				Description: "Role name of the user or group to add to the environment security (one of: 'Administrator', 'User', 'Reader').",
			},
		},
	}
}

func (r *EnvironmentSecurityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnvironmentSecurityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	environmentSecurityAccess, err := r.client.GetEnvironmentSecurityMembers(state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read environment security", err.Error())
		return
	}

	accessIdx := slices.IndexFunc(environmentSecurityAccess, func(c client.EnvironmentSecurityAccess) bool {
		return c.Identity.Id == state.MemberId.ValueString()
	})
	if accessIdx != -1 {
		state.RoleName = types.StringValue(environmentSecurityAccess[accessIdx].Role.Name)
	} else {
		resp.State.RemoveResource(ctx)
	}
}

func (r *EnvironmentSecurityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EnvironmentSecurityResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleNameEnum, err := client.MakeRoleName(plan.RoleName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to convert role name from string to enum", err.Error())
		return
	}

	access, err := r.client.AddMemberToEnvironmentSecurity(plan.EnvironmentId.ValueString(), plan.MemberId.ValueString(), roleNameEnum)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create environment security", err.Error())
		return
	}
	plan.ID = types.StringValue(makeIdEnvSecurity(plan.MemberId.ValueString(), roleNameEnum))
	plan.RoleName = types.StringValue(access.Role.Name)
	plan.MemberId = types.StringValue(access.Identity.Id)
	resp.State.Set(ctx, plan)
}

func (r *EnvironmentSecurityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No update, only delete + create for us
}

func (r *EnvironmentSecurityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnvironmentSecurityResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteMemberInEnvironmentSecurity(state.EnvironmentId.ValueString(), state.MemberId.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to create environment security", err.Error())
	}
}

func (r *EnvironmentSecurityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	pattern := regexp.MustCompile("^([^-]+),([^@]+)@(.+)$")
	matches := pattern.FindStringSubmatch(req.ID)
	if matches == nil {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: '<member id>,<role name>@<environment id>'")
		return
	}
	memberId := matches[1]
	roleName := matches[2]
	environmentId := matches[3]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("member_id"), memberId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_name"), roleName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), environmentId)...)
}

func makeIdEnvSecurity(memberId string, roleName client.RoleName) string {
	return memberId + "," + string(roleName)
}
