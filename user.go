package purecloud

import "fmt"

// User describe a PureCloud User
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

func (client *Client) GetMyUser() (*User, error) {
	user := &User{}
	if err := client.Get("/users/me", &user); err != nil {
		return nil, err
	}
	return user, nil
}

func (user User) ActivityID() string {
	return fmt.Sprintf("v2.users.%s.conversations.chats", user.ID)
}
