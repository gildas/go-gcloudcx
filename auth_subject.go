package gcloudcx

import (
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// AuthorizationSubject describes the roles and permissions of a Subject
type AuthorizationSubject struct {
	ID      uuid.UUID            `json:"id"`
	SelfUri string               `json:"selfUri"`
	Name    string               `json:"name"`
	Grants  []AuthorizationGrant `json:"grants"`
	Version int                  `json:"version"`
	logger  *logger.Logger       `json:"-"`
}

// Initialize initializes the object
//
// accepted parameters: *gcloufcx.Client, *logger.Logger
//
// implements Initializable
func (subject *AuthorizationSubject) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case *logger.Logger:
			subject.logger = parameter.Child("authorization_subject", "authorization_subject", "id", subject.ID)
		}
	}
}

// GetID gets the identifier
//
// implements core.Identifiable
func (subject AuthorizationSubject) GetID() uuid.UUID {
	return subject.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (subject AuthorizationSubject) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/authorization/subjects/%s", ids[0])
	}
	if subject.ID != uuid.Nil {
		return NewURI("/api/v2/authorization/subjects/%s", subject.ID)
	}
	return URI("/api/v2/authorization/subjects/")
}

// CheckScopes checks if the subject allows or denies the given scopes
//
// See https://developer.genesys.cloud/authorization/platform-auth/scopes#scope-descriptions
func (subject AuthorizationSubject) CheckScopes(scopes ...string) (permitted []string, denied []string) {
	log := subject.logger.Child(nil, "check_scopes")

	for _, scope := range scopes {
		authScope := AuthorizationScope{}.With(scope)
		granted := false
		for _, grant := range subject.Grants {
			log.Tracef("Checking against grant %s", grant)
			if granted = grant.CheckScope(authScope); granted {
				log.Debugf("Scope %s permitted by Authorization Grant %s", authScope, grant)
				permitted = append(permitted, scope)
				break
			}
		}
		if !granted {
			log.Tracef("Scope %s is denied", authScope)
			denied = append(denied, scope)
		}
	}
	return
}

// String returns a string representation of the AuthorizationSubject
//
// implements fmt.Stringer
func (subject AuthorizationSubject) String() string {
	if len(subject.Name) > 0 {
		return subject.Name
	}
	return subject.ID.String()
}
