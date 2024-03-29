package internals

import (
	"fmt"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
	"go.starlark.net/starlark"
)

type columnParameters struct {
	date_format string
}

type loadExceSheetParams struct {
	excel_file_name        string
	sheet_name             string
	skip_rows              int64
	use_cols               []string
	db_col_names           map[string]string
	append_rows            bool
	db_table_name          string
	drop_table             bool
	excelizeFile           *excelize.File
	table_range_start      string
	table_range_start_x    int64
	table_range_start_y    int64
	table_range_end        string
	table_range_end_x      int64
	table_range_end_y      int64
	auto_rename_table_name bool
	column_parameters_dict map[string]columnParameters
}

// TODO: Rework this using reflection
func (self loadExceSheetParams) String() string {
	return fmt.Sprintf(
		`loadExceSheetParams
		sheet_name: %s
		skip_rows: %d
		file_name: %s
		use_cols: %v
		db_col_names : %v
		append_rows: %v
		db_table_name : %s
		drop_table: %v
		table_range_start: %s
		table_range_start_x %d
		table_range_start_y %d
		table_range_end: %s
		table_range_end_x   %d
		table_range_end_y   %d
		auto_rename_table_name %v
		column_parameters_dict %v`,

		self.sheet_name,
		self.skip_rows,
		self.excel_file_name,
		self.use_cols,
		self.db_col_names,
		self.append_rows,
		self.db_table_name,
		self.drop_table,
		self.table_range_start,
		self.table_range_start_x,
		self.table_range_start_y,
		self.table_range_end,
		self.table_range_end_x,
		self.table_range_end_y,
		self.auto_rename_table_name,
		self.column_parameters_dict)
} //func (self loadExceSheetParams) String() string {

func getParameters(args starlark.Tuple,
	kwargs []starlark.Tuple) (*loadExceSheetParams, error) {
	params := loadExceSheetParams{
		excel_file_name:        "",
		sheet_name:             "",
		skip_rows:              0,
		use_cols:               make([]string, 0),
		db_col_names:           make(map[string]string, 0),
		append_rows:            false,
		db_table_name:          "",
		drop_table:             false,
		table_range_start:      "",
		table_range_end:        "",
		auto_rename_table_name: false,
		column_parameters_dict: make(map[string]columnParameters),
	}

	//Below will not work as it checks for passed keyword arguments and if cannot find it panics
	// if err := starlark.UnpackArgs(
	// 	b.Name(), args, kwargs,
	// 	"auto_rename_table_name?",
	// 	&params.auto_rename_table_name); err != nil {
	// 	return nil, err
	// }
	// DLf("Hello from getParameters")
	//TODO: move parameters parsing to separate function
	for _, arg := range kwargs {
		// For some reason names of arguments are within "" - e.g. "worksheet" - double quotes are also part of the string
		argName := StripDoubleQuotes(arg[0].String())
		// DLf("Stripped argName:%s", argName)
		switch argName {
		case "drop_table":
			params.drop_table = bool(arg[1].(starlark.Bool).Truth())
		case "auto_rename_table_name":
			params.auto_rename_table_name = bool(arg[1].(starlark.Bool).Truth())
		case "sheet_name":
			params.sheet_name = StripDoubleQuotes(arg[1].String())
		case "skip_rows":
			var error_occurred bool
			params.skip_rows, error_occurred = (arg[1]).(starlark.Int).Int64()
			if error_occurred != true {
				log.Fatal("load_excel_sheet - Error parsing skip_rows")
			}
		case "usecols":
			log.Fatal("usecols parameter of load_excel_sheet is not implemented")
		case "names":
			log.Fatal("names parameter of load_excel_sheet is not implemented")
		case "file_name":
			params.excel_file_name = strings.ReplaceAll(arg[1].String(), "\"", "")
		case "table_range":
			// DLf("StripDoubleQuotes(arg[1].String()): %v\n", StripDoubleQuotes(arg[1].String()))
			range_arr := strings.Split(StripDoubleQuotes(arg[1].String()), ":")
			params.table_range_start = range_arr[0]
			params.table_range_end = range_arr[1]
			i, k, err := excelize.CellNameToCoordinates(params.table_range_start)
			params.table_range_start_x, params.table_range_start_y =
				int64(i), int64(k)
			if err != nil {
				log.Fatal("load_excel_sheet. Error parsing table_range. " + err.Error())
			}
			i, k, err =
				excelize.CellNameToCoordinates(params.table_range_end)
			params.table_range_end_x, params.table_range_end_y =
				int64(i), int64(k)
			if err != nil {
				log.Fatal("load_excel_sheet. Error parsing table_range. " + err.Error())
			}
		//case "table_range":
		case "cols":
			// DLf("Have cols parameter: %v\n", arg[1])
			// DLf("Have cols parameter: %t\n", arg[1])
			colMap := arg[1].(*starlark.Dict)
			// fmt.Printf("colMap.Keys(): %v\n", colMap.Keys())
			for _, col := range colMap.Keys() {
				// fmt.Printf("col: %v\n", StripDoubleQuotes(col.String()))
				params.column_parameters_dict[StripDoubleQuotes(col.String())] = columnParameters{
					date_format: "NONE YET",
				}
			}
		//case "cols":
		default:
			return nil, fmt.Errorf("Unknown argument of load_excel_worksheet:%v\n", arg[1])
		} //switch argName {
	} //for _, arg := range kwargs {
	if params.excel_file_name == "" {
		log.Fatal("load_excel_sheet: required parameter 'file_name' (Excel file name) is not provided.")
	}
	if params.sheet_name == "" {
		log.Fatal("load_excel_sheet: required parameter 'sheet_name' (Excel worksheet name) is not provided.")
	}
	// fmt.Printf("Load_excel_sheet Params: %v\n", params)
	// DLf("*********** ------------ ***************\n\n")
	return &params, nil
} //func getParameters(args starlark.Tuple,
