package gcloudcx

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// Integration represents a GCloudCX integration.
type Integration struct {
	ID                        uuid.UUID                          `json:"id"`
	Type                      IntegrationType                    `json:"integrationType"`
	Name                      string                             `json:"name"`
	Description               string                             `json:"description,omitempty"`
	Config                    IntegrationConfigurationInfo       `json:"config"`
	ConfigPropertiesSchemaURI *url.URL                           `json:"configPropertiesSchemaUri,omitempty"`
	ConfigAdvancedSchemaURI   *url.URL                           `json:"configAdvancedSchemaUri,omitempty"`
	HelpURI                   *url.URL                           `json:"helpUri,omitempty"`
	TermsOfServiceURI         *url.URL                           `json:"termsOfServiceUri,omitempty"`
	VendorName                string                             `json:"vendorName"`
	VendorWebsiteURI          *url.URL                           `json:"vendorWebsiteUri,omitempty"`
	MarketplaceURI            *url.URL                           `json:"marketplaceUri,omitempty"`
	FAQURI                    *url.URL                           `json:"faqUri,omitempty"`
	PrivacyPolicyURI          *url.URL                           `json:"privacyPolicyUri,omitempty"`
	SupportContactURI         *url.URL                           `json:"supportContactUri,omitempty"`
	SalesContactURI           *url.URL                           `json:"salesContactUri,omitempty"`
	HelpLinks                 []HelpLink                         `json:"helpLinks,omitempty"`
	Notes                     string                             `json:"notes,omitempty"`
	IntendedState             string                             `json:"intendedState"`
	NonInstallable            bool                               `json:"nonInstallable"`
	MaxInstances              int                                `json:"maxInstances"`
	Images                    []Image                            `json:"images"`
	Attributes                map[string]interface{}             `json:"attributes"`
	Credentials               map[string]CredentialSpecification `json:"credentials,omitempty"`
	UserPermissions           []string                           `json:"userPermissions,omitempty"`
	VendorOAUTHClientIDs      []string                           `json:"vendorOauthClientIds,omitempty"`
	Client                    *Client                            `json:"-"`
	logger                    *logger.Logger                     `json:"-"`
}

// CredentialSpecification represents a specification for a credential.
type CredentialSpecification struct {
	Title    string   `json:"title"`
	Required bool     `json:"required"`
	Types    []string `json:"credentialTypes"`
}

// Initialize initializes the object
//
// accepted parameters: *gcloudcx.Client, *logger.Logger
//
// implements Initializable
func (integration *Integration) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case uuid.UUID:
			integration.ID = parameter
		case *Client:
			integration.Client = parameter
		case *logger.Logger:
			integration.logger = parameter.Child("integration", "integration", "id", integration.ID)
		}
	}
	if integration.logger == nil {
		integration.logger = logger.Create("gcloudcx", &logger.NilStream{})
	}
}

// GetID gets the identifier of this
//
//	implements Identifiable
func (integration Integration) GetID() uuid.UUID {
	return integration.ID
}

// GetURI gets the URI of this
//
//	implements Addressable
func (integration Integration) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/integrations/%s", ids[0])
	}
	if integration.ID != uuid.Nil {
		return NewURI("/api/v2/integrations/%s", integration.ID)
	}
	return URI("/api/v2/integrations/")
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (integration Integration) String() string {
	if len(integration.Name) > 0 {
		return integration.Name
	}
	return integration.ID.String()
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (integration Integration) MarshalJSON() ([]byte, error) {
	type surrogate Integration
	data, err := json.Marshal(struct {
		surrogate
		ID                        core.UUID `json:"id"`
		ConfigPropertiesSchemaURI *core.URL `json:"configPropertiesSchemaUri,omitempty"`
		ConfigAdvancedSchemaURI   *core.URL `json:"configAdvancedSchemaUri,omitempty"`
		HelpURI                   *core.URL `json:"helpUri,omitempty"`
		TermsOfServiceURI         *core.URL `json:"termsOfServiceUri,omitempty"`
		VendorWebsiteURI          *core.URL `json:"vendorWebsiteUri,omitempty"`
		MarketplaceURI            *core.URL `json:"marketplaceUri,omitempty"`
		FAQURI                    *core.URL `json:"faqUri,omitempty"`
		PrivacyPolicyURI          *core.URL `json:"privacyPolicyUri,omitempty"`
		SupportContactURI         *core.URL `json:"supportContactUri,omitempty"`
		SalesContactURI           *core.URL `json:"salesContactUri,omitempty"`
	}{
		surrogate:                 surrogate(integration),
		ID:                        core.UUID(integration.ID),
		ConfigPropertiesSchemaURI: (*core.URL)(integration.ConfigPropertiesSchemaURI),
		ConfigAdvancedSchemaURI:   (*core.URL)(integration.ConfigAdvancedSchemaURI),
		HelpURI:                   (*core.URL)(integration.HelpURI),
		TermsOfServiceURI:         (*core.URL)(integration.TermsOfServiceURI),
		VendorWebsiteURI:          (*core.URL)(integration.VendorWebsiteURI),
		MarketplaceURI:            (*core.URL)(integration.MarketplaceURI),
		FAQURI:                    (*core.URL)(integration.FAQURI),
		PrivacyPolicyURI:          (*core.URL)(integration.PrivacyPolicyURI),
		SupportContactURI:         (*core.URL)(integration.SupportContactURI),
		SalesContactURI:           (*core.URL)(integration.SalesContactURI),
	})
	if err != nil {
		return nil, errors.JSONMarshalError.Wrap(err)
	}
	return data, nil
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (integration *Integration) UnmarshalJSON(payload []byte) error {
	type surrogate Integration
	var inner struct {
		surrogate
		ID                        core.UUID `json:"id"`
		ConfigPropertiesSchemaURI *core.URL `json:"configPropertiesSchemaUri,omitempty"`
		ConfigAdvancedSchemaURI   *core.URL `json:"configAdvancedSchemaUri,omitempty"`
		HelpURI                   *core.URL `json:"helpUri,omitempty"`
		TermsOfServiceURI         *core.URL `json:"termsOfServiceUri,omitempty"`
		VendorWebsiteURI          *core.URL `json:"vendorWebsiteUri,omitempty"`
		MarketplaceURI            *core.URL `json:"marketplaceUri,omitempty"`
		FAQURI                    *core.URL `json:"faqUri,omitempty"`
		PrivacyPolicyURI          *core.URL `json:"privacyPolicyUri,omitempty"`
		SupportContactURI         *core.URL `json:"supportContactUri,omitempty"`
		SalesContactURI           *core.URL `json:"salesContactUri,omitempty"`
	}
	if err := json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	*integration = Integration(inner.surrogate)
	integration.ID = uuid.UUID(inner.ID)
	integration.ConfigPropertiesSchemaURI = (*url.URL)(inner.ConfigPropertiesSchemaURI)
	integration.ConfigAdvancedSchemaURI = (*url.URL)(inner.ConfigAdvancedSchemaURI)
	integration.HelpURI = (*url.URL)(inner.HelpURI)
	integration.TermsOfServiceURI = (*url.URL)(inner.TermsOfServiceURI)
	integration.VendorWebsiteURI = (*url.URL)(inner.VendorWebsiteURI)
	integration.MarketplaceURI = (*url.URL)(inner.MarketplaceURI)
	integration.FAQURI = (*url.URL)(inner.FAQURI)
	integration.PrivacyPolicyURI = (*url.URL)(inner.PrivacyPolicyURI)
	integration.SupportContactURI = (*url.URL)(inner.SupportContactURI)
	integration.SalesContactURI = (*url.URL)(inner.SalesContactURI)
	if integration.ID == uuid.Nil {
		integration.ID = uuid.New()
	}
	return nil
}
