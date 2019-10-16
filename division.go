package purecloud

// Division describes an Authorization Division
type Division struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	SelfURI string  `json:"selfUri"`
	Client  *Client `json:"-"`
}