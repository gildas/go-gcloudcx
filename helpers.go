package purecloud

import (
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// parseParameters extracts Client, Logger, ID, from the given parameters
//
// Note, *uuid.UUID is optional and no error will be generated when it is not present
func parseParameters(parameters ...interface{}) (*Client, *logger.Logger, uuid.UUID, error) {
	var (
		client *Client
		log    *logger.Logger
		id     uuid.UUID = uuid.UUID{}
	)

	for _, parameter := range parameters {
		switch object := parameter.(type) {
		case *Client:
			client = object
		case *logger.Logger:
			log = object
		case uuid.UUID:
			id = object
		}
	}
	if client == nil {
		return nil, nil, uuid.Nil, errors.ArgumentMissing.With("Client").WithStack()
	}
	if log == nil {
		if client.Logger == nil {
			return nil, nil, uuid.Nil, errors.ArgumentMissing.With("Client Logger").WithStack()
		}
		log = client.Logger
	}
	return client, log, id, nil
}
