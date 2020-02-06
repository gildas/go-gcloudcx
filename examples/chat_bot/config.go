package main

import (
	"net/url"
	"github.com/gildas/go-logger"
	"context"
	"net/http"
	"strings"

	"github.com/gildas/go-purecloud"
	"github.com/gildas/go-errors"
)

// AppConfig describes An Application Configuration
type AppConfig struct {
	// WebRootPath contains the root path to prepend to URL when redirecting or in the web pages
	WebRootPath string

	// BotURL is the Chat Bot's url
	BotURL      *url.URL

	// BotQueue is the Queue used for sending the initial chat to
	BotQueue    *purecloud.Queue

	// AgentQueue is the Queue used when the customer wants to talk to an agent
	AgentQueue  *purecloud.Queue

	// User is the currently Logged in User
	User *purecloud.User

	// NotificationChannel is the channel used to receive notifications from PureCloud
	NotificationChannel *purecloud.NotificationChannel

	// Logger is the Logger for the configuration
	Logger *logger.Logger
}

// AppConfigKey is the key to store AppConfig in context.Context
const AppConfigKey = iota + 60506

// ToContext store this AppConfig in the given context
func (config *AppConfig) ToContext(parent context.Context) context.Context {
	return context.WithValue(parent, AppConfigKey, config)
}

// AppConfigFromContext retrieves an AppConfig from the given context
func AppConfigFromContext(context context.Context) (*AppConfig, error) {
	value := context.Value(AppConfigKey)
	if value == nil {
		return nil, errors.ArgumentMissing.With("AppConfig").WithStack()
	}
	if config, ok := value.(*AppConfig); ok {
		return config, nil
	}
	return nil, errors.ArgumentInvalid.With("AppConfig", value).WithStack()
}

// HttpHandler wraps this AppConfig into an http.Handler
func (config *AppConfig) HttpHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(config.ToContext(r.Context())))
		})
	}
}

// Initialize configures this AppConfig by calling PureCloud as needed
func (config *AppConfig) Initialize(client *purecloud.Client) (err error) {
	config.Logger = client.Logger.Topic("config")
	log := config.Logger.Scope("initialize")
	if config.AgentQueue == nil {
		return errors.ArgumentMissing.With("Queue").WithStack()
	}
	if len(config.AgentQueue.ID) == 0 {
		queueName := config.AgentQueue.Name
		config.AgentQueue, err = client.FindQueueByName(queueName)
		if err != nil {
			return errors.Wrapf(err, "Failed to retrieve the Agent Queue %s", queueName)
		}
	}
	if config.BotQueue == nil {
		return errors.New("Bot Queue is nil")
	}
	if len(config.BotQueue.ID) == 0 {
		queueName := config.BotQueue.Name
		config.BotQueue, err = client.FindQueueByName(queueName)
		if err != nil {
			return errors.Wrapf(err, "Failed to retrieve the Bot Queue %s", queueName)
		}
	}

	config.User, err = client.GetMyUser()
	if err != nil {
		log.Errorf("Failed to retrieve my User", err)
		return
	}
	log.Infof("Current User: %s", config.User)

	config.NotificationChannel, err = client.CreateNotificationChannel()
	if err != nil {
		log.Errorf("Failed to create a notification channel", err)
		return
	}

	topics, err := config.NotificationChannel.Subscribe(
		purecloud.UserPresenceTopic{}.TopicFor(config.User),
		purecloud.UserConversationChatTopic{}.TopicFor(config.User),
	)
	if err != nil {
		log.Errorf("Failed to subscribe to topics", err)
		return
	}
	log.Infof("Subscribed to topics: [%s]", strings.Join(topics, ","))

	// Call the PureCloud Notification Topic loop
	go MessageLoop(config)

	return nil
}

func (config *AppConfig) Reset() (err error) {
	log := config.Logger.Scope("reset")

	if config.NotificationChannel != nil {
		log.Debugf("Closing Notification Channel %s", config.NotificationChannel)
		err = config.NotificationChannel.Close()
		if err != nil {
			config.NotificationChannel.Logger.Errorf("Failed to close the notification channel %s", config.NotificationChannel, err)
			return err
		}
		log.Infof("Closed Notification Channel %s", config.NotificationChannel)
		config.NotificationChannel = nil
	}
	config.User = nil
	log.Infof("Config is reset")
	return
}