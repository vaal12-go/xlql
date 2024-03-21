package tests

import (
	"testing"

	"go.starlark.net/starlark"
	"test.com/excel-ark/internals"
)

//HIGH: clean file
//HIGH: review/update logging with logging levels

//TODO: make separate directory (package?) for tests and check how commands will execute there

func prepareDB() {

}

//Test command line (in root directory): go test -v ./...

//General test information: https://go.dev/doc/tutorial/add-a-test

var (
	query_test_cases = map[string]string{
		"load_excel_sheet_table_range": "d3847b68cf21864dd885f26436dd1037cee6b5084139f5ce9c8a31a44d7aea29",
	}
)

func TestQuery(t *testing.T) {
	internals.Init(true)
	globals := internals.ExecStarlarkFile("./test_cases/query_test.star")
	// fmt.Printf("globals: %v\n", globals)

	// t.Log("load_excel_sheet_table_range:",
	// 	internals.StripDoubleQuotes((*globals)["load_excel_sheet_table_range"].(starlark.String).String()))
	for key, expectedVal := range query_test_cases {
		if internals.StripDoubleQuotes((*globals)[key].(starlark.String).String()) != expectedVal {
			t.Errorf("TestQuery: test variable returned from query_test:%s returned wrong value:%s",
				key, expectedVal)
		}
	}
	internals.Close()

}
