package assertions

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func decodeWithMetadata(assertion Assertion, result interface{}) (mapstructure.Metadata, error) {
	var decoderMetadata mapstructure.Metadata
	decoderConfig := &mapstructure.DecoderConfig{
		Metadata: &decoderMetadata,
		Result:   &result,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return mapstructure.Metadata{}, fmt.Errorf("error creating decoder: %s", err)
	}

	err = decoder.Decode(assertion.Metadata)
	if err != nil {
		return mapstructure.Metadata{}, fmt.Errorf("error decoding assertion metadata: %s", err)
	}

	return decoderMetadata, nil
}
