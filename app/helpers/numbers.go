package helpers

import "strconv"

// Converts a string to an integer. If string in question
// cannot be converted, it will be converted to a default value
func ConvertToInt(value string, defaultValue int) int {
	if len(value) == 0 {
		return defaultValue
	}

	valueInt, err := strconv.Atoi(value)

	if err != nil {
		valueInt = defaultValue
	}

	return valueInt
}
