package purecloud

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// NotificationChannel defines a Notification Channel
//   See: https://developer.mypurecloud.com/api/rest/v2/notifications/notification_service.html
type NotificationChannel struct {
	ID            string                 `json:"id"`
	ConnectURL    *url.URL               `json:"-"`
	ExpiresOn     time.Time              `json:"expires"`
	LogHeartbeat  bool                   `json:"logHeartbeat"`
	Logger        *logger.Logger         `json:"-"`
	Client        *Client                `json:"-"`
	Socket        *websocket.Conn        `json:"-"`
	TopicReceived chan NotificationTopic `json:"-"`
}

// CreateNotificationChannel creates a new channel for notifications
//   If the environment variable PURECLOUD_LOG_HEARTBEAT is set to true, the Heartbeat topic will be logged
func (client *Client) CreateNotificationChannel() (*NotificationChannel, error) {
	var err error
	channel := &NotificationChannel{}
	if err = client.Post("/notifications/channels", struct{}{}, &channel); err != nil {
		return nil, err
	}
	channel.LogHeartbeat  = core.GetEnvAsBool("PURECLOUD_LOG_HEARTBEAT", false)
	channel.Client        = client
	channel.Logger        = client.Logger.Topic("notification_channel")
	channel.TopicReceived = make(chan NotificationTopic)
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
		close(channel.TopicReceived)
		if err = channel.Socket.Close(); err != nil {
			return errors.WithStack(err)
		}
		channel.Socket = nil
	}
	return
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
	log := channel.Logger.Scope("receive")
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

		topic, err := NotificationTopicFromJSON(body)
		if err != nil {
			log.Warnf("%s, Body size: %d, Content: %s", err.Error(), len(body), string(body))
			continue
		}
		switch topic.(type) {
		case *MetadataTopic:
			if channel.LogHeartbeat {
				log.Tracef("Request %d bytes: %s", len(body), string(body))
			}
		default:
			log.Tracef("Request %d bytes: %s", len(body), string(body))
		}
		topic.Send(channel)
	}
}
