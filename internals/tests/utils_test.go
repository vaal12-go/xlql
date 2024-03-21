package tests

import (
	"testing"

	"test.com/excel-ark/internals"
)

var (
	stripDoubleQuotesTestCases = map[string]string{
		"":         "",
		"\"qwe1\"": "qwe1",
		"q\"we1":   "q\"we1",
		"\"\"":     "",
		"qwe1\"":   "qwe1",
		"\"qwe1":   "qwe1",
	}
)

func TestStripDoubleQuotes(t *testing.T) {
	t.Log("TestStripDoubleQuotes starts")
	for key, expectedResult := range stripDoubleQuotesTestCases {
		if internals.StripDoubleQuotes(key) != expectedResult {
			resReceived := internals.StripDoubleQuotes(key)
			t.Errorf(
				"StripDoubleQuotes with argument %s had not returned expected %s \n instead returned: %s",
				key, expectedResult, resReceived)
		}
	}

	// t.Error("I am a error")
}
