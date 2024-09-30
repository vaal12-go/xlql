package tests

import (
	"fmt"
	"testing"

	"go.starlark.net/starlark"
	"test.com/excel-ark/internals"
)

// TO run: go test ./internals/tests/date_parse_test.go -v

var (
	date_parse_test_cases = map[string]string{
		"./test_cases/date_parse_no_err.star":       "984d248c0a8b206b88599d27761d6e5e6fb6f05e0aad4a4bd0f5f8213dca624f", //Old one "2f1a271142504609bbd22e8dda94c544d22c5e8c0328413f7d34d613d889706e"
		"./test_cases/date_parse_wrong_format.star": "5822a1d38f0ca8adf4708d569926e7bb401bb9b1380f3aa0a713cfa6a0283ae3", // "414270b0da3add1c35245060a050f19eeb3d849b257f95643f0b35cc0e46e4f5",
	}
)

func RunTestScriptsWithMap(testMap map[string]string, result_key string, t *testing.T) error {
	for starlark_file, expected_sha := range date_parse_test_cases {
		internals.Init(true, nil)
		t.Logf("Runnin test file:%s\n", starlark_file)
		globals := internals.ExecStarlarkFile(starlark_file)
		ret_val, ok := (*globals)["starlark_result_sha"]
		if !ok {
			t.Errorf("TestQuery: no expected test variable with name %s",
				result_key)
			return fmt.Errorf(
				"TestQuery: no expected test variable with name %s",
				result_key)
		}
		if internals.StripDoubleQuotes(ret_val.(starlark.String).String()) != expected_sha {
			t.Errorf("TestQuery: test variable returned from query_test:%s returned wrong value:%s\n\t while expected is:%s",
				result_key,
				internals.StripDoubleQuotes(ret_val.(starlark.String).String()),
				expected_sha)
			t.Logf("\t file [%s] RUN FAIL\n", starlark_file)
		} else {
			t.Logf("\t file [%s] RUN OK\n", starlark_file)
		}
		internals.Close()
	}
	return nil
}

func TestDateParsing(t *testing.T) {
	const KEY_NAME = "starlark_result_sha"
	//TODO: make a function of the test cases runner
	RunTestScriptsWithMap(date_parse_test_cases, KEY_NAME, t)
} //func TestDateParsing(t *testing.T) {
