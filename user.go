package purecloud

import (
	"strings"
	"net/url"
	"fmt"
)

// User describes a PureCloud User
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Division struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		SelfURI string `json:"selfUri"`
	} `json:"division"`
	Chat struct {
		JabberID string `json:"jabberId"`
	} `json:"chat"`
	Department         string `json:"department"`
	Mail               string `json:"email"`
	PrimaryContactInfo []struct {
		Address   string `json:"address"`
		MediaType string `json:"mediaType"`
		Type      string `json:"type"`
	} `json:"primaryContactInfo"`
	Addresses []struct {
		Display   string `json:"display"`
		MediaType string `json:"mediaType"`
		Type      string `json:"type"`
	} `json:"addresses"`
	State         string `json:"state"`
	Title         string `json:"title"`
	UserName      string `json:"username"`
	Version       int    `json:"version"`
	AcdAutoAnswer bool   `json:"acdAutoAnswer"`
	SelfURI       string `json:"selfUri"`
}

type UserImage struct {
	ImageURL   *url.URL `json:"imageUri"`
	Resolution string   `json:"resolution"`
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
	if err := client.Get("/users/me?" + query.Encode(), &user); err != nil {
		return nil, err
	}
	user.Client = client
	return user, nil
}

func (user User) ActivityID() string {
	return fmt.Sprintf("v2.users.%s.conversations.chats", user.ID)
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