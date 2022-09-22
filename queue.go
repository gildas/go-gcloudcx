package gcloudcx

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// Queue defines a GCloud Queue
type Queue struct {
	ID                    uuid.UUID      `json:"id"`
	Name                  string         `json:"name"`
	CreatedBy             *User          `json:"-"`
	ModifiedBy            string         `json:"modifiedBy"`
	DateCreated           time.Time      `json:"dateCreated"`
	Division              *Division      `json:"division"`
	MemberCount           int            `json:"memberCount"`
	MediaSettings         MediaSettings  `json:"mediaSettings"`
	ACWSettings           ACWSettings    `json:"acwSettings"`
	SkillEvaluationMethod string         `json:"skillEvaluationMethod"`
	AutoAnswerOnly        bool           `json:"true"`
	DefaultScripts        interface{}    `json:"defaultScripts"`
	SelfURI               URI            `json:"selfUri"`
	client                *Client        `json:"-"`
	logger                *logger.Logger `json:"-"`
}

// RoutingTarget describes a routing target
type RoutingTarget struct {
	Type    string `json:"targetType,omitempty"`
	Address string `json:"targetAddress,omitempty"`
}

// Fetch fetches a queue
//
// implements Fetchable
func (queue *Queue) Fetch(ctx context.Context, client *Client, parameters ...interface{}) error {
	id, name, selfURI, log := client.ParseParameters(ctx, queue, parameters...)

	if id != uuid.Nil {
		if err := client.Get(ctx, NewURI("/groups/%s", id), &queue); err != nil {
			return err
		}
		queue.logger = log
	} else if len(selfURI) > 0 {
		if err := client.Get(ctx, selfURI, &queue); err != nil {
			return err
		}
		queue.logger = log.Record("id", queue.ID)
	} else if len(name) > 0 {
		return errors.NotImplemented.WithStack()
	}
	queue.client = client
	return nil
}

// FindQueueByName finds a Queue by its name
func (client *Client) FindQueueByName(context context.Context, name string) (*Queue, error) {
	response := struct {
		Entities   []*Queue `json:"entities"`
		PageSize   int64    `json:"pageSize"`
		PageNumber int64    `json:"pageNumber"`
		PageCount  int64    `json:"pageCount"`
		PageTotal  int64    `json:"pageTotal"`
		SelfURI    string   `json:"selfUri"`
		FirstURI   string   `json:"firstUrl"`
		LastURI    string   `json:"lastUri"`
	}{}
	query := url.Values{}
	query.Add("name", name)
	err := client.Get(context, NewURI("/routing/queues?%s", query.Encode()), &response)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, queue := range response.Entities {
		if queue.Name == name {
			queue.client = client
			queue.logger = client.Logger.Child("queue", "queue", "id", queue.ID)
			if queue.CreatedBy != nil {
				queue.CreatedBy.client = client
				queue.CreatedBy.logger = client.Logger.Child("user", "user", "id", queue.CreatedBy.ID)
			}
			return queue, nil
		}
	}
	// TODO: read all pages!!!
	return nil, errors.NotFound.With("queue", name)
}

// GetID gets the identifier of this
//   implements Identifiable
func (queue Queue) GetID() uuid.UUID {
	return queue.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (queue Queue) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/routing/queues/%s", ids[0])
	}
	if queue.ID != uuid.Nil {
		return NewURI("/api/v2/routing/queues/%s", queue.ID)
	}
	return URI("/api/v2/routing/queues/")
}

func (queue Queue) String() string {
	if len(queue.Name) > 0 {
		return queue.Name
	}
	return queue.ID.String()
}

// UnmarshalJSON unmarshals JSON into this
func (queue *Queue) UnmarshalJSON(payload []byte) (err error) {
	type surrogate Queue
	var inner struct {
		surrogate
		CreatedByID uuid.UUID `json:"createdBy"`
	}

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*queue = Queue(inner.surrogate)
	if len(inner.CreatedByID) > 0 {
		queue.CreatedBy = &User{ID: inner.CreatedByID}
	}
	return
}
