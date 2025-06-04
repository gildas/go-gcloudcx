package gcloudcx

// IntegrationCredentialSpecification represents information about an integration credential type.
type IntegrationCredentialSpecification struct {
	Title           string   `json:"title"`
	Required        bool     `json:"required,omitempty"`
	CredentialTypes []string `json:"credentialTypes,omitempty"`
}
