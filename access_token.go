package purecloud

import (
	"time"
)

// AccessToken is used to consume the PureCloud API
//
// It must be obtained via an AuthorizationGrant
type AccessToken struct {
	Type      string    `json:"tokenType"`
	Token     string    `json:"token"`
	ExpiresOn time.Time `json:"tokenExpires"`  // UTC!
}

// Reset resets the Token so it is expired and empty
func (token *AccessToken) Reset() {
	token.Type      = ""
	token.Token     = ""
	token.ExpiresOn = time.Time{}
}

// IsValid tells if this AccessToken is valid
func (token AccessToken) IsValid() bool {
	return len(token.Token) > 0 // TODO: We should have && !token.IsExpired()
}

// IsExpired tells if this AccessToken is expired or not
func (token AccessToken) IsExpired() bool {
	return time.Now().UTC().After(token.ExpiresOn)
}

// ExpiresIn tells when the token should expire
func (token AccessToken) ExpiresIn() time.Duration {
	if token.IsExpired() {
		return time.Duration(0)
	}
	return token.ExpiresOn.Sub(time.Now().UTC())
}

func (token AccessToken) String() string {
	return token.Type + " " + token.Token
}