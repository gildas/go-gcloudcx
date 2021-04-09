package main

import (
	"context"
	"net/http"
	"net/url"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-purecloud"
)

// Config carries an application Configuration
type Config struct {
	IntegrationName         string
	IntegrationWebhookURL   *url.URL
	IntegrationWebhookToken string

	Integration *purecloud.OpenMessagingIntegration
	Client      *purecloud.Client
}

type  key int
const contextKey key = iota + 56334

// FromContext retrieves the Logger stored in the context
func ConfigFromContext(context context.Context) (*Config, error) {
	if logger, ok := context.Value(contextKey).(*Config); ok {
		return logger, nil
	}
	return nil, errors.ArgumentMissing.With("Config").WithStack()
}

// ToContext stores the Logger in the given context
func (item *Config) ToContext(parent context.Context) context.Context {
	return context.WithValue(parent, contextKey, item)
}

func (config *Config) HttpHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r*http.Request){
			next.ServeHTTP(w, r.WithContext(config.ToContext(r.Context())))
		})
	}
}
