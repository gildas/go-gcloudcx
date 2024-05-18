package gcloudcx

import (
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
)

type OpenMessage interface {
	GetID() string
	core.TypeCarrier
	logger.Redactable
}

var openMessageRegistry = core.TypeRegistry{}

func UnmarshalOpenMessage(payload []byte) (OpenMessage, error) {
	message, err := openMessageRegistry.UnmarshalJSON(payload)
	if err == nil {
		return message.(OpenMessage), nil
	}
	if strings.HasPrefix(err.Error(), "Missing JSON Property") {
		return nil, errors.JSONUnmarshalError.Wrap(errors.ArgumentMissing.With("type"))
	}
	if strings.HasPrefix(err.Error(), "Unsupported Type") {
		return nil, errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(strings.TrimSuffix(strings.TrimPrefix(err.Error(), `Unsupported Type "`), `"`)))
	}
	if errors.Is(err, errors.JSONUnmarshalError) {
		return nil, err
	}
	return nil, errors.JSONUnmarshalError.Wrap(err)
}
