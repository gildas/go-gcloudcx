package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// ProcessAutomationTarget is a Process Automation Target
type ProcessAutomationTarget struct {
	ID         uuid.UUID `json:"id"`
	Type       string    `json:"type"` // Workflow
	DataFormat string    `json:"-"`    // Json, TopLevelPrimitives
}

// MarshalJSON marshals the object into JSON
func (target ProcessAutomationTarget) MarshalJSON() ([]byte, error) {
	type surrogate ProcessAutomationTarget
	type Settings struct {
		DataFormat string `json:"dataFormat"`
	}
	var settings *Settings
	if target.DataFormat != "" {
		settings = &Settings{target.DataFormat}
	}
	data, err := json.Marshal(struct {
		surrogate
		Settings *Settings `json:"workflowTargetSettings,omitempty"`
	}{
		surrogate: surrogate(target),
		Settings:  settings,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals the object from JSON
func (target *ProcessAutomationTarget) UnmarshalJSON(data []byte) error {
	type surrogate ProcessAutomationTarget
	type settings struct {
		DataFormat string `json:"dataFormat"`
	}
	var inner struct {
		surrogate
		Settings settings `json:"workflowTargetSettings"`
	}
	err := json.Unmarshal(data, &inner)
	if err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*target = ProcessAutomationTarget(inner.surrogate)
	target.DataFormat = inner.Settings.DataFormat
	return nil
}
