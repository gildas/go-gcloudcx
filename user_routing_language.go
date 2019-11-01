package purecloud

// UserRoutingLanguage describe a Routing Language for a User
type UserRoutingLanguage struct {
	ID           string  `json:"id"`
	SelfURI      string  `json:"selfUri"`
	Name         string  `json:"name"`
	LanguageURI  string  `json:"languageUri"`
	State        string  `json:"state"`
	Proficiency  float64 `json:"proficiency"`
}