package purecloud

// Logout logs out a Client from PureCloud
func (client *Client) Logout() {
	client.Delete("/tokens/me", nil) // we don't care much about the error as we are logging out
	if client.AuthorizationGrant != nil {
		client.AuthorizationGrant.AccessToken().Reset()
	}
}