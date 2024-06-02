package tests

import (
	"testing"

	"go.starlark.net/starlark"
	"test.com/excel-ark/internals"
)

func TestLoadExcelSheetDropTable(t *testing.T) {
	internals.Init(true, nil)
	globals := internals.ExecStarlarkFile("./test_cases/query_load_excel2.star")
	ret_key := internals.StripDoubleQuotes((*globals)["test_result"].(starlark.String).String())
	// fmt.Printf("key from test:%v", ret_key)
	if ret_key != "OK" {
		t.Errorf(
			"TestLoadExcelSheetDropTable failed with message:%v", ret_key)
	}
	internals.Close()
} //func TestLoadExcelSheetDropTable(t *testing.T) {

func TestCreateTableDropTable(t *testing.T) {
	internals.Init(true, nil)
	globals := internals.ExecStarlarkFile("./test_cases/database_create_table.star")
	ret_key := internals.StripDoubleQuotes((*globals)["test_result"].(starlark.String).String())
	// fmt.Printf("key from test:%v", ret_key)
	if ret_key != "OK" {
		t.Errorf(
			"TestLoadExcelSheetDropTable failed with message:%v", ret_key)
	}
	internals.Close()
} //func TestCreateTableDropTable(t *testing.T) {

func TestExecSQL(t *testing.T) {
	internals.Init(true, nil)
	globals := internals.ExecStarlarkFile("./test_cases/db_execSQL.star")
	ret_key := internals.StripDoubleQuotes((*globals)["test_result"].(starlark.String).String())
	// fmt.Printf("key from test:%v", ret_key)
	if ret_key != "OK" {
		t.Errorf(
			"TestExecSQL failed with message:%v", ret_key)
	}
	internals.Close()
} //func TestExecSQL(t *testing.T) {
