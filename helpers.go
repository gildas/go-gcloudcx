package purecloud

import (
	"github.com/gildas/go-logger"
	"github.com/pkg/errors"
)

// ExtractClientAndLogger extracts a Client and a logger.Logger from its parameters
func ExtractClientAndLogger(object interface{}, parameters ...interface{}) (*Client, *logger.Logger, error) {
	var client *Client
	var log    *logger.Logger

	for _, parameter := range parameters {
		if paramClient, ok := parameter.(*Client); ok {
			client = paramClient
		}
		if paramLogger, ok := parameter.(*logger.Logger); ok {
			log = paramLogger
		}
	}
	if client == nil {
		return nil, nil, errors.New("Missing Client")
	}
	if log == nil {
		if client.Logger == nil {
			return nil, nil, errors.New("Missing Client Logger")
		}
		log = client.Logger
	}
	return client, log, nil
}