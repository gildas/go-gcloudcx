package purecloud

// AuthorizationGrant describes the capabilities authorization grants must have
type AuthorizationGrant interface {
	// Authorize this Grant with PureCloud
	Authorize(client *Client) error
	AccessToken() *AccessToken
}
