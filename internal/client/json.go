package client

type ISO8601Time string

type EnvironmentSecurityAccessList struct {
	Count int                         `json:"count"`
	Value []EnvironmentSecurityAccess `json:"value"`
}

type EnvironmentSecurityAccess struct {
	Identity          *AzureIdentity           `json:"identity"`
	Role              *EnvironmentSecurityRole `json:"role"`
	Access            string                   `json:"access"`
	AccessDisplayName string                   `json:"accessDisplayName"`
}

type AzureIdentity struct {
	DisplayName string `json:"displayName"`
	Id          string `json:"id"`
	UniqueName  string `json:"uniqueName"`
	Descritpor  string `json:"descriptor,omitempty"`
}

type EnvironmentSecurityRole struct {
	DisplayName      string `json:"displayName"`
	Name             string `json:"name"`
	AllowPermissions int    `json:"allowPermissions"`
	DenyPermissions  int    `json:"denyPermissions"`
	Identifier       string `json:"identifier"`
	Description      string `json:"description"`
	Scope            string `json:"scope"`
}

type CreateEnvironmentSecurityAccess struct {
	RoleName RoleName `json:"roleName"`
	UserId   string   `json:"userId"`
}

type PipelinePermissionList struct {
	Resource  *PipelineResource    `json:"resource,omitempty"`
	Pipelines []PipelinePermission `json:"pipelines"`
}

type PipelineResource struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type PipelinePermission struct {
	Id           int                             `json:"id"`
	Authorized   bool                            `json:"authorized"`
	AuthorizedBy *PipelinePermissionAuthorizedBy `json:"authorizedBy,omitempty"`
	AuthorizedOn ISO8601Time                     `json:"authorizedOn,omitempty"`
}

type PipelinePermissionAuthorizedBy struct {
	Identity *AzureIdentity `json:"identity"`
}
