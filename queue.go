package purecloud

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// Queue defines a PureCloud Queue
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
	SelfURI               string         `json:"selfUri"`
	Client                *Client        `json:"-"`
	Logger                *logger.Logger `json:"-"`
}

// RoutingTarget describes a routing target
type RoutingTarget struct {
	Type    string `json:"targetType,omitempty"`
	Address string `json:"targetAddress,omitempty"`
}

// FindQueueByName finds a Queue by its name
func (client *Client) FindQueueByName(name string) (*Queue, error) {
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
	err := client.Get("/routing/queues?"+query.Encode(), &response)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, queue := range response.Entities {
		if queue.Name == name {
			queue.Client = client
			queue.Logger = client.Logger.Child("queue", "queue", "queue", queue.ID)
			if queue.CreatedBy != nil {
				queue.CreatedBy.Client = client
				queue.CreatedBy.Logger = client.Logger.Child("user", "user", "user", queue.CreatedBy.ID)
			}
			return queue, nil
		}
	}
	return nil, errors.NotFound.With("queue", name).WithStack()
}

// GetID gets the identifier of this
//   implements Identifiable
func (queue Queue) GetID() uuid.UUID {
	return queue.ID
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
