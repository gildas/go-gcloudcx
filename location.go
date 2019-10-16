package purecloud

// LocationDefinition describes a location (office, etc)
type LocationDefinition struct {
	ID              string                  `json:"id"`
	Name            string                  `json:"name"`
	ContactUser     AddressableEntityRef    `json:"contactUser"`
	EmergencyNumber LocationEmergencyNumber `json:"emergencyNumber"`
	Address         LocationAddress         `json:"address"`
	AddressVerified bool                    `json:"addressVerified"`
	State           string                  `json:"state"`
	Notes           string                  `json:"notes"`
	Path            []string                `json:"path"`
	ProfileImage    []LocationImage         `json:"profileImage"`
	FloorplanImage  []LocationImage         `json:"flooreImage"`
	Version         int                     `json:"version"`
	SelfURI         string                  `json:"selfUri"`
}

// LocationEmergencyNumber describes a Location's Emergency Number
type LocationEmergencyNumber struct {
	Type string   `json:"type"` // default, elin
	Number string `json:"number"`
	E164   string `json:"e164"`
}

// LocationAddress describes the address of a Location
type LocationAddress struct {
	Country     string `json:"country"`
	CountryName string `json:"countryName"`
	State       string `json:"State"`
	City        string `json:"City"`
	ZipCode     string `json:"zipcode"`
	Street1     string `json:"street1"`
	Street2     string `json:"street2"`
}

type LocationImage struct {
	ImageURL   string `json:"imageUrl"`
	Resolution string `json:"resolution"`
}

// GeoLocation describes a location with coordinates
type GeoLocation struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Locations  []LocationDefinition `json:"locations"`

	SelfURI        string         `json:"selfUri"`
}