package purecloud

// User describe a PureCloud User
type User struct {
	ID                         string          `json:"id"`
}

func (client *Client) GetMyUser() (*User, error) {
	user := &User{}
	if err := client.Get("/users/me", &user); err != nil {
		return nil, err
	}
	return user, nil
}