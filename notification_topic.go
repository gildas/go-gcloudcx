package gcloudcx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// NotificationTopic describes a Notification Topic received on a WebSocket
type NotificationTopic interface {
	core.TypeCarrier
	fmt.Stringer

	// With creates a new NotificationTopic with the given targets
	With(targets ...Identifiable) NotificationTopic

	// GetTargets gets the targets of this topic
	GetTargets() []Identifiable
}

// NotificationTopicDefinition defines a Notification Topic that can subscribed to
type NotificationTopicDefinition struct {
	ID          string                 `json:"id"`
	Description string                 `json:"description"`
	Permissions []string               `json:"requiresPermissions"`
	Schema      map[string]interface{} `json:"schema"`
}

var notificationTopicRegistry = core.TypeRegistry{}

// GetAvailableNotificationTopics retrieves available notification topics
//
//	properties is one of more properties that should be expanded
//	see https://developer.mypurecloud.com/api/rest/v2/notifications/#get-api-v2-notifications-availabletopics
func (client *Client) GetAvailableNotificationTopics(context context.Context, properties ...string) (definitions []NotificationTopicDefinition, correlationID string, err error) {
	query := url.Values{}
	if len(properties) > 0 {
		query.Add("expand", strings.Join(properties, ","))
	}
	results := &struct {
		Entities []NotificationTopicDefinition `json:"entities"`
	}{}
	if correlationID, err = client.Get(context, NewURI("/notifications/availabletopics?%s", query.Encode()), &results); err != nil {
		return []NotificationTopicDefinition{}, correlationID, err
	}
	return results.Entities, correlationID, nil
}

// topicNameWith builds the topicName for the given identifiables
func topicNameWith(topic NotificationTopic, identifiables ...Identifiable) string {
	var topicName strings.Builder
	components := strings.Split(topic.GetType(), "{id}")

	topicName.WriteString(components[0])
	for i, identifiable := range identifiables {
		topicName.WriteString(identifiable.GetID().String())
		if i < len(components)-1 {
			topicName.WriteString(components[i+1])
		}
	}
	return topicName.String()
}

// NotificationTopicFrom builds a NotificationTopic from the given topicName
func NotificationTopicFrom(topicName string) (NotificationTopic, error) {
	for _, topic := range notificationTopicRegistry {
		result := reflect.New(topic).Interface()
		resultType := result.(core.TypeCarrier).GetType()
		if found, targets := getTargets(resultType, topicName); found {
			return result.(NotificationTopic).With(targets...), nil
		}
	}
	return nil, errors.ArgumentInvalid.With("topicName", topicName)
}

func getTargets(topicType, topicName string) (found bool, targets []Identifiable) {
	if topicType == topicName {
		return true, targets
	}
	components := strings.Split(topicType, "{id}")
	for _, component := range components {
		if !strings.Contains(topicName, component) {
			return false, []Identifiable{}
		}
		value := strings.TrimPrefix(topicName, component)
		if index := strings.Index(value, "."); index > 0 {
			value = value[:index]
		}
		if id, err := uuid.Parse(value); err == nil {
			targets = append(targets, EntityRef{id})
		}
	}
	return true, targets
}

// UnmarshalNotificationTopic Unmarshal JSON into a NotificationTopic
//
// The result is a NotificationTopic that can be casted into the appropriate type
func UnmarshalNotificationTopic(payload []byte) (NotificationTopic, error) {
	var header struct {
		TopicName string `json:"topicName"`
	}

	if err := json.Unmarshal(payload, &header); err != nil {
		return nil, errors.JSONUnmarshalError.Wrap(err)
	}
	for _, topic := range notificationTopicRegistry {
		result := reflect.New(topic).Interface()
		resultType := result.(core.TypeCarrier).GetType()
		if found, targets := getTargets(resultType, header.TopicName); found {
			if err := json.Unmarshal(payload, result); err != nil {
				return nil, errors.JSONUnmarshalError.Wrap(err)
			}
			return result.(NotificationTopic).With(targets...), nil
		}
	}
	return nil, errors.Unsupported.With("Topic", header.TopicName)
}
