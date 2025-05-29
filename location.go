package gcloudcx

// LocationDefinition describes a location (office, etc)
type LocationDefinition struct {
	ID              string                   `json:"id"`
	SelfURI         URI                      `json:"selfUri"`
	Name            string                   `json:"name"`
	ContactUser     *AddressableEntityRef    `json:"contactUser"`
	EmergencyNumber *LocationEmergencyNumber `json:"emergencyNumber"`
	Address         *LocationAddress         `json:"address"`
	AddressVerified bool                     `json:"addressVerified"`
	State           string                   `json:"state"`
	Notes           string                   `json:"notes"`
	Path            []string                 `json:"path"`
	ProfileImage    []Image                  `json:"profileImage"`
	FloorplanImage  []Image                  `json:"flooreImage"`
	Version         int                      `json:"version"`
}

// LocationEmergencyNumber describes a Location's Emergency Number
type LocationEmergencyNumber struct {
	Type   string `json:"type"` // default, elin
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

// GeoLocation describes a location with coordinates
type GeoLocation struct {
	ID        string                `json:"id"`
	SelfURI   URI                   `json:"selfUri"`
	Name      string                `json:"name"`
	Locations []*LocationDefinition `json:"locations"`
}

// Location describes a location in a building
type Location struct {
	ID                 string              `json:"id"`
	FloorplanID        string              `json:"floorplanId"`
	Coordinates        map[string]float64  `json:"coordinates"`
	Notes              string              `json:"notes"`
	LocationDefinition *LocationDefinition `json:"locationDefinition"`
}

// GetID gets the identifier of this
//
//	implements Identifiable
func (location LocationDefinition) GetID() string {
	return location.ID
}

// GetURI gets the URI of this
//
//	implements Addressable
func (location LocationDefinition) GetURI() URI {
	return location.SelfURI
}

// String gets a string version
//
//	implements the fmt.Stringer interface
func (location LocationDefinition) String() string {
	if len(location.Name) != 0 {
		return location.Name
	}
	return location.ID
}
