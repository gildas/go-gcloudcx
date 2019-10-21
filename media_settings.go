package purecloud

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// MediaSetting defines a media setting in a Queue
type MediaSetting struct {
	AlertingTimeout time.Duration
	ServiceLevel    ServiceLevel
}

// MediaSettings is a map of media names and settings
type MediaSettings map[string]MediaSetting

// MarshalJSON marshals this into JSON
func (setting MediaSetting) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		AlertingTimeout int64        `json:"durationMs"`
		ServiceLevel    ServiceLevel `json:"serviceLevel"`
	}{
		AlertingTimeout: setting.AlertingTimeout.Milliseconds(),
		ServiceLevel:    setting.ServiceLevel,
	})
}

// UnmarshalJSON unmarshals JSON into this
func (setting *MediaSetting) UnmarshalJSON(payload []byte) (err error) {
	var inner struct {
		AlertingTimeout int64        `json:"durationMs"`
		ServiceLevel    ServiceLevel `json:"serviceLevel"`
	}

	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.WithStack(err)
	}
	setting.AlertingTimeout = time.Duration(inner.AlertingTimeout) * time.Millisecond
	setting.ServiceLevel    = inner.ServiceLevel
	return
}