package purecloud_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gildas/go-purecloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanMarshallAccessToken(t *testing.T) {
	token := purecloud.AccessToken{
		Type:      "Bearer",
		Token:     "Very Long String",
		ExpiresOn: time.Date(1996, 9, 23, 0, 0, 0, 0, time.UTC),
	}

	expected := `{"tokenType": "Bearer", "token": "Very Long String", "tokenExpires": "1996-09-23T00:00:00Z"}`
	data, err := json.Marshal(token)
	require.Nil(t, err, "Failed to marshall token")
	require.NotEmpty(t, data, "Failed to marshall token")
	assert.JSONEq(t, expected, string(data))
}

func TestCanUnmarshallAccessToken(t *testing.T) {
	source := `{"tokenType": "Bearer", "token": "Very Long String", "tokenExpires": "1996-09-23T00:00:00Z"}`
	token := purecloud.AccessToken{}

	err := json.Unmarshal([]byte(source), &token)
	require.Nil(t, err, "Failed to unmarshall token")
	assert.Equal(t, "Bearer", token.Type)
	assert.Equal(t, "Very Long String", token.Token)
	assert.Equal(t, time.Date(1996, 9, 23, 0, 0, 0, 0, time.UTC), token.ExpiresOn)
	assert.True(t, token.IsExpired(), "Token should be expired!")
}

func TestCanTellExpirationOfAccessToken(t *testing.T) {
	token := purecloud.AccessToken{
		Type:      "Bearer",
		Token:     "Very Long String",
		ExpiresOn: time.Now().UTC().Add(2 * time.Hour),
	}

	assert.False(t, token.IsExpired(), "Token should not be expired")
	assert.True(t, 1*time.Hour < token.ExpiresIn(), "Token should expire in an hour at least")

	token.ExpiresOn = time.Now().UTC().AddDate(0, 0, -1)
	assert.True(t, token.IsExpired(), "Token should be expired")
	assert.True(t, time.Duration(0) == token.ExpiresIn(), "Token should expire in 0")
}

func TestCanResetAccessToken(t *testing.T) {
	token := purecloud.AccessToken{
		Type:      "Bearer",
		Token:     "Very Long String",
		ExpiresOn: time.Now().UTC().Add(2 * time.Hour),
	}

	token.Reset()
	assert.Empty(t, token.Token, "The Token string should be empty")
	assert.Empty(t, token.Type, "The Token type should be empty")
	assert.True(t, token.IsExpired(), "The Token should be expired")
	assert.False(t, token.IsValid(), "The Token should not be valid")
}

func TestCanResetGrantAccessToken(t *testing.T) {
	token := purecloud.AccessToken{
		Type:      "Bearer",
		Token:     "Very Long String",
		ExpiresOn: time.Now().UTC().Add(2 * time.Hour),
	}
	client := purecloud.NewClient(&purecloud.ClientOptions{}).SetAuthorizationGrant(&purecloud.ClientCredentialsGrant{Token: token})
	assert.Equal(t, "Bearer", client.AuthorizationGrant.AccessToken().Type)
	assert.Equal(t, "Very Long String", client.AuthorizationGrant.AccessToken().Token)

	client.AuthorizationGrant.AccessToken().Reset()
	assert.Empty(t, client.AuthorizationGrant.AccessToken().Token, "The Token string should be empty")
	assert.Empty(t, client.AuthorizationGrant.AccessToken().Type, "The Token type should be empty")
	assert.True(t, client.AuthorizationGrant.AccessToken().IsExpired(), "The Token should be expired")
	assert.False(t, client.AuthorizationGrant.AccessToken().IsValid(), "The Token should not be valid")
}
