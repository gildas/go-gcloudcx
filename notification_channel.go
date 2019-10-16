package purecloud

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gildas/go-core"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// NotificationChannel  defines a Notification Channel
type NotificationChannel struct {
	ID         string          `json:"id"`
	ConnectURL *url.URL        `json:"-"`
	ExpiresOn  time.Time       `json:"expires"`
	Client     *Client         `json:"-"`
	Socket     *websocket.Conn `json:"-"`
}

// NotificationTopic defines a Notification Topic that can subscribed to
type NotificationTopic struct {
	ID          string                 `json:"id"`
	Description string                 `json:"description"`
	Permissions []string               `json:"requiresPermissions"`
	Schema      map[string]interface{} `json:"schema"`
}

// CreateNotificationChannel creates a new channel for notifications
func (client *Client) CreateNotificationChannel() (*NotificationChannel, error) {
	var err error
	channel := &NotificationChannel{}
	if err = client.Post("/notifications/channels", struct{}{}, &channel); err != nil {
		return nil, err
	}
	channel.Client = client
	if channel.ConnectURL != nil {
		channel.Socket, _, err = websocket.DefaultDialer.Dial(channel.ConnectURL.String(), nil)
	}
	// Start the message loop
	go channel.messageLoop()

	return channel, nil
}

// Close unsubscribes from all subscriptions and closes the websocket
func (channel *NotificationChannel) Close() (err error) {
	if err = channel.Unsubscribe(); err != nil {
		return err
	}
	if channel.Socket != nil {
		if err = channel.Socket.Close(); err != nil {
			return errors.WithStack(err)
		}
		channel.Socket = nil
	}
	return
}

// GetNotificationAvailableTopics retrieves available notification topics
//   properties is one of more properties that should be expanded
//   see https://developer.mypurecloud.com/api/rest/v2/notifications/#get-api-v2-notifications-availabletopics
func (client *Client) GetNotificationAvailableTopics(properties ...string) ([]NotificationTopic, error) {
	query := url.Values{}
	if len(properties) > 0 {
		query.Add("expand", strings.Join(properties, ","))
	}
	results := &struct {
		Entities []NotificationTopic `json:"entities"`
	}{}
	if err := client.Get("/notifications/availabletopics?"+query.Encode(), &results); err != nil {
		return []NotificationTopic{}, err
	}
	return results.Entities, nil
}

// Subscribe subscribes to a list of topics in the NotificationChannel
func (channel *NotificationChannel) Subscribe(topics ...string) ([]string, error) {
	type idHolder struct {ID string `json:"id"`}
	channelTopics := make([]idHolder, len(topics))
	for i, topic := range topics {
		channelTopics[i].ID = topic
	}
	results := &struct {
		Entities []idHolder `json:"entities"`
	}{}
	if err := channel.Client.Post(
		fmt.Sprintf("/notifications/channels/%s/subscriptions", channel.ID),
		channelTopics,
		&results,
	); err != nil {
		return []string{}, errors.WithStack(err)
	}
	ids := make([]string, len(results.Entities))
	for i, entity := range results.Entities {
		ids[i] = entity.ID
	}
	return ids, nil
}

// Unsubscribe unsubscribes from all topics
func (channel *NotificationChannel) Unsubscribe() error {
	return channel.Client.Delete(fmt.Sprintf("/notifications/channels/%s/subscriptions", channel.ID), nil)
}

// MarshalJSON marshals this into JSON
func (channel NotificationChannel) MarshalJSON() ([]byte, error) {
	type surrogate NotificationChannel
	return json.Marshal(struct {
		surrogate
		C *core.URL `json:"connectUri"`
	}{
		surrogate: surrogate(channel),
		C:         (*core.URL)(channel.ConnectURL),
	})
}

// UnmarshalJSON unmarshals JSON into this
func (channel *NotificationChannel) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NotificationChannel
	var inner struct {
		surrogate
		C *core.URL `json:"connectUri"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.WithStack(err)
	}
	*channel = NotificationChannel(inner.surrogate)
	channel.ConnectURL = (*url.URL)(inner.C)
	return
}

func (channel *NotificationChannel) messageLoop() (err error) {
	log := channel.Client.Logger.Scope("receive")
	for {
		var body []byte

		if _, body, err = channel.Socket.ReadMessage(); err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Infof("Websocket was closed, stopping the message loop")
				return nil
			}
			log.Errorf("Failed to read incoming message", err)
			continue
		}
		log.Tracef("Received %d bytes: %s", len(body), string(body))
	}
}
