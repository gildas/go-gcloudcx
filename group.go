package purecloud

import (
	"time"
	"github.com/gildas/go-logger"
)

// Group describe a Group of users
type Group struct {
	ID           string          `json:"id"`
	SelfURI      string          `json:"selfUri"`
	Name         string          `json:"name"`
	Type         string          `json:"type"`
	Description  string          `json:"description"`
	State        string          `json:"state"`
	MemberCount  int             `json:"memberCount"`
	Owners       []*User         `json:"owners"`
	Images       []*UserImage    `json:"images"`
	Addresses    []*Contact      `json:"addresses"`
	RulesVisible bool            `json:"rulesVisible"`
	Visibility   bool            `json:"visibility"`
	DateModified time.Time       `json:"dateModified"`
	Version      int             `json:"version"`
	Client       *Client         `json:"-"`
	Logger       *logger.Logger  `json:"-"`
}

// Initialize initializes this from the given Client
//   implements Initializable
//   if the user ID is not given, /users/me is fetched (if grant allows)
func (group *Group) Initialize(parameters ...interface{}) error {
	client, logger, err := ExtractClientAndLogger(parameters...)
	if err != nil {
		return err
	}
	if len(group.ID) > 0 {
		if err := client.Get("/groups/" + group.ID, &group); err != nil {
			return err
		}
	}
	group.Client = client
	group.Logger = logger.Topic("group").Scope("group").Record("group", group.ID)
	return nil
}

// GetID gets the identifier of this
//   implements Identifiable
func (group Group) GetID() string {
	return group.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
func (group Group) String() string {
	if len(group.Name) > 0 {
		return group.Name
	}
	return group.ID
}