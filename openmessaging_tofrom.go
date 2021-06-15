package purecloud

type OpenMessageTo struct {
	ID string `json:"id"`
}

type OpenMessageFrom struct {
	ID        string   `json:"id"`
	Type      string   `json:"idType"`
	Firstname string   `json:"firstName"`
	Lastname  string   `json:"lastName"`
	Nickname  string   `json:"nickname"`
}
