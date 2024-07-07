package client

import "fmt"

// RoleName

func MakeRoleName(s string) (RoleName, error) {
	switch s {
	case "Administrator":
		return RoleNameAdministrator, nil
	case "User":
		return RoleNameReader, nil
	case "Reader":
		return RoleNameContributor, nil
	default:
		return "", fmt.Errorf("unknown role name: %s", s)
	}
}

type RoleName string

const (
	RoleNameAdministrator RoleName = "Administrator"
	RoleNameReader        RoleName = "User"
	RoleNameContributor   RoleName = "Reader"
)

// ApiVersion

type APIVersion string

const (
	APIVersion7_1Preview1 APIVersion = "7.1-preview.1"
)
