package purecloud

// Biography describes a User's biography
type Biography struct {
	Biography string   `json:"biography"`
	Interests []string `json:"interests"`
	Hobbies   []string `json:"hobbies"`
	Spouse    string   `json:"spouse"`
}