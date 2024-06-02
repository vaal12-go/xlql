package internals

import (
	"log"
	"strconv"
	"strings"
	"time"

	"go.starlark.net/starlark"
)

//[x]: review/update logging with logging levels

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	DL            *log.Logger
	DLf           (func(format string, v ...any))
)

func reduce[T, M any](s []T, f func(M, T) M, initValue M) M {
	acc := initValue
	for _, v := range s {
		acc = f(acc, v)
	}
	return acc
}

//TODO: copy ErrorLogger to stdout
//TODO: if -v is provided in the command line of the programm - copy Info and Warning loggers to stdout
//TODO: copy DebugLogger (DL) to stdout during development

func getStringValueFromInterface(val any) string {
	// fmt.Printf("getStringValueFromInterface. val: %v\n", val)
	// fmt.Printf("reflect.TypeOf(val): %v\n", reflect.TypeOf(val))
	switch val.(type) {
	case string:
		return val.(string)
	case int:
		return strconv.Itoa(val.(int))
	case int64:
		return strconv.FormatInt(val.(int64), 10)
	case nil:
		return ""
	case float64:
		return strconv.FormatFloat(val.(float64), 'f', -1, 64)
	case time.Time:
		return val.(time.Time).Format("02-Jan-2006")
	default:
		log.Fatalf("getStringValueFromInterface - unknown type: %v/n %t/n",
			val, val)
	}
	return ""
} //func getStringValueFromInterface(val any) string {

func CopyToSliceOfGenerics(from []interface{}) *[]interface{} {
	retSlice := make([]interface{}, len(from))
	copy(retSlice, from)
	return &retSlice
}

func getSQLTypeFromExcelType(xl_cell_type int) string {
	switch xl_cell_type {
	case 2:
		return "DATE"
	case 6:
		return "NUMERIC"
	case 7:
		return "TEXT"
	default:
		return "TEXT"
	} //switch cellType {
} //func getSQLTypeFromExcelType(xl_cell_type int) string {

// From here: https://stackoverflow.com/questions/27689058/convert-string-to-interface/73034622#73034622
func toAnyList[T any](input []T) []any {
	list := make([]any, len(input))
	for i, v := range input {
		list[i] = v
	}
	return list
}

func StarlarkDictToMap(dct *starlark.Dict) map[string]string {
	//TODO: check where this one is used and remove conversion to string of the map values
	retMap := make(map[string]string, 0)
	for _, key := range dct.Keys() {
		val, _, _ := dct.Get(key)
		retMap[RemoveQuotesFromString(key.String())] = RemoveQuotesFromString(val.String())
	}
	return retMap
}

func StarlarkDictToMap2(dct *starlark.Dict) map[string]any {
	retMap := make(map[string]any, 0)
	for _, key := range dct.Keys() {
		val, _, _ := dct.Get(key)
		retMap[RemoveQuotesFromString(key.String())] = val
	}
	return retMap
}

func StripDoubleQuotes(str string) string {
	start, end := 0, len(str)
	if end == 0 {
		return ""
	}
	if str[0] == '"' {
		start++
	}
	if str[end-1] == '"' {
		end--
	}
	if start > end {
		return ""
	}
	return str[start:end]
}

func RemoveQuotesFromString(instr string) string {
	return strings.ReplaceAll(instr, "\"", "")
}
