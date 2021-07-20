package gcloudcx

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// Client is the primary object to use Gcloud
type Client struct {
	Region         string         `json:"region"`
	DeploymentID   uuid.UUID      `json:"deploymentId"`
	Organization   *Organization  `json:"-"`
	API            *url.URL       `json:"apiUrl,omitempty"`
	LoginURL       *url.URL       `json:"loginUrl,omitempty"`
	Proxy          *url.URL       `json:"proxyUrl,omitempty"`
	Grant          Authorizer     `json:"-"`
	RequestTimeout time.Duration  `json:"requestTimout"`
	Logger         *logger.Logger `json:"-"`
}

// ClientOptions contains the options to create a new Client
type ClientOptions struct {
	Region         string
	OrganizationID uuid.UUID
	DeploymentID   uuid.UUID
	Proxy          *url.URL
	Grant          Authorizer
	RequestTimeout time.Duration
	Logger         *logger.Logger
}

// NewClient creates a new Gcloud Client
func NewClient(options *ClientOptions) *Client {
	if options == nil {
		options = &ClientOptions{}
	}
	if len(options.Region) == 0 {
		options.Region = "mypurecloud.com"
	}
	if options.RequestTimeout < 2*time.Second {
		options.RequestTimeout = 10 * time.Second
	}
	client := Client{
		Proxy:          options.Proxy,
		DeploymentID:   options.DeploymentID,
		Organization:   &Organization{ID: options.OrganizationID},
		Grant:          options.Grant,
		RequestTimeout: options.RequestTimeout,
	}
	return client.SetLogger(options.Logger).SetRegion(options.Region)
}

// SetLogger sets the logger
func (client *Client) SetLogger(log *logger.Logger) *Client {
	client.Logger = logger.CreateIfNil(log, "gcloudcx").Child("gcloudcx", "gcloudcx")
	return client
}

// SetRegion sets the region and its main API
func (client *Client) SetRegion(region string) *Client {
	client.Region = region
	client.API, _ = url.Parse(fmt.Sprintf("https://api.%s", region))
	client.LoginURL, _ = url.Parse(fmt.Sprintf("https://login.%s", region))
	return client
}

// SetAuthorizationGrant sets the Authorization Grant
func (client *Client) SetAuthorizationGrant(grant Authorizer) *Client {
	client.Grant = grant
	return client
}

// IsAuthorized tells if the client has an Authorization Token
// It migt be expired and the app should login again as needed
func (client *Client) IsAuthorized() bool {
	return client.Grant.AccessToken().IsValid()
}

// Fetch fetches an initializable object
func (client *Client) Fetch(object Initializable) error {
	return object.Initialize(client)
}

func (client *Client) CheckPermissions(permissions ...string) (permitted []string, missing []string) {
	log := client.Logger.Child(nil, "checkpermissions")
	subject, err := client.FetchRolesAndPermissions()
	if err != nil {
		return []string{}, permissions
	}
	permitted = []string{}
	missing   = []string{}
	for _, desired := range permissions {
		elements := strings.Split(desired, ":")
		if len(elements) < 3 {
			log.Warnf("This permission is invalid: %s (%d elements)", desired, len(elements))
			missing = append(missing, desired)
			break
		}
		desiredDomain := elements[0]
		desiredEntity := elements[1]
		desiredAction := elements[2]
		found := false
		log.Tracef("Checking Domain: %s, Entity: %s, Action: %s", desiredDomain, desiredEntity, desiredAction)
		for _, grant := range subject.Grants {
			for _, policy := range grant.Role.Policies {
				if policy.Domain == desiredDomain && (policy.EntityName == "*" || desiredEntity == policy.EntityName) {
					for _, action := range policy.Actions {
						if action == "*" || action == desiredAction {
							log.Tracef("  OK: %s:%s:%s", policy.Domain, policy.EntityName, action)
							permitted = append(permitted, desired)
							found = true
							break
						}
					}
				}
				if found {
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			missing = append(missing, desired)
		}
	}
	return
}

// FetchRolesAndPermissions fetches roles and permissions for the current client
func (client *Client) FetchRolesAndPermissions() (*AuthorizationSubject, error) {
	return client.FetchRolesAndPermissionsOf(client.Grant)
}

// FetchRolesAndPermissions fetches roles and permissions for the current client
func (client *Client) FetchRolesAndPermissionsOf(id core.Identifiable) (*AuthorizationSubject, error) {
	log := client.Logger.Child(nil, "fetch_roles_permissions")
	subject := AuthorizationSubject{}

	log.Debugf("Fetching roles and permissions for %s", id.GetID())
	if err := client.Get(NewURI("/authorization/subjects/%s", id.GetID().String()), &subject); err != nil {
		return nil, err
	}
	return &subject, nil
}
