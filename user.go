package gcloudcx

import (
	"context"
	"net/url"
	"strings"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// User describes a GCloud User
type User struct {
	ID                  uuid.UUID                `json:"id"`
	Name                string                   `json:"name,omitempty"`
	UserName            string                   `json:"username,omitempty"`
	Department          string                   `json:"department,omitempty"`
	Title               string                   `json:"title,omitempty"`
	Division            *Division                `json:"division,omitempty"`
	Mail                string                   `json:"email,omitempty"`
	Images              []*UserImage             `json:"images,omitempty"`
	PrimaryContact      []*Contact               `json:"primaryContactInfo,omitempty"`
	Addresses           []*Contact               `json:"addresses,omitempty"`
	State               string                   `json:"state,omitempty"`
	Presence            *UserPresence            `json:"presence,omitempty"`
	OutOfOffice         *OutOfOffice             `json:"outOfOffice,omitempty"`
	AcdAutoAnswer       bool                     `json:"acdAutoAnswer,omitempty"`
	RoutingStatus       *RoutingStatus           `json:"routingStatus,omitempty"`
	ProfileSkills       []string                 `json:"profileSkills,omitempty"`
	Skills              []*UserRoutingSkill      `json:"skills,omitempty"`
	Languages           []*UserRoutingLanguage   `json:"languages,omitempty"`
	LanguagePreference  string                   `json:"languagePreference,omitempty"`
	Groups              []*Group                 `json:"groups,omitempty"`
	Station             *UserStations            `json:"station,omitempty"`
	Authorization       *UserAuthorization       `json:"authorization,omitempty"`
	Employer            *EmployerInfo            `json:"employerInfo,omitempty"`
	Manager             *User                    `json:"manager,omitempty"`
	Certifications      []string                 `json:"certifications,omitempty"`
	Biography           *Biography               `json:"biography,omitempty"`
	ConversationSummary *UserConversationSummary `json:"conversationSummary,omitempty"`
	Locations           []*Location              `json:"locations,omitempty"`
	GeoLocation         *GeoLocation             `json:"geolocation,omitempty"`
	Chat                *Jabber                  `json:"chat,omitempty"`
	Version             int                      `json:"version,omitempty"`
	client              *Client                  `json:"-"`
	logger              *logger.Logger           `json:"-"`
}

// Jabber describe a Jabber ID for chats
type Jabber struct {
	ID string `json:"jabberId"`
}

// Fetch fetches a user
//
// implements Fetchable
func (user *User) Fetch(ctx context.Context, client *Client, parameters ...interface{}) error {
	id, name, selfURI, log := client.ParseParameters(ctx, user, parameters...)

	if id != uuid.Nil {
		if err := client.Get(ctx, NewURI("/users/%s", id), &user); err != nil {
			return err
		}
		user.logger = log
	} else if len(selfURI) > 0 {
		if err := client.Get(ctx, selfURI, &user); err != nil {
			return err
		}
		user.logger = log.Record("id", user.ID)
	} else if len(name) > 0 {
		return errors.NotImplemented.WithStack()
	} else if _, ok := client.Grant.(*ClientCredentialsGrant); !ok { // /users/me is not possible with ClientCredentialsGrant
		if err := client.Get(ctx, "/users/me", &user); err != nil {
			return err
		}
		user.logger = log.Record("id", user.ID)
	}
	user.client = client
	return nil
}

// GetMyUser retrieves the User that authenticated with the client
//   properties is one of more properties that should be expanded
//   see https://developer.mypurecloud.com/api/rest/v2/users/#get-api-v2-users-me
func (client *Client) GetMyUser(context context.Context, properties ...string) (*User, error) {
	query := url.Values{}
	if len(properties) > 0 {
		query.Add("expand", strings.Join(properties, ","))
	}
	user := &User{}
	if err := client.Get(context, NewURI("/users/me?%s", query.Encode()), &user); err != nil {
		return nil, err
	}
	user.client = client
	user.logger = client.Logger.Child("user", "user", "id", user.ID)
	return user, nil
}

// GetID gets the identifier of this
//
// implements Identifiable
func (user User) GetID() uuid.UUID {
	return user.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (user User) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/users/%s", ids[0])
	}
	if user.ID != uuid.Nil {
		return NewURI("/api/v2/users/%s", user.ID)
	}
	return URI("/api/v2/users/")
}

// String gets a string version
//
// implements the fmt.Stringer interface
func (user User) String() string {
	if len(user.Name) > 0 {
		return user.Name
	}
	if len(user.UserName) > 0 {
		return user.UserName
	}
	if len(user.Mail) > 0 {
		return user.Mail
	}
	return user.ID.String()
}
