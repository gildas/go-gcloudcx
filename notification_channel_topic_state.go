package gcloudcx

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// NotificationChannelTopicState describes a Topic subscription channel
//
// See https://developer.genesys.cloud/api/rest/v2/notifications/notification_service#topic-subscriptions
type NotificationChannelTopicState struct {
	Topic        NotificationTopic `json:"-"`
	State        string            `json:"state,omitempty"`
	RejectReason string            `json:"rejectReason,omitempty"`
	SelfURI      URI               `json:"selfUri,omitempty"`
}

// GetURI gets the URI of this
//
//	implements Addressable
func (state NotificationChannelTopicState) GetURI() URI {
	return state.SelfURI
}

// Contains tells if this contains the given topic
func (state NotificationChannelTopicState) Contains(topic NotificationTopic) bool {
	return state.Topic.GetType() == topic.GetType()
}

// MarshalJSON marshals into JSON
//
// implements json.Marshaler
func (state NotificationChannelTopicState) MarshalJSON() ([]byte, error) {
	type surrogate NotificationChannelTopicState
	data, err := json.Marshal(struct {
		ID string `json:"id"` // ID is a string of the form v2.xxx.uuid.yyy
		surrogate
	}{
		ID:        state.Topic.String(),
		surrogate: surrogate(state),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals from JSON
//
// implements json.Unmarshaler
func (state *NotificationChannelTopicState) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NotificationChannelTopicState
	var inner struct {
		surrogate
		ID string `json:"id"` // ID is a string of the form v2.xxx.uuid.yyy
	}
	if err = json.Unmarshal(payload, &inner); errors.Is(err, errors.JSONUnmarshalError) {
		return err
	} else if err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*state = NotificationChannelTopicState(inner.surrogate)
	state.Topic, err = NotificationTopicFrom(inner.ID)
	return errors.JSONUnmarshalError.Wrap(err)
}
