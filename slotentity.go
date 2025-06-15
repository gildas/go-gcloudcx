package gcloudcx

import (
	"fmt"
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// TODO: This will need to go to go-gcloudcx

// SlotEntity represents a slot entity in a Genesys Cloud CX message
type SlotEntity interface {
	core.TypeCarrier
	core.Named
	fmt.Stringer                                 // used to get the string representation of the entity's value
	ParseValue(value string) (SlotEntity, error) // ParseValue parses the value and returns a new SlotEntity instance
}

var slotEntityRegistry = core.TypeRegistry{}

// UnmarshalSlotEntity unmarshals a slot entity from JSON
func UnmarshalSlotEntity(payload []byte) (SlotEntity, error) {
	entity, err := slotEntityRegistry.UnmarshalJSON(payload)
	if err == nil {
		return entity.(SlotEntity), nil
	}
	if strings.HasPrefix(err.Error(), "Missing JSON Property") {
		return nil, errors.JSONUnmarshalError.Wrap(errors.ArgumentMissing.With("type"))
	}
	if strings.HasPrefix(err.Error(), "Unsupported Type") {
		return nil, errors.JSONUnmarshalError.Wrap(
			errors.InvalidType.With(
				strings.TrimSuffix(strings.TrimPrefix(err.Error(), `Unsupported Type "`), `"`),
				strings.Join(slotEntityRegistry.SupportedTypes(), ","),
			),
		)
	}
	if errors.Is(err, errors.JSONUnmarshalError) {
		return nil, err
	}
	return nil, errors.JSONUnmarshalError.Wrap(err)
}
