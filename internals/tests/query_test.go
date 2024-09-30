package tests

import (
	"testing"

	"go.starlark.net/starlark"
	"test.com/excel-ark/internals"
)

//HIGH: review/update logging with logging levels

//TODO: make separate directory (package?) for tests and check how commands will execute there

//Test command line (in root directory): go test -v ./...

//General test information: https://go.dev/doc/tutorial/add-a-test

var (
	query_test_cases = map[string]string{
		"./test_cases/query_test.star":        "2d3e7fc8284ea9805a41fd54cbbf07aec3365b925fbabe90a29769dc7a997109",
		"./test_cases/excel_load_simple.star": "ad211fb93f6948a9111606f217f178ad1f0cc0baf20a2ba3ad9e46d9e2a7a8ad",
		//06May2024 SHA change due to SkipRows worksheet change
		"./test_cases/excel_load_skip_rows.star": "6aa5aaf2e4eb1a575e81c89153119622a54b0e681eea2d4947ef3f275491ab18",
		"./test_cases/excel_load_range.star":     "2d3e7fc8284ea9805a41fd54cbbf07aec3365b925fbabe90a29769dc7a997109",
	}
)

func TestQuery(t *testing.T) {
	const KEY_NAME = "starlark_result_sha"
	for starlark_file, expected_sha := range query_test_cases {
		internals.Init(true, nil)
		t.Logf("Runnin test file:%s\n", starlark_file)
		globals := internals.ExecStarlarkFile(starlark_file)
		ret_val, ok := (*globals)["starlark_result_sha"]
		if !ok {
			t.Errorf("TestQuery: no expected test variable with name %s",
				KEY_NAME)
			return
		}
		if internals.StripDoubleQuotes(ret_val.(starlark.String).String()) != expected_sha {
			t.Errorf("TestQuery: test variable returned from query_test:%s returned wrong value:%s\n\t while expected is:%s",
				KEY_NAME,
				internals.StripDoubleQuotes(ret_val.(starlark.String).String()),
				expected_sha)
			t.Logf("\t file RUN FAIL\n")
		} else {
			t.Logf("\t file RUN OK\n")
		}
		internals.Close()
	} //for starlark_file, expected_sha := range query_test_cases {
} //func TestQuery(t *testing.T) {
