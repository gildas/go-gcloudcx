package purecloud

import (
	"net/url"
	"time"

	"github.com/gildas/go-logger"
)

// Client is the primary object to use PureCloud
type Client struct {
	Region         string         `json:"region"`
	Organization   *Organization  `json:"organization,omitempty"`
	DeploymentID   string         `json:"deploymentId"`
	API            *url.URL       `json:"apiUrl,omitempty"`
	Proxy          *url.URL       `json:"proxyUrl,omitempty"`
	Authorization  *Authorization `json:"auth"`
	Logger         *logger.Logger `json:"-"`
}

// ClientOptions contains the options to create a new Client
type ClientOptions struct {
	Region         string
	OrganizationID string
	DeploymentID   string
	ClientID       string
	ClientSecret   string
	Proxy          *url.URL
	Logger         *logger.Logger
}

// Authorization contains the login options to connect the client to PureCloud
type Authorization struct {
	GrantType    string    `json:"grantType"`
	ClientID     string    `json:"clientId"`
	Secret       string    `json:"clientSecret"`
	TokenType    string    `json:"tokenType"`
	Token        string    `json:"token"`
	TokenExpires time.Time `json:"tokenExpires"`
}
