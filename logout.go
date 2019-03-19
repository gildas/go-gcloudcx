package purecloud

// Logout logs out a Client from PureCloud
func (client *Client) Logout() (err error) {
	if err := client.Delete("tokens/me", nil, nil); err != nil { return err }
	client.Authorization.Token = ""
	return nil
}