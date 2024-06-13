package gcloudcx_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/gildas/go-gcloudcx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanMarshallAccessToken(t *testing.T) {
	token := gcloudcx.NewAccessToken("Very Long String", time.Date(1996, 9, 23, 0, 0, 0, 0, time.UTC))
	expected := fmt.Sprintf(`{"id": "%s", "tokenType": "Bearer", "token": "Very Long String", "expiresOn": "1996-09-23T00:00:00Z"}`, token.ID)
	data, err := json.Marshal(token)
	require.Nil(t, err, "Failed to marshall token")
	require.NotEmpty(t, data, "Failed to marshall token")
	assert.JSONEq(t, expected, string(data))
}

func TestCanUnmarshallAccessToken(t *testing.T) {
	source := `{"tokenType": "Bearer", "token": "Very Long String", "expiresOn": "1996-09-23T00:00:00Z"}`
	token := gcloudcx.AccessToken{}
	err := json.Unmarshal([]byte(source), &token)
	require.Nil(t, err, "Failed to unmarshall token")
	assert.NotNil(t, token.ID, "Token ID should not be nil")
	assert.Equal(t, "Bearer", token.Type)
	assert.Equal(t, "Very Long String", token.Token)
	assert.Equal(t, time.Date(1996, 9, 23, 0, 0, 0, 0, time.UTC), token.ExpiresOn)
	assert.True(t, token.IsExpired(), "Token should be expired!")
}

func TestCanTellExpirationOfAccessToken(t *testing.T) {
	token := gcloudcx.NewAccessTokenWithDurationAndType("Very Long String", "Bearer", 2*time.Hour)

	assert.False(t, token.IsExpired(), "Token should not be expired")
	assert.True(t, 1*time.Hour < token.ExpiresIn(), "Token should expire in an hour at least")

	token = gcloudcx.NewAccessTokenWithType("Bearer", "Very Long String", time.Date(1996, 9, 23, 0, 0, 0, 0, time.UTC))
	assert.True(t, token.IsExpired(), "Token should be expired")
	assert.True(t, time.Duration(0) == token.ExpiresIn(), "Token should expire in 0")
}

func TestCanResetAccessToken(t *testing.T) {
	token := gcloudcx.NewAccessTokenWithDuration("Very Long String", 2*time.Hour)
	token.Reset()
	assert.Empty(t, token.Token, "The Token string should be empty")
	assert.Empty(t, token.Type, "The Token type should be empty")
	assert.True(t, token.IsExpired(), "The Token should be expired")
	assert.False(t, token.IsValid(), "The Token should not be valid")
}

func TestCanResetGrantAccessToken(t *testing.T) {
	token := gcloudcx.NewAccessTokenWithDuration("Very Long String", 2*time.Hour)
	client := gcloudcx.NewClient(&gcloudcx.ClientOptions{}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{Token: *token})
	assert.Equal(t, "Bearer", client.Grant.AccessToken().Type)
	assert.Equal(t, "Very Long String", client.Grant.AccessToken().Token)

	client.Grant.AccessToken().Reset()
	assert.Empty(t, client.Grant.AccessToken().Token, "The Token string should be empty")
	assert.Empty(t, client.Grant.AccessToken().Type, "The Token type should be empty")
	assert.True(t, client.Grant.AccessToken().IsExpired(), "The Token should be expired")
	assert.False(t, client.Grant.AccessToken().IsValid(), "The Token should not be valid")
}
