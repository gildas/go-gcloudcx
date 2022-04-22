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
		subject.Logger = log
	} else if len(selfURI) > 0 {
		if err := client.Get(context, selfURI, &subject); err != nil {
			return err
		}
		subject.Logger = log.Record("id", subject.ID)
	} else if len(name) > 0 {
		return errors.NotImplemented.WithStack()
	} else if _, ok := client.Grant.(*ClientCredentialsGrant); !ok { // /users/me is not possible with ClientCredentialsGrant
		if err := client.Get(context, "/users/me", &subject); err != nil {
			return err
		}
		subject.Logger = log.Record("id", subject.ID)
	}
	return nil
}

// CheckScopes checks if the subject allows or denies the given scopes
//
// See https://developer.genesys.cloud/authorization/platform-auth/scopes#scope-descriptions
func (subject AuthorizationSubject) CheckScopes(scopes ...string) (permitted []string, denied []string) {
	for _, scope := range scopes {
		authScope := AuthorizationScope{}.With(scope)
		granted := false
		for _, grant := range subject.Grants {
			if granted = grant.CheckScope(authScope); granted {
				permitted = append(permitted, scope)
				break
			}
		}
		if !granted {
			denied = append(denied, scope)
		}
	}
	return
}

// String returns a string representation of the AuthorizationSubject
//
// implements fmt.Stringer
func (subject AuthorizationSubject) String() string {
	return subject.Name
}