package gcloudcx

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-errors"
)

// OpenMessageTypingEvent is a typing event sent or received by the Open Messaging API
type OpenMessageTypingEvent struct {
	IsTyping bool          `json:"-"`
	Duration time.Duration `json:"-"`
}

func init() {
	openMessageEventRegistry.Add(OpenMessageTypingEvent{})
}

// GetType returns the type of this event
//
// implements core.TypeCarrier
func (event OpenMessageTypingEvent) GetType() string {
	return "Typing"
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (event OpenMessageTypingEvent) MarshalJSON() (data []byte, err error) {
	if !event.IsTyping || event.Duration > 0 {
		type TypingInfo struct {
			Type     string `json:"type"`
			Duration int    `json:"duration"`
		}
		newTypingInfo := func(isTyping bool, duration time.Duration) TypingInfo {
			if !isTyping {
				return TypingInfo{
					Type:     "On",
					Duration: int(duration.Milliseconds()),
				}
			}
			return TypingInfo{
				Type:     "Off",
				Duration: int(duration.Milliseconds()),
			}
		}
		data, err = json.Marshal(struct {
			Type   string     `json:"eventType"`
			Typing TypingInfo `json:"typing"`
		}{
			Type:   event.GetType(),
			Typing: newTypingInfo(event.IsTyping, event.Duration),
		})
	} else {
		data, err = json.Marshal(struct {
			Type string `json:"eventType"`
		}{
			Type: event.GetType(),
		})
	}
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
//
// implements json.Unmarshaler
func (event *OpenMessageTypingEvent) UnmarshalJSON(payload []byte) (err error) {
	type surrogate OpenMessageTypingEvent
	var inner struct {
		surrogate
		Type   string `json:"eventType"`
		Typing *struct {
			Type     string `json:"type"`
			Duration int    `json:"duration"`
		} `json:"typing"`
	}
	if err = json.Unmarshal(payload, &inner); errors.Is(err, errors.JSONUnmarshalError) {
		return err
	} else if err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*event = OpenMessageTypingEvent(inner.surrogate)

	if inner.Typing != nil {
		event.IsTyping = inner.Typing.Type == "On"
		event.Duration = time.Duration(inner.Typing.Duration) * time.Millisecond
	} else {
		event.IsTyping = true
	}
	return nil
}
