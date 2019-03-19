package purecloud

import (
	"net/url"
	"time"

	"github.com/gildas/go-logger"
)

// Client is the primary object to use PureCloud
type Client struct {
	Region         string
	Organization   *Organization
	DeploymentID   string
	API            *url.URL
	Proxy          *url.URL
	Authorization  *Authorization
	Logger         *logger.Logger
}

// ClientOptions contains the options to create a new Client
type ClientOptions struct {
	Region         string
	OrganizationID string
	DeploymentID   string
	Logger         *logger.Logger
}

// Authorization contains the login options to connect the client to PureCloud
type Authorization struct {
	GrantType    string
	ClientID     string
	Secret       string
	TokenType    string
	Token        string
	TokenExpires time.Time
}
