package purecloud

// Contact describes something that can be contacted
type Contact struct {
	Type      string `json:"type"`      // PRIMARY, WORK, WORK2, WORK3, WORK4, HOME, MOBILE, MAIN
	MediaType string `json:"mediaType"` // PHONE, EMAIL, SMS
	Display   string `json:"display,omitempty"`
	Address   string `json:"address,omitempty"`   // If present, there is no Extension
	Extension string `json:"extension,omitempty"` // If present, there is no Address
}