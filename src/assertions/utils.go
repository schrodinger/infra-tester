package assertions

import "fmt"

func partialDeepCompare(a, b interface{}) error {
	switch typedA := a.(type) {
	case bool:
		return partialDeepCompareBool(typedA, b)
	case float64:
		return partialDeepCompareFloat64(typedA, b)
	case int:
		return partialDeepCompareInt(typedA, b)
	case string:
		return partialDeepCompareString(typedA, b)
	case []interface{}:
		return partialDeepCompareSlice(typedA, b)
	case map[string]interface{}:
		return partialDeepCompareMap(typedA, b)
	default:
		return fmt.Errorf("type %T is not supported, please raise an issue in the infra-tester GitHub repo", a)
	}
}

func partialDeepCompareBool(typedA bool, b interface{}) error {
	typedB, ok := b.(bool)
	if !ok {
		return fmt.Errorf("%+v of type %T could not be cast to bool", b, b)
	}

	if typedA == typedB {
		return nil
	}

	return fmt.Errorf("%+v does not equal %+v", typedA, typedB)
}

func partialDeepCompareFloat64(typedA float64, b interface{}) error {
	switch typedB := b.(type) {
	// Handle both int and float64 types for second data type
	case float64:
		if typedA == typedB {
			return nil
		}

		return fmt.Errorf("%+v does not equal %+v", typedA, typedB)
	case int:
		if typedA == float64(typedB) {
			return nil
		}

		return fmt.Errorf("%+v does not equal %+v", typedA, typedB)
	default:
		return fmt.Errorf("%+v of type %T could not be cast to float64 or int", b, b)
	}
}

func partialDeepCompareInt(typedA int, b interface{}) error {
	switch typedB := b.(type) {
	// Handle both int and float64 types for second data type
	case float64:
		if float64(typedA) == typedB {
			return nil
		}

		return fmt.Errorf("%+v does not equal %+v", typedA, typedB)
	case int:
		if typedA == typedB {
			return nil
		}

		return fmt.Errorf("%+v does not equal %+v", typedA, typedB)
	default:
		return fmt.Errorf("%+v of type %T could not be cast to float64 or int", b, b)
	}
}

func partialDeepCompareString(typedA string, b interface{}) error {
	typedB, ok := b.(string)
	if !ok {
		return fmt.Errorf("%+v of type %T could not be cast to string", b, b)
	}

	if typedA == typedB {
		return nil
	}

	return fmt.Errorf("%+v does not equal %+v", typedA, typedB)
}

func partialDeepCompareSlice(typedA []interface{}, b interface{}) error {
	typedB, ok := b.([]interface{})
	if !ok {
		return fmt.Errorf("%+v of type %T could not be cast to []interface{}", b, b)
	}

	if len(typedA) != len(typedB) {
		return fmt.Errorf("length of %+v does not equal length of %+v", typedA, typedB)
	}

	for i := range typedA {
		result := partialDeepCompare(typedA[i], typedB[i])
		if result != nil {
			return result
		}
	}

	return nil
}

func partialDeepCompareMap(typedA map[string]interface{}, b interface{}) error {
	typedB, ok := b.(map[string]interface{})
	if !ok {
		return fmt.Errorf("%+v of type %T could not be cast to map[string]interface{}", b, b)
	}

	// We are only interested in the keys that are present in the first map
	for key, aValue := range typedA {
		bValue, ok := typedB[key]
		if !ok {
			return fmt.Errorf("key %s is not present in %+v", key, typedB)
		}

		result := partialDeepCompare(aValue, bValue)
		if result != nil {
			return result
		}
	}

	return nil
}
