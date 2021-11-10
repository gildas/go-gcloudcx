package gcloudcx

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// NotificationChannel defines a Notification Channel
//
//   See: https://developer.mypurecloud.com/api/rest/v2/notifications/notification_service.html
type NotificationChannel struct {
	ID            uuid.UUID              `json:"id"`
	ConnectURL    *url.URL               `json:"-"`
	ExpiresOn     time.Time              `json:"expires"`
	LogHeartbeat  bool                   `json:"logHeartbeat"`
	Logger        *logger.Logger         `json:"-"`
	Client        *Client                `json:"-"`
	Socket        *websocket.Conn        `json:"-"`
	TopicReceived chan NotificationTopic `json:"-"`
}

// CreateNotificationChannel creates a new channel for notifications
//
//   If the environment variable PURECLOUD_LOG_HEARTBEAT is set to true, the Heartbeat topic will be logged
func (client *Client) CreateNotificationChannel(context context.Context) (*NotificationChannel, error) {
	var err error
	channel := &NotificationChannel{}
	if err = client.Post(context, "/notifications/channels", struct{}{}, &channel); err != nil {
		return nil, err
	}
	channel.LogHeartbeat = core.GetEnvAsBool("PURECLOUD_LOG_HEARTBEAT", false)
	channel.Client = client
	channel.Logger = client.Logger.Topic("notification_channel")
	channel.TopicReceived = make(chan NotificationTopic)
	if channel.ConnectURL != nil {
		channel.Socket, _, err = websocket.DefaultDialer.Dial(channel.ConnectURL.String(), nil)
		if err != nil {
			return nil, errors.NotConnected.With("Channel").(errors.Error).Wrap(err)
		}
	}
	// Start the message loop
	go channel.messageLoop()

	return channel, nil
}

// Close unsubscribes from all subscriptions and closes the websocket
func (channel *NotificationChannel) Close(context context.Context) (err error) {
	if channel.Client != nil && channel.Client.IsAuthorized() {
		_ = channel.Unsubscribe(context)
	}
	if channel.Socket != nil {
		close(channel.TopicReceived)
		if err = channel.Socket.Close(); err != nil {
			return errors.WithMessage(err, "Failed while closing websocket")
		}
		channel.Socket = nil
	}
	channel.ID = uuid.Nil
	return
}

// GetTopics gets all subscription topics set on this
func (channel *NotificationChannel) GetTopics(context context.Context) ([]string, error) {
	results := struct{ Entities []ChannelTopic }{}
	if err := channel.Client.Get(
		context,
		NewURI("/notifications/channels/%s/subscriptions", channel.ID),
		&results,
	); err != nil {
		return []string{}, err // err should already be decorated by Client
	}
	ids := make([]string, len(results.Entities))
	for i, entity := range results.Entities {
		ids[i] = entity.ID
	}
	return ids, nil
}

// SetTopics sets the subscriptions. It overrides any previous subscriptions
func (channel *NotificationChannel) SetTopics(context context.Context, topics ...string) ([]string, error) {
	channelTopics := make([]ChannelTopic, len(topics))
	for i, topic := range topics {
		channelTopics[i].ID = topic
	}
	results := struct {
		Entities []ChannelTopic `json:"entities"`
	}{}
	if err := channel.Client.Put(
		context,
		NewURI("/notifications/channels/%s/subscriptions", channel.ID),
		channelTopics,
		&results,
	); err != nil {
		return []string{}, err // err should already be decorated by Client
	}
	ids := make([]string, len(results.Entities))
	for i, entity := range results.Entities {
		ids[i] = entity.ID
	}
	return ids, nil
}

// IsSubscribed tells if the channel is subscribed to the given topic
func (channel *NotificationChannel) IsSubscribed(context context.Context, topic string) bool {
	topics, err := channel.GetTopics(context)
	if err != nil {
		return false
	}
	for _, t := range topics {
		if t == topic {
			return true
		}
	}
	return false
}

// Subscribe subscribes to a list of topics in the NotificationChannel
func (channel *NotificationChannel) Subscribe(context context.Context, topics ...string) ([]string, error) {
	channelTopics := make([]ChannelTopic, len(topics))
	for i, topic := range topics {
		channelTopics[i].ID = topic
	}
	results := struct {
		Entities []ChannelTopic `json:"entities"`
	}{}
	if err := channel.Client.Post(
		context,
		NewURI("/notifications/channels/%s/subscriptions", channel.ID),
		channelTopics,
		&results,
	); err != nil {
		return []string{}, err // err should already be decorated by Client
	}
	ids := make([]string, len(results.Entities))
	for i, entity := range results.Entities {
		ids[i] = entity.ID
	}
	return ids, nil
}

// Unsubscribe unsubscribes from some topics,
//
// If there is no argument, unsubscribe from all topics
func (channel *NotificationChannel) Unsubscribe(context context.Context, topics ...string) error {
	if len(topics) == 0 {
		return channel.Client.Delete(context, NewURI("/notifications/channels/%s/subscriptions", channel.ID), nil)
	}
	currentTopics, err := channel.GetTopics(context)
	if err != nil {
		return err
	}
	filteredTopics := []string{}
	for _, current := range currentTopics {
		found := false
		for _, topic := range topics {
			if current == topic {
				found = true
				break
			}
		}
		if !found {
			filteredTopics = append(filteredTopics, current)
		}
	}
	_, err = channel.SetTopics(context, filteredTopics...)
	return err
}

// MarshalJSON marshals this into JSON
func (channel NotificationChannel) MarshalJSON() ([]byte, error) {
	type surrogate NotificationChannel
	data, err := json.Marshal(struct {
		surrogate
		C *core.URL `json:"connectUri"`
	}{
		surrogate: surrogate(channel),
		C:         (*core.URL)(channel.ConnectURL),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON into this
func (channel *NotificationChannel) UnmarshalJSON(payload []byte) (err error) {
	type surrogate NotificationChannel
	var inner struct {
		surrogate
		C *core.URL `json:"connectUri"`
	}
	if err = json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*channel = NotificationChannel(inner.surrogate)
	channel.ConnectURL = (*url.URL)(inner.C)
	return
}

func (channel *NotificationChannel) messageLoop() {
	log := channel.Logger.Scope("receive")
	for {
		var err error
		var body []byte

		if _, body, err = channel.Socket.ReadMessage(); err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Infof("Websocket was closed, stopping the Channel's websocket message loop")
				return
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

// GetID gets the identifier of this
//
//   implements Identifiable
func (channel NotificationChannel) GetID() uuid.UUID {
	return channel.ID
}

// String gets a string version
//
//   implements the fmt.Stringer interface
func (channel NotificationChannel) String() string {
	return channel.ID.String()
}
