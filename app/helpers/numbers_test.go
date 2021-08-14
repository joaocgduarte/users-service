package helpers

import (
	"fmt"
	"testing"
)

type ConvertToIntTest struct {
	inputString       string
	inputDefaultValue int
	expectedResult    int
}

var tests = []ConvertToIntTest{
	{"2", 0, 2},
	{"uncastable string", 4, 4},
	{"1001", 3, 1001},
	{"45t", 4, 4},
	{"", 3, 3},
}

// Executes `ConvertToInt` on each of the test cases.
// in failure reports
func TestConvertToInt(t *testing.T) {
	for _, test := range tests {
		result := ConvertToInt(test.inputString, test.inputDefaultValue)

		if result != test.expectedResult {
			errorMessage := fmt.Sprintf("Expected string %s to be converted to %d. Result: %d", test.inputString, test.expectedResult, result)
			t.Error(errorMessage)
		}
	}
}
