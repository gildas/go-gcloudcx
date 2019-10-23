package main

import (
	"github.com/pkg/errors"
	"context"
	"net/http"

	"github.com/gildas/go-purecloud"
)

// AppConfig describes An Application Configuration
type AppConfig struct {
	// WebRootPath contains the root path to prepend to URL when redirecting or in the web pages
	WebRootPath string
	// BotQueue is the Queue used for sending the initial chat to
	BotQueue    *purecloud.Queue
	// AgentQueue is the Queue used when the customer wants to talk to an agent
	AgentQueue  *purecloud.Queue
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
		return nil, errors.New("Context does not contain any AppConfing")
	}
	if config, ok := value.(*AppConfig); ok {
		return config, nil
	}
	return nil, errors.New("Invalid AppConfig stored in context")
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
	if config.AgentQueue == nil {
		return errors.New("Agent Queue is nil")
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
	return nil
}