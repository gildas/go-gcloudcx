package purecloud

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/gildas/go-errors"
)

// NotificationTopic describes a Notification Topic received on a WebSocket
type NotificationTopic interface {
	// Match tells if the given topicName matches this topic
	Match(topicName string) bool

	// Get the PureCloud Client associated with this
	GetClient() *Client

	// Send sends the current topic to the Channel's chan
	Send(channel *NotificationChannel)

	// TopicFor builds the topicName for the given identifiables
	TopicFor(identifiables ...Identifiable) string
}

// NotificationTopicDefinition defines a Notification Topic that can subscribed to
type NotificationTopicDefinition struct {
	ID          string                 `json:"id"`
	Description string                 `json:"description"`
	Permissions []string               `json:"requiresPermissions"`
	Schema      map[string]interface{} `json:"schema"`
}

// GetNotificationAvailableTopics retrieves available notification topics
//   properties is one of more properties that should be expanded
//   see https://developer.mypurecloud.com/api/rest/v2/notifications/#get-api-v2-notifications-availabletopics
func (client *Client) GetNotificationAvailableTopics(properties ...string) ([]NotificationTopicDefinition, error) {
	query := url.Values{}
	if len(properties) > 0 {
		query.Add("expand", strings.Join(properties, ","))
	}
	results := &struct {
		Entities []NotificationTopicDefinition `json:"entities"`
	}{}
	if err := client.Get("/notifications/availabletopics?"+query.Encode(), &results); err != nil {
		return []NotificationTopicDefinition{}, err
	}
	return results.Entities, nil
}

// NotificationTopicFromJSON Unmarshal JSON into a NotificationTopic
func NotificationTopicFromJSON(payload []byte) (NotificationTopic, error) {
	var header struct {
		TopicName string `json:"topicName"`
		Data      json.RawMessage
	}
	if err := json.Unmarshal(payload, &header); err != nil {
		return nil, errors.JSONUnmarshalError.Wrap(err)
	}
	switch {
	case ConversationChatMessageTopic{}.Match(header.TopicName):
		var topic ConversationChatMessageTopic
		if err := json.Unmarshal(payload, &topic); err != nil {
			return nil, err // err should already be decorated by that struct type
		}
		return &topic, nil
	case MetadataTopic{}.Match(header.TopicName):
		var topic MetadataTopic
		if err := json.Unmarshal(payload, &topic); err != nil {
			return nil, err // err should already be decorated by that struct type
		}
		return &topic, nil
	case UserConversationChatTopic{}.Match(header.TopicName):
		var topic UserConversationChatTopic
		if err := json.Unmarshal(payload, &topic); err != nil {
			return nil, err // err should already be decorated by that struct type
		}
		return &topic, nil
	case UserPresenceTopic{}.Match(header.TopicName):
		var topic UserPresenceTopic
		if err := json.Unmarshal(payload, &topic); err != nil {
			return nil, err // err should already be decorated by that struct type
		}
		return &topic, nil
	default:
		return nil, errors.Unsupported.With("Topic", header.TopicName).WithStack()
	}
}
