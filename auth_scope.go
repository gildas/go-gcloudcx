package gcloudcx

import "strings"

// AuthorizationScope represents a scope for a client
//
// See https://developer.genesys.cloud/authorization/platform-auth/scopes#scope-descriptions
type AuthorizationScope struct {
	Domain string
	Entity string
	Action string
}

func (scope AuthorizationScope) With(subscopes ...string) AuthorizationScope {
	newScope := AuthorizationScope{"*", "*", "*"}
	expanded := []string{}
	for _, subscope := range subscopes {
		expanded = append(expanded,  strings.Split(subscope, ":")...)
	}
	if len(expanded) > 0 {
		newScope.Domain = expanded[0]
		if len(expanded) > 1 {
			newScope.Entity = expanded[1]
			if len(expanded) > 2 {
				newScope.Action = expanded[2]
			}
		}
	}
	return newScope
}

func (scope AuthorizationScope) String() string {
	return scope.Domain + ":" + scope.Entity + ":" + scope.Action
}