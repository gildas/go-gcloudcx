package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// OpenMessageTemplate describes a template for an OpenMessage
type OpenMessageTemplate struct {
	Text       string            `json:"text"`
	Parameters map[string]string `json:"parameters"`
}

// MarshalJSON marshals this into JSON
func (template OpenMessageTemplate) MarshalJSON() ([]byte, error) {
	type OpenMessageTemplateParameter struct {
		Name string `json:"name"`
		Text string `json:"text"`
	}

	var inner struct {
		Body struct {
			Text       string                         `json:"text"`
			Parameters []OpenMessageTemplateParameter `json:"parameters"`
		} `json:"body"`
	}

	inner.Body.Text = template.Text
	inner.Body.Parameters = make([]OpenMessageTemplateParameter, 0, len(template.Parameters))
	for name, text := range template.Parameters {
		inner.Body.Parameters = append(inner.Body.Parameters, OpenMessageTemplateParameter{name, text})
	}
	data, err := json.Marshal(inner)
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (template *OpenMessageTemplate) UnmarshalJSON(payload []byte) (err error) {
	type OpenMessageTemplateParameter struct {
		Name string `json:"name"`
		Text string `json:"text"`
	}

	var inner struct {
		Body struct {
			Text       string                         `json:"text"`
			Parameters []OpenMessageTemplateParameter `json:"parameters"`
		} `json:"body"`
	}

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	template.Text = inner.Body.Text
	template.Parameters = make(map[string]string)
	for _, parameter := range inner.Body.Parameters {
		template.Parameters[parameter.Name] = parameter.Text
	}
	return
}
