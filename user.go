package gcloudcx

import (
	"context"
	"encoding/json"
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
	PreferredName       string                   `json:"preferredName,omitempty"`
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

// GetMyUser retrieves the User that authenticated with the client
//
//	properties is one of more properties that should be expanded
//	see https://developer.mypurecloud.com/api/rest/v2/users/#get-api-v2-users-me
func (client *Client) GetMyUser(context context.Context, properties ...string) (*User, error) {
	query := url.Values{}
	if len(properties) > 0 {
		query.Add("expand", strings.Join(properties, ","))
	}
	user := &User{}
	if _, err := client.Get(context, NewURI("/users/me?%s", query.Encode()), &user); err != nil {
		return nil, err
	}
	user.client = client
	user.logger = client.Logger.Child("user", "user", "id", user.ID)
	return user, nil
}

// Initialize initializes the object
//
// accepted parameters: *gcloudcx.Client, *logger.Logger
//
// implements Initializable
func (user *User) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case uuid.UUID:
			user.ID = parameter
		case *Client:
			user.client = parameter
		case *logger.Logger:
			user.logger = parameter.Child("user", "user", "id", user.ID)
		}
	}
	if user.logger == nil {
		user.logger = logger.Create("gcloudcx", &logger.NilStream{})
	}
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

// Redact redacts sensitive data
//
// implements logger.Redactable
func (user User) Redact() interface{} {
	redacted := user
	if len(user.Name) > 0 {
		redacted.Name = logger.RedactWithHash(user.Name)
	}
	if len(user.UserName) > 0 {
		redacted.UserName = logger.RedactWithHash(user.UserName)
	}
	if len(user.Mail) > 0 {
		redacted.Mail = logger.RedactWithHash(user.Mail)
	}
	if len(user.PrimaryContact) > 0 {
		redacted.PrimaryContact = make([]*Contact, len(user.PrimaryContact))
		for i, contact := range user.PrimaryContact {
			redactedContact := contact.Redact().(Contact)
			redacted.PrimaryContact[i] = &redactedContact
		}
	}
	if len(user.Addresses) > 0 {
		redacted.Addresses = make([]*Contact, len(user.Addresses))
		for i, contact := range user.Addresses {
			redactedContact := contact.Redact().(Contact)
			redacted.Addresses[i] = &redactedContact
		}
	}
	if user.Manager != nil {
		redactedUser := user.Manager.Redact().(User)
		redacted.Manager = &redactedUser
	}
	if user.Biography != nil {
		redactedBiography := user.Biography.Redact().(Biography)
		redacted.Biography = &redactedBiography
	}
	return redacted
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (user User) MarshalJSON() ([]byte, error) {
	type surrogate User
	data, err := json.Marshal(&struct {
		surrogate
		SelfURI URI `json:"selfUri"`
	}{
		surrogate: surrogate(user),
		SelfURI:   user.GetURI(),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}
