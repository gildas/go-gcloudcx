package purecloud

import (
	"net/url"
	"time"

	logger "bitbucket.org/gildas_cherruel/go-logger"
)

// Client is the primary object to use PureCloud
type Client struct {
	Region         string
	OrganizationID string
	DeploymentID   string
	API            url.URL
	Token          Token
	Logger         *logger.Logger
}

// Token contains the Access Token after logging in PureCloud
type Token struct {
	Type    string
	Token   string
	Expires time.Time
}

// ClientOptions contains the options to create a new Client
type ClientOptions struct {
	Region         string
	OrganizationID string
	DeploymentID   string
	Logger         *logger.Logger
}

// LoginOptions contains the login options to connect the client to PureCloud
type LoginOptions struct {
	GrantType string
	ClientID  string
	Secret    string
}
