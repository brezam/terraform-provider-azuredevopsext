package client

type ValueList[T any] struct {
	Count int
	Value []T
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
