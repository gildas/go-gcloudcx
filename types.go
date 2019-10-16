package purecloud

import (
	"net/url"

	"github.com/gildas/go-logger"
)

// Client is the primary object to use PureCloud
type Client struct {
	Region              string             `json:"region"`
	DeploymentID        string             `json:"deploymentId"`
	API                 *url.URL           `json:"apiUrl,omitempty"`
	Proxy               *url.URL           `json:"proxyUrl,omitempty"`
	AuthorizationGrant  AuthorizationGrant `json:"auth"`
	Logger              *logger.Logger     `json:"-"`
}

// DomainEntityRef describes a DomainEntity Reference
type DomainEntityRef struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SelfURI  string `json:"self_uri"`
}

// Identifiable describes things that carry an identifier
type Identifiable interface {
	// GetID gets the identifier of this
	GetID() string
}