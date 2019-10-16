package purecloud

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// NotificationTopic describes a Notification Topic received on a WebSocket
type NotificationTopic interface {
	Match(topicName string) bool
	Send(channel *NotificationChannel)
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

func NotificationTopicFromJSON(payload []byte) (NotificationTopic, error) {
	var header struct {
		TopicName string `json:"topicName"`
		Data      json.RawMessage
	}
	if err := json.Unmarshal(payload, &header); err != nil {
		return nil, errors.WithStack(err)
	}
	switch {
	case MetadataTopic{}.Match(header.TopicName):
		var topic MetadataTopic
		if err := json.Unmarshal(payload, &topic); err != nil {
			return nil, errors.WithStack(err)
		}
		return topic, nil
	default:
		return nil, errors.Errorf("Unsupported Topic: %s", header.TopicName)
	}
}