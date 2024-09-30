package internals

import (
	"crypto/sha256"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/xuri/excelize/v2"
	"go.starlark.net/starlark"
)

type columnParameter struct {
	Type   string
	Format string
}

func NewColumnParameter(func_param map[string]any) *columnParameter {
	ret := columnParameter{}

	tp, ok := func_param["type"]
	if ok {
		// fmt.Printf("reflect.TypeOf(tp): %v\n", reflect.TypeOf(tp))
		ret.Type = StripDoubleQuotes(tp.(starlark.String).String())
	}

	fm, ok := func_param["format"]
	if ok {
		// fmt.Printf("reflect.TypeOf(fm): %v\n", reflect.TypeOf(fm))
		ret.Format = StripDoubleQuotes(fm.(starlark.String).String())
	}
	return &ret
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
	column_parameters_dict map[string]*columnParameter
	calc_header_start_row  int64
	calc_header_start_col  int64
	calc_header_end_col    int64
	calc_data_start_row    int64
	calc_data_end_row      int64
	calc_data_start_col    int64
}

func (esp loadExceSheetParams) SHA256() string {
	hash := sha256.New() // Use the Hash in crypto/sha256
	hash.Write([]byte(esp.String()))
	sha_str := fmt.Sprintf("%x", hash.Sum(nil)) // Get encoded hash sum
	return sha_str
	// fmt.Printf("loadExceSheetParams sha_str: %v\n", sha_str)
}

// [x]: Rework this using reflection
func (self loadExceSheetParams) String() string {
	tp := reflect.TypeOf(self)
	VALUE_START_POS := 25

	ret := ""
	for i := 0; i < (tp).NumField(); i++ {
		field := (tp).Field(i)
		fieldVal := reflect.ValueOf(self).FieldByIndex([]int{i})
		valStr := fmt.Sprintf("%v", fieldVal)
		padSpacesNum := VALUE_START_POS - len(field.Name) - 1 + (10 - len(valStr)) //-1 to account for colon
		if len(valStr) > 10 {
			padSpacesNum = VALUE_START_POS - len(field.Name) - 1 //-1 to account for colon
		}
		if field.Name == "excelizeFile" && len(valStr) > 10 {
			valStr = valStr[:10]
		}
		outStr := fmt.Sprintf("  %s%s%v\n",
			field.Name,
			strings.Repeat(".", padSpacesNum),
			valStr)

		ret = fmt.Sprintf("%s%s", ret, outStr)
		// fmt.Println(
		// 	strings.Repeat("-", 40))
		// fmt.Println("\ttype:", field.Type)
		// fmt.Println("\tanon:", field.Anonymous)
		// fmt.Printf("field.Type.Kind(): %v\n", field.Type.Kind())
		// fmt.Printf("\t Value: %v\n",
		// 	reflect.ValueOf(params).FieldByIndex([]int{i}))
	} //for i := 0; i < (tp).NumField(); i++ {
	return ret
} //func (self loadExceSheetParams) String() string {

//TODO: make loadExceSheetParams fields same name as load_excel_sheet parameter names

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
		column_parameters_dict: make(map[string](*columnParameter)),
		calc_header_start_row:  -1,
		calc_header_start_col:  -1,
		calc_data_start_row:    -1,
		calc_data_start_col:    -1,
	}

	// fmt.Printf("args: %v\n", args)

	//TODO: add all parameters to unnamed parsing
	for i := 0; i < args.Len(); i++ { //Parsing unnamed parameters
		switch i {
		case 0: //file_name
			val := args.Index(0)
			// fmt.Printf("reflect.TypeOf(val): %v\n", reflect.TypeOf(val))
			if reflect.TypeOf(val).String() != "starlark.String" {
				return nil,
					fmt.Errorf("getParameters. file_name parameter is not string.")
			}
			params.excel_file_name = StripDoubleQuotes(args.Index(0).String())
		case 1: //sheet_name
			val := args.Index(1)
			// fmt.Printf("reflect.TypeOf(val): %v\n", reflect.TypeOf(val))
			if reflect.TypeOf(val).String() != "starlark.String" {
				return nil,
					fmt.Errorf("getParameters. sheet_name parameter is not string.")
			}
			params.sheet_name = StripDoubleQuotes(args.Index(1).String())
		}
	} //for i := 0; i < args.Len(); i++ {//Parsing unnamed parameters

	// fmt.Printf("params After args parsing: %v\n", params)

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
			params.excel_file_name = StripDoubleQuotes(arg[1].String())
		case "table_range":
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
			// fmt.Printf("arg[1]: %v\n", arg[1])
			colMap := arg[1].(*starlark.Dict)
			// fmt.Printf("colMap.Keys(): %v\n", colMap.Keys())
			colRealMap := StarlarkDictToMap2(colMap)
			// fmt.Printf("colRealMap: %#v\n", colRealMap)
			for key, val := range colRealMap {
				// fmt.Printf("key: %v\n", key)
				retVal := val
				if reflect.TypeOf(val).String() == "*starlark.Dict" {
					retVal = StarlarkDictToMap2(val.(*starlark.Dict))
				} else {
					return nil, fmt.Errorf("getParameters. Parsing cols parameter: non Dict 2nd level value. All values of 'cols' parameter dict should also be dicts.")
				}
				colParam := NewColumnParameter(retVal.(map[string]any))
				// fmt.Printf("colParam: %#v\n", colParam)
				params.column_parameters_dict[key] = colParam

				// fmt.Printf("col: %v\n", StripDoubleQuotes(col.String()))
				// params.column_parameters_dict[StripDoubleQuotes(col.String())] = columnParameters{
				// 	date_format: "NONE YET",
				// }
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
