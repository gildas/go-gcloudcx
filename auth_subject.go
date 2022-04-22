package gcloudcx

import (
	"context"

	"github.com/gildas/go-errors"
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
	Logger  *logger.Logger       `json:"-"`
}

// GetID gets the identifier
//
// implements core.Identifiable
func (subject AuthorizationSubject) GetID() uuid.UUID {
	return subject.ID
}

func (subject *AuthorizationSubject) Fetch(context context.Context, client *Client, parameters ...interface{}) error {
	id, name, selfURI, log := client.ParseParameters(context, subject, parameters...)

	if id != uuid.Nil {
		if err := client.Get(context, NewURI("/authorization/subjects/%s", id), &subject); err != nil {
			return err
		}
	} else if len(selfURI) > 0 {
		if err := client.Get(context, selfURI, &subject); err != nil {
			return err
		}
	} else if len(name) > 0 {
		return errors.NotImplemented.WithStack()
	} else {
		return errors.CreationFailed.With("AuthorizationSubject")
	}
	subject.Logger = log.Child("authorization_subject", "authorization_subject", "id", subject.ID)
	return nil
}

// CheckScopes checks if the subject allows or denies the given scopes
//
// See https://developer.genesys.cloud/authorization/platform-auth/scopes#scope-descriptions
func (subject AuthorizationSubject) CheckScopes(scopes ...string) (permitted []string, denied []string) {
	log := subject.Logger.Child(nil, "check_scopes")

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
