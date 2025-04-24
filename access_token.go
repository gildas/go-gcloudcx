package gcloudcx

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
)

// AccessToken is used to consume the GCloud API
//
// It must be obtained via an AuthorizationGrant
type AccessToken struct {
	ID               uuid.UUID `json:"id" db:"key"`
	Type             string    `json:"tokenType"`
	Token            string    `json:"token"`
	ExpiresOn        time.Time `json:"expiresOn"` // UTC!
	AuthorizedScopes []string  `json:"authorizedScopes"`
}

// UpdatedAccessToken describes an updated Access Token
//
// This object is sent by the AuthorizationGrant.Authorize() when the token is updated
type UpdatedAccessToken struct {
	AccessToken
	CustomData interface{}
}

var (
	hashKey      = []byte(core.GetEnvAsString("PURECLOUD_SESSION_HASH_KEY", "Pur3Cl0udS3ss10nH@5hK3y"))
	blockKey     = []byte(core.GetEnvAsString("PURECLOUD_SESSION_BLOCK_KEY", "Pur3Cl0udS3ss10nBl0ckK3y"))
	secureCookie = securecookie.New(hashKey, blockKey)
)

// NewAccessToken creates a new AccessToken
func NewAccessToken(token string, expiresOn time.Time) *AccessToken {
	return &AccessToken{
		ID:        uuid.New(),
		Type:      "Bearer",
		Token:     token,
		ExpiresOn: expiresOn,
	}
}

// NewAccessTokenWithType creates a new AccessToken with a type
func NewAccessTokenWithType(tokenType, token string, expiresOn time.Time) *AccessToken {
	return &AccessToken{
		ID:        uuid.New(),
		Type:      tokenType,
		Token:     token,
		ExpiresOn: expiresOn,
	}
}

// NewAccessTokenWithDuration creates a new AccessToken that expires in a given duration
func NewAccessTokenWithDuration(token string, expiresIn time.Duration) *AccessToken {
	return &AccessToken{
		ID:        uuid.New(),
		Type:      "Bearer",
		Token:     token,
		ExpiresOn: time.Now().UTC().Add(expiresIn),
	}
}

// NewAccessTokenWithDurationAndType creates a new AccessToken with a type and that expires in a given duration
func NewAccessTokenWithDurationAndType(tokenType, token string, expiresIn time.Duration) *AccessToken {
	return &AccessToken{
		ID:        uuid.New(),
		Type:      tokenType,
		Token:     token,
		ExpiresOn: time.Now().UTC().Add(expiresIn),
	}
}

// Reset resets the Token so it is expired and empty
func (token *AccessToken) Reset() {
	token.Type = ""
	token.Token = ""
	token.ExpiresOn = time.Time{}
}

// LoadFromCookie loads this token from a cookie in the given HTTP Request
func (token *AccessToken) LoadFromCookie(r *http.Request, cookieName string) *AccessToken {
	if cookie, err := r.Cookie(cookieName); err == nil {
		var jsonToken string

		if err = secureCookie.Decode(cookieName, cookie.Value, &jsonToken); err == nil {
			_ = json.Unmarshal([]byte(jsonToken), token)
		}
	}
	return token
}

// SaveToCookie saves this token to a cookie in the given HTTP ResponseWriter
func (token AccessToken) SaveToCookie(w http.ResponseWriter, cookieName string) {
	jsonToken, _ := json.Marshal(token)
	encodedID, _ := secureCookie.Encode("pcsession", string(jsonToken))
	http.SetCookie(w, &http.Cookie{Name: "pcsession", Value: encodedID, Path: "/", HttpOnly: true, Secure: true})
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

// Redact redacts sensitive information
//
// implements logger.Redactable
func (token AccessToken) Redact() any {
	redacted := token
	if len(redacted.Token) > 0 {
		redacted.Token = logger.RedactWithHash(token.Token)
	}
	return redacted
}

// String gets a string representation of this AccessToken
func (token AccessToken) String() string {
	return token.Type + " " + token.Token
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (token AccessToken) MarshalJSON() ([]byte, error) {
	type surrogate AccessToken

	data, err := json.Marshal(struct {
		ID core.UUID `json:"id"`
		surrogate
		ExpiresOn core.Time `json:"expiresOn"`
	}{
		ID:        core.UUID(token.ID),
		surrogate: surrogate(token),
		ExpiresOn: core.Time(token.ExpiresOn),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON decodes JSON
//
// implements json.Unmarshaler
func (token *AccessToken) UnmarshalJSON(payload []byte) (err error) {
	type surrogate AccessToken

	var inner struct {
		ID core.UUID `json:"id"`
		surrogate
		ExpiresOn core.Time `json:"expiresOn"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*token = AccessToken(inner.surrogate)
	token.ID = uuid.UUID(inner.ID)
	token.ExpiresOn = inner.ExpiresOn.AsTime()
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}
	return nil
}
