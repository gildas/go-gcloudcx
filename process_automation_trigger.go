package gcloudcx

import (
	"context"
	"encoding/json"

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
	Description     string                      `json:"description"`
	TopicName       string                      `json:"topicName"`
	MatchCriteria   []ProcessAutomationCriteria `json:"matchCriteria"`
	Target          ProcessAutomationTarget     `json:"target"`
	Enabled         bool                        `json:"enabled"`
	EventTTLSeconds int                         `json:"eventTTLSeconds"`
	DelayBySeconds  int                         `json:"delayBySeconds"`
	Version         int                         `json:"version"`
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
		case *Client:
			trigger.client = parameter
		case *logger.Logger:
			trigger.logger = parameter.Child("trigger", "trigger", "id", trigger.ID)
		}
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
// See: https://developer.genesys.cloud/platform/process-automation/trigger-apis#post-api-v2-processautomation-triggers
func (client *Client) CreateProcessAutomationTrigger(context context.Context, name, description, topicName string, target ProcessAutomationTarget, matchCriteria []ProcessAutomationCriteria, eventTTLSeconds, delayBySeconds int, enabled bool) (*ProcessAutomationTrigger, error) {
	trigger := ProcessAutomationTrigger{}
	err := client.Post(
		context,
		"processautomation/triggers",
		struct {
			Name            string                      `json:"name"`
			Description     string                      `json:"description,omitempty"`
			TopicName       string                      `json:"topicName"`
			MatchCriteria   []ProcessAutomationCriteria `json:"matchCriteria"`
			Target          ProcessAutomationTarget     `json:"target"`
			Enabled         bool                        `json:"enabled"`
			EventTTLSeconds int                         `json:"eventTTLSeconds,omitempty"`
			DelayBySeconds  int                         `json:"delayBySeconds,omitempty"`
		}{
			Name:            name,
			Description:     description,
			TopicName:       topicName,
			MatchCriteria:   matchCriteria,
			Target:          target,
			Enabled:         enabled,
			EventTTLSeconds: eventTTLSeconds,
			DelayBySeconds:  delayBySeconds,
		},
		&trigger,
	)
	if err != nil {
		return nil, err
	}
	trigger.client = client
	trigger.logger = client.Logger.Child("trigger", "trigger", "id", trigger.ID)
	return &trigger, nil
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
		SelfURI URI `json:"selfUri"`
	}{
		surrogate: surrogate(trigger),
		SelfURI:   trigger.GetURI(),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}
