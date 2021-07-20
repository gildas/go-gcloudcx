package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-errors"
)

// ACWSettings defines the After Call Work settings of a Queue
type ACWSettings struct {
	WrapupPrompt string
	Timeout      time.Duration
}

// MarshalJSON marshals this into JSON
func (settings ACWSettings) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(struct {
		Timeout      int64  `json:"timeoutMs"`
		WrapupPrompt string `json:"wrapupPrompt"`
	}{
		Timeout:      settings.Timeout.Milliseconds(),
		WrapupPrompt: settings.WrapupPrompt,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (settings *ACWSettings) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		Timeout      int64  `json:"timeoutMs"`
		WrapupPrompt string `json:"wrapupPrompt"`
	}

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	settings.Timeout = time.Duration(inner.Timeout) * time.Millisecond
	settings.WrapupPrompt = inner.WrapupPrompt
	return
}
