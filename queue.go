package gcloudcx

import (
	"encoding/json"
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

// Initialize initializes the object
//
// accepted parameters: *gcloufcx.Client, *logger.Logger
//
// implements Initializable
func (queue *Queue) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case *Client:
			queue.client = parameter
		case *logger.Logger:
			queue.logger = parameter.Child("queue", "queue", "id", queue.ID)
		}
	}
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
