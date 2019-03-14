package purecloud

// Logout logs out a Client from PureCloud
func (client *Client) Logout() (err error) {
	if err := client.delete("tokens/me", nil, nil); err != nil { return err }
	client.Token = Token{}
	return nil
}