package purecloud

import (
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// ExtractClientAndLogger extracts a Client and a logger.Logger from its parameters
func ExtractClientAndLogger(parameters ...interface{}) (*Client, *logger.Logger, error) {
	var client *Client
	var log *logger.Logger

	for _, parameter := range parameters {
		if paramClient, ok := parameter.(*Client); ok {
			client = paramClient
		} else if paramLogger, ok := parameter.(*logger.Logger); ok {
			log = paramLogger
		}
	}
	if client == nil {
		return nil, nil, errors.ArgumentMissing.With("Client").WithStack()
	}
	if log == nil {
		if client.Logger == nil {
			return nil, nil, errors.ArgumentMissing.With("Client Logger").WithStack()
		}
		log = client.Logger
	}
	return client, log, nil
}
