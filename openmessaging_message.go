package gcloudcx

import (
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// OpenMessage is a message sent or received by the Open Messaging API
//
// See: https://developer.genesys.cloud/commdigital/digital/openmessaging/normalizedmsgformat#openoutboundnormalizedmessage-json-schema
type OpenMessage interface {
	GetID() string
	core.TypeCarrier
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
		supportedTypes := make([]string, 0, len(openMessageRegistry))
		for key := range openMessageRegistry {
			supportedTypes = append(supportedTypes, key)
		}
		return nil, errors.JSONUnmarshalError.Wrap(
			errors.InvalidType.With(
				strings.TrimSuffix(strings.TrimPrefix(err.Error(), `Unsupported Type "`), `"`),
				strings.Join(supportedTypes, ","),
			),
		)
	}
	if errors.Is(err, errors.JSONUnmarshalError) {
		return nil, err
	}
	return nil, errors.JSONUnmarshalError.Wrap(err)
}
