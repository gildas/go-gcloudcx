package purecloud

// UserAuthorization  desribes authorizations for a User
type UserAuthorization struct {
	Roles              []*DomainRole               `json:"roles"`
	UnusedRoles        []*DomainRole               `json:"unusedRoles"`
	Permissions        []string                    `json:"permissions"`
	PermissionPolicies []*ResourcePermissionPolicy `json:"permissionPolicies"`
}