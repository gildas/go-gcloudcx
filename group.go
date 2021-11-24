package gcloudcx

import (
	"context"
	"time"

	"github.com/gildas/go-errors"
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
	Visibility   bool           `json:"visibility"`
	DateModified time.Time      `json:"dateModified"`
	Version      int            `json:"version"`
	client       *Client        `json:"-"`
	logger       *logger.Logger `json:"-"`
}

// Fetch fetches a group
//
// implements Fetchable
func (group *Group) Fetch(ctx context.Context, client *Client, parameters ...interface{}) error {
	id, name, selfURI, log := client.ParseParameters(ctx, group, parameters...)

	if id != uuid.Nil {
		if err := client.Get(ctx, NewURI("/groups/%s", id), &group); err != nil {
			return err
		}
		group.logger = log
	} else if len(selfURI) > 0 {
		if err := client.Get(ctx, selfURI, &group); err != nil {
			return err
		}
		group.logger = log.Record("id", group.ID)
	} else if len(name) > 0 {
		return errors.NotImplemented.WithStack()
	}
	group.client = client
	return nil
}

// GetID gets the identifier of this
//   implements Identifiable
func (group Group) GetID() uuid.UUID {
	return group.ID
}

// GetURI gets the URI of this
//   implements Addressable
func (group Group) GetURI() URI {
	return group.SelfURI
}

// String gets a string version
//   implements the fmt.Stringer interface
func (group Group) String() string {
	if len(group.Name) > 0 {
		return group.Name
	}
	return group.ID.String()
}
