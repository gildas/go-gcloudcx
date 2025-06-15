package gcloudcx

import (
	"context"
	"encoding/json"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// ProcessAutomationTrigger is a Process Automation Trigger
//
// See: https://developer.genesys.cloud/commdigital/taskmanagement/task-management-maestro
//
// See: https://developer.genesys.cloud/platform/process-automation/trigger-apis
type ProcessAutomationTrigger struct {
	ID              uuid.UUID                   `json:"id"`
	Name            string                      `json:"name"`
	Description     string                      `json:"description,omitempty"`
	Topic           NotificationTopic           `json:"-"`
	MatchCriteria   []ProcessAutomationCriteria `json:"matchCriteria"`
	Target          ProcessAutomationTarget     `json:"target"`
	Enabled         bool                        `json:"enabled"`
	EventTTLSeconds int                         `json:"eventTTLSeconds,omitempty"`
	DelayBySeconds  int                         `json:"delayBySeconds,omitempty"`
	Version         int                         `json:"version,omitempty"`
	client          *Client                     `json:"-"`
	logger          *logger.Logger              `json:"-"`
}

// Initialize initializes the object
//
// accepted parameters: *gcloudcx.Client, *logger.Logger
//
// implements Initializable
func (trigger *ProcessAutomationTrigger) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case uuid.UUID:
			trigger.ID = parameter
		case *Client:
			trigger.client = parameter
		case *logger.Logger:
			trigger.logger = parameter.Child("trigger", "trigger", "id", trigger.ID)
		}
	}
	if trigger.logger == nil {
		trigger.logger = logger.Create("gcloudcx", &logger.NilStream{})
	}
}

// GetID gets the identifier of this
//
// implements Identifiable
func (trigger ProcessAutomationTrigger) GetID() uuid.UUID {
	return trigger.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (trigger ProcessAutomationTrigger) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/processautomation/triggers/%s", ids[0])
	}
	if trigger.ID != uuid.Nil {
		return NewURI("/api/v2/processautomation/triggers/%s", trigger.ID)
	}
	return URI("/api/v2/processautomation/triggers/")
}

// Create creates a new ProcessAutomationTrigger
//
// The Genesys Cloud correlation ID is return as the second return value.
//
// See: https://developer.genesys.cloud/platform/process-automation/trigger-apis#post-api-v2-processautomation-triggers
func (trigger ProcessAutomationTrigger) Create(context context.Context, client *Client) (*ProcessAutomationTrigger, string, error) {
	created := ProcessAutomationTrigger{}
	correlationID, err := client.Post(
		context,
		"processautomation/triggers",
		trigger,
		&created,
	)
	if err != nil {
		return nil, correlationID, err
	}
	created.client = client
	created.logger = client.Logger.Child("trigger", "trigger", "id", trigger.ID)
	return &created, correlationID, nil
}

// Delete deletes this ProcessAutomationTrigger
func (trigger ProcessAutomationTrigger) Delete(context context.Context) (correlationID string, err error) {
	return trigger.client.Delete(context, trigger.GetURI(), nil)
}

// String gets a string version
//
// implements the fmt.Stringer interface
func (trigger ProcessAutomationTrigger) String() string {
	return trigger.ID.String()
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (trigger ProcessAutomationTrigger) MarshalJSON() ([]byte, error) {
	type surrogate ProcessAutomationTrigger
	data, err := json.Marshal(&struct {
		surrogate
		ID        string `json:"id,omitempty"`
		TopicName string `json:"topicName"`
	}{
		surrogate: surrogate(trigger),
		ID:        core.UUID(trigger.ID).String(),
		TopicName: trigger.Topic.GetType(),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (trigger *ProcessAutomationTrigger) UnmarshalJSON(payload []byte) (err error) {
	type surrogate ProcessAutomationTrigger
	var inner struct {
		surrogate
		TopicName string `json:"topicName"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*trigger = ProcessAutomationTrigger(inner.surrogate)
	trigger.Topic, err = NotificationTopicFrom(inner.TopicName)
	return errors.JSONUnmarshalError.Wrap(err)
}
