package purecloud

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gildas/go-core"
	"github.com/gorilla/securecookie"
)

// AccessToken is used to consume the PureCloud API
//
// It must be obtained via an AuthorizationGrant
type AccessToken struct {
	Type      string    `json:"tokenType"`
	Token     string    `json:"token"`
	ExpiresOn time.Time `json:"tokenExpires"`  // UTC!
}

var (
	hashKey      = []byte(core.GetEnvAsString("PURECLOUD_SESSION_HASH_KEY",  "Pur3Cl0udS3ss10nH@5hK3y"))
	blockKey     = []byte(core.GetEnvAsString("PURECLOUD_SESSION_BLOCK_KEY", "Pur3Cl0udS3ss10nBl0ckK3y"))
	secureCookie = securecookie.New(hashKey, blockKey)
)

// Reset resets the Token so it is expired and empty
func (token *AccessToken) Reset() {
	token.Type      = ""
	token.Token     = ""
	token.ExpiresOn = time.Time{}
}

// LoadFromCookie loads this token from a cookie in the given HTTP Request
func (token *AccessToken) LoadFromCookie(r *http.Request, cookieName string) (*AccessToken) {
	if cookie, err := r.Cookie(cookieName); err == nil {
		var jsonToken string

		if err = secureCookie.Decode(cookieName, cookie.Value, &jsonToken); err == nil {
			json.Unmarshal([]byte(jsonToken), token)
		}
	}
	return token
}

// SaveToCookie saves this token to a cookie in the given HTTP ResponseWriter
func (token AccessToken) SaveToCookie(w http.ResponseWriter, cookieName string) {
	jsonToken, _ := json.Marshal(token)
	encodedID, _ := secureCookie.Encode("pcsession", string(jsonToken))
	http.SetCookie(w, &http.Cookie{Name: "pcsession", Value: encodedID, Path: "/", HttpOnly: true})
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