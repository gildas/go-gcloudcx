package purecloud

// Logout logs out a Client from PureCloud
func (client *Client) Logout() (err error) {
	if err := client.Delete("/tokens/me", nil); err != nil { return err }
	if client.AuthorizationGrant != nil {
		client.AuthorizationGrant.AccessToken().Reset()
	}
	return nil
}