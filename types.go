package purecloud

import (
	"net/url"

	"github.com/gildas/go-logger"
)

// Client is the primary object to use PureCloud
type Client struct {
	Region              string             `json:"region"`
	Organization        *Organization      `json:"organization,omitempty"`
	DeploymentID        string             `json:"deploymentId"`
	API                 *url.URL           `json:"apiUrl,omitempty"`
	Proxy               *url.URL           `json:"proxyUrl,omitempty"`
	AuthorizationGrant  AuthorizationGrant `json:"auth"`
	Logger              *logger.Logger     `json:"-"`
}