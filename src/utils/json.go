package utils

import (
	"encoding/json"
	"fmt"
)

type GenericMappable interface {
	ToGenericMap() map[string]interface{}
}

func ConvertToGenericInterface(data any) interface{} {
	switch t := data.(type) {
	case map[interface{}]interface{}:
		interfaceToInterfaceMap := map[string]interface{}{}
		for k, v := range t {
			interfaceToInterfaceMap[k.(string)] = ConvertToGenericInterface(v)
		}
		return interfaceToInterfaceMap
	case map[string]interface{}:
		stringToInterfaceMap := map[string]interface{}{}
		for k, v := range t {
			stringToInterfaceMap[k] = ConvertToGenericInterface(v)
		}
		return stringToInterfaceMap
	case []interface{}:
		for i, v := range t {
			t[i] = ConvertToGenericInterface(v)
		}
	default:
		return data
	}

	return data
}

// ToJSON marshals the given data into a JSON string.
func ToJSON(data GenericMappable) (string, error) {
	dataAsInterface := ConvertToGenericInterface(data.ToGenericMap())
	json, err := json.Marshal(dataAsInterface)

	if err != nil {
		return "", fmt.Errorf("error while trying to marshal %#v to JSON: %s", data, err)
	}

	return string(json), nil
}

// FromJSON unmarshals the given JSON string into the given data structure.
func FromJSON(jsonString string, data any) error {
	err := json.Unmarshal([]byte(jsonString), data)

	if err != nil {
		return fmt.Errorf("error while trying to unmarshal JSON string %s: %s", jsonString, err)
	}

	return nil
}
