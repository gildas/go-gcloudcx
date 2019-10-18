package purecloud

import (
	"net/url"
	"strings"
)

// User describes a PureCloud User
type User struct {
	ID                  string                   `json:"id"`
	Name                string                   `json:"name"`
	UserName            string                   `json:"username"`
	Department          string                   `json:"department"`
	Title               string                   `json:"title"`
	Division            Division                 `json:"division"`
	Mail                string                   `json:"email"`
	Images              []UserImage              `json:"images"`
	PrimaryContact      []Contact                `json:"primaryContactInfo"`
	Addresses           []Contact                `json:"addresses"`
	State               string                   `json:"state"`
	RoutingStatus       *RoutingStatus           `json:"routingStatus,omitempty"`
	Presence            *UserPresence            `json:"presence,omitempty"`
	AcdAutoAnswer       bool                     `json:"acdAutoAnswer"`
	Employer            *EmployerInfo            `json:"employerInfo,omitempty"`
	Manager             *User                    `json:"manager,omitempty"`
	ConversationSummary *UserConversationSummary `json:"conversationSummary,omitempty"`
	GeoLocation         *GeoLocation             `json:"geolocation"`
	Chat                struct {
		JabberID string `json:"jabberId"`
	} `json:"chat"`
	SelfURI string  `json:"selfUri"`
	Version int     `json:"version"`
	Client  *Client `json:"-"`
	// TODO: Continue to add objects...
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
	if err := client.Get("/users/me?"+query.Encode(), &user); err != nil {
		return nil, err
	}
	user.Client = client
	return user, nil
}

// GetID gets the identifier of this
//   implements Identifiable
func (user User) GetID() string {
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
	return user.ID
}