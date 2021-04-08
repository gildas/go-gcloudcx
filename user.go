package purecloud

import (
	"net/url"
	"strings"

	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// User describes a PureCloud User
type User struct {
	ID                  uuid.UUID                `json:"id"`
	SelfURI             string                   `json:"selfUri"`
	Name                string                   `json:"name"`
	UserName            string                   `json:"username"`
	Department          string                   `json:"department"`
	Title               string                   `json:"title"`
	Division            *Division                `json:"division"`
	Mail                string                   `json:"email"`
	Images              []UserImage              `json:"images"`
	PrimaryContact      []*Contact               `json:"primaryContactInfo"`
	Addresses           []*Contact               `json:"addresses"`
	State               string                   `json:"state"`
	Presence            *UserPresence            `json:"presence,omitempty"`
	OutOfOffice         *OutOfOffice             `json:"outOfOffice"`
	AcdAutoAnswer       bool                     `json:"acdAutoAnswer"`
	RoutingStatus       *RoutingStatus           `json:"routingStatus,omitempty"`
	ProfileSkills       []string                 `json:"profileSkills"`
	Skills              []*UserRoutingSkill      `json:"skills"`
	Languages           []*UserRoutingLanguage   `json:"languages"`
	LanguagePreference  string                   `json:"languagePreference"`
	Groups              []*Group                 `json:"groups"`
	Station             *UserStations            `json:"station"`
	Authorization       *UserAuthorization       `json:"authorization"`
	Employer            *EmployerInfo            `json:"employerInfo,omitempty"`
	Manager             *User                    `json:"manager,omitempty"`
	Certifications      []string                 `json:"certifications"`
	Biography           *Biography               `json:"biography"`
	ConversationSummary *UserConversationSummary `json:"conversationSummary,omitempty"`
	Locations           []*Location              `json:"locations"`
	GeoLocation         *GeoLocation             `json:"geolocation"`
	Chat                struct {
		JabberID string `json:"jabberId"`
	} `json:"chat"`
	Version int            `json:"version"`
	Client  *Client        `json:"-"`
	Logger  *logger.Logger `json:"-"`
}

// Initialize initializes this from the given Client
//   implements Initializable
//   if the user ID is not given, /users/me is fetched (if grant allows)
func (user *User) Initialize(parameters ...interface{}) error {
	client, logger, id, err := parseParameters(parameters...)
	if err != nil {
		return err
	}
	if id != uuid.Nil {
		if err := client.Get(NewURI("/users/%s", user.ID), &user); err != nil {
			return err
		}
	} else if _, ok := client.AuthorizationGrant.(*ClientCredentialsGrant); !ok { // /users/me is not possible with ClientCredentialsGrant
		if err := client.Get("/users/me", &user); err != nil {
			return err
		}
	}
	user.Client = client
	user.Logger = logger.Child("user", "user", "user", user.ID)
	return nil
}

// GetMyUser retrieves the User that authenticated with the client
//   properties is one of more properties that should be expanded
//   see https://developer.mypurecloud.com/api/rest/v2/users/#get-api-v2-users-me
func (client *Client) GetMyUser(properties ...string) (*User, error) {
	query := url.Values{}
	if len(properties) > 0 {
		query.Add("expand", strings.Join(properties, ","))
	}
	user := &User{}
	if err := client.Get(NewURI("/users/me?%s", query.Encode()), &user); err != nil {
		return nil, err
	}
	user.Client = client
	return user, nil
}

// GetID gets the identifier of this
//   implements Identifiable
func (user User) GetID() uuid.UUID {
	return user.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
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
