package assertions

func partialDeepCompare(a, b interface{}) bool {
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
		return false
	}
}

func partialDeepCompareBool(typedA bool, b interface{}) bool {
	typedB, ok := b.(bool)
	if !ok {
		return false
	}

	return typedA == typedB
}

func partialDeepCompareFloat64(typedA float64, b interface{}) bool {
	switch typedB := b.(type) {
	// Handle both int and float64 types for second data type
	case float64:
		return typedA == typedB
	case int:
		return typedA == float64(typedB)
	default:
		return false
	}
}

func partialDeepCompareInt(typedA int, b interface{}) bool {
	switch typedB := b.(type) {
	// Handle both int and float64 types for second data type
	case float64:
		return float64(typedA) == typedB
	case int:
		return typedA == typedB
	default:
		return false
	}
}

func partialDeepCompareString(typedA string, b interface{}) bool {
	typedB, ok := b.(string)
	if !ok {
		return false
	}

	return typedA == typedB
}

func partialDeepCompareSlice(typedA []interface{}, b interface{}) bool {
	typedB, ok := b.([]interface{})
	if !ok {
		return false
	}

	if len(typedA) != len(typedB) {
		return false
	}

	for i := range typedA {
		if !partialDeepCompare(typedA[i], typedB[i]) {
			return false
		}
	}

	return true
}

func partialDeepCompareMap(typedA map[string]interface{}, b interface{}) bool {
	typedB, ok := b.(map[string]interface{})
	if !ok {
		return false
	}

	// We are only interested in the keys that are present in the first map
	for key, aValue := range typedA {
		bValue, ok := typedB[key]
		if !ok {
			return false
		}

		if !partialDeepCompare(aValue, bValue) {
			return false
		}
	}

	return true
}
