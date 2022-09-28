package gcloudcx

import (
	"time"

	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// Group describe a Group of users
type Group struct {
	ID           uuid.UUID      `json:"id"`
	SelfURI      URI            `json:"selfUri"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Description  string         `json:"description"`
	State        string         `json:"state"`
	MemberCount  int            `json:"memberCount"`
	Owners       []*User        `json:"owners"`
	Images       []*UserImage   `json:"images"`
	Addresses    []*Contact     `json:"addresses"`
	RulesVisible bool           `json:"rulesVisible"`
	Visibility   string         `json:"visibility"`
	DateModified time.Time      `json:"dateModified"`
	Version      int            `json:"version"`
	client       *Client        `json:"-"`
	logger       *logger.Logger `json:"-"`
}

// Initialize initializes the object
//
// accepted parameters: *gcloufcx.Client, *logger.Logger
//
// implements Initializable
func (group *Group) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case *Client:
			group.client = parameter
		case *logger.Logger:
			group.logger = parameter.Child("group", "group", "id", group.ID)
		}
	}
}

// GetID gets the identifier of this
//
// implements Identifiable
func (group Group) GetID() uuid.UUID {
	return group.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (group Group) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/groups/%s", ids[0])
	}
	if group.ID != uuid.Nil {
		return NewURI("/api/v2/groups/%s", group.ID)
	}
	return URI("/api/v2/groups/")
}

// String gets a string version
//
// implements the fmt.Stringer interface
func (group Group) String() string {
	if len(group.Name) > 0 {
		return group.Name
	}
	return group.ID.String()
}
