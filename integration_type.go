package gcloudcx

// IntegrationType represents the type of integration.
type IntegrationType struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Provider    string  `json:"provider"`
	Category    string  `json:"category"`
	Images      []Image `json:"images"`
}

const (
	IntegrationTypeGenesysDigitalBotConnector = "genesys-digital-bot-connector"
)
