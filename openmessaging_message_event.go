package gcloudcx

import (
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type OpenMessageEvent interface {
	core.TypeCarrier
}

var openMessageEventRegistry = core.TypeRegistry{}

func UnmarshalOpenMessageEvent(payload []byte) (OpenMessageEvent, error) {
	message, err := openMessageEventRegistry.UnmarshalJSON(payload, "eventType")
	if err == nil {
		return message.(OpenMessageEvent), nil
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
