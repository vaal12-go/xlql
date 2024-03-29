package internals

import (
	"strings"
	"testing"

	"go.starlark.net/starlark"
)

//TO run: go test -run TestLoadExcelSheet -v

func UnescapeStarlarkString(val starlark.Value) string {
	replacementMap := map[string]string{
		"\\n": "\n",
		"\\t": "\t",
		"\"":  "",
	}
	val_str := val.String()
	for k, v := range replacementMap {
		val_str = strings.ReplaceAll(val_str, k, v)
	}
	return val_str
}

var TestGetParametersMap = map[string]string{
	"test_result1": `loadExceSheetParams
		sheet_name: test1
		skip_rows: 0
		file_name: qwe1.xlsx
		use_cols: []
		db_col_names : map[]
		append_rows: false
		db_table_name : 
		drop_table: false
		table_range_start: 
		table_range_start_x 0
		table_range_start_y 0
		table_range_end: 
		table_range_end_x   0
		table_range_end_y   0
		auto_rename_table_name false
		column_parameters_dict map[]`,
	"test_result2": `loadExceSheetParams
		sheet_name: test1
		skip_rows: 2
		file_name: qwe1.xlsx
		use_cols: []
		db_col_names : map[]
		append_rows: false
		db_table_name : 
		drop_table: false
		table_range_start: 
		table_range_start_x 0
		table_range_start_y 0
		table_range_end: 
		table_range_end_x   0
		table_range_end_y   0
		auto_rename_table_name false
		column_parameters_dict map[]`,
}

func TestGetParameters(t *testing.T) {
	Init(true, nil)
	PredeclaredDict["TEST_getParametersLoad_excel_sheet"] =
		starlark.NewBuiltin("TEST_getParametersLoad_excel_sheet",
			MOCK_getParametersLoad_excel_sheet)
	globals := ExecStarlarkFile("./tests/test_cases/query_load_excel.star")
	// fmt.Printf("globals: %v\n", globals)
	for k, v := range TestGetParametersMap {
		result := UnescapeStarlarkString((*globals)[k])
		if result != v {
			DLf("Expected Result:")
			DLf(v)
			DLf("returnedResult:")
			DLf(result)
			t.Errorf("Parameter:%s did not return expected result:%s\n But instead returned:%s\n",
				k, v, result)
		}
	}
} //func TestGetParameters(t *testing.T) {

func TestLoadExcelSheet(t *testing.T) {

}
