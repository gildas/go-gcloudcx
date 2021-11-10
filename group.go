package gcloudcx

import (
	"time"

	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// Group describe a Group of users
type Group struct {
	ID           uuid.UUID      `json:"id"`
	SelfURI      string         `json:"selfUri"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Description  string         `json:"description"`
	State        string         `json:"state"`
	MemberCount  int            `json:"memberCount"`
	Owners       []*User        `json:"owners"`
	Images       []*UserImage   `json:"images"`
	Addresses    []*Contact     `json:"addresses"`
	RulesVisible bool           `json:"rulesVisible"`
	Visibility   bool           `json:"visibility"`
	DateModified time.Time      `json:"dateModified"`
	Version      int            `json:"version"`
	Client       *Client        `json:"-"`
	Logger       *logger.Logger `json:"-"`
}

// Initialize initializes this from the given Client
//   implements Initializable
//   if the group ID is given in group, the group is fetched
func (group *Group) Initialize(parameters ...interface{}) error {
	context, client, logger, id, err := parseParameters(group, parameters...)
	if err != nil {
		return err
	}
	if id != uuid.Nil {
		if err := client.Get(context, NewURI("/groups/%s", id), &group); err != nil {
			return err
		}
	}
	group.Client = client
	group.Logger = logger.Topic("group").Scope("group").Record("group", group.ID)
	return nil
}

// GetID gets the identifier of this
//   implements Identifiable
func (group Group) GetID() uuid.UUID {
	return group.ID
}

// String gets a string version
//   implements the fmt.Stringer interface
func (group Group) String() string {
	if len(group.Name) > 0 {
		return group.Name
	}
	return group.ID.String()
}
