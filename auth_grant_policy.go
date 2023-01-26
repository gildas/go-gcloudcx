package gcloudcx

import (
	"fmt"
	"strings"
)

type AuthorizationGrantPolicy struct {
	EntityName string   `json:"entityName"`
	Domain     string   `json:"domain"`
	Condition  string   `json:"condition"`
	Actions    []string `json:"actions"`
}

// CheckScope checks if the grant allows or denies the given scope
func (policy AuthorizationGrantPolicy) CheckScope(scope AuthorizationScope) bool {
	if policy.Domain == "*" || (policy.Domain == scope.Domain && (policy.EntityName == "*" || policy.EntityName == scope.Entity)) {
		for _, action := range policy.Actions {
			if action == "*" || action == scope.Action {
				return true
			}
		}
	}
	return false
}

// String returns a string representation of the AuthorizationDivision
//
// implements fmt.Stringer
func (policy AuthorizationGrantPolicy) String() string {
	if len(policy.Condition) == 0 {
		return fmt.Sprintf("%s:%s:[%s]", policy.Domain, policy.EntityName, strings.Join(policy.Actions, ","))
	}
	return fmt.Sprintf("%s:%s:[%s]{? %s}", policy.Domain, policy.EntityName, strings.Join(policy.Actions, ","), policy.Condition)
}
