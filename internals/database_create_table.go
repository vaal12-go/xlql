package internals

import (
	"fmt"

	starlark "go.starlark.net/starlark"
)

//[x] File is cleaned
//[x]: review/update logging with logging levels

// HIGH: limit debug messages in production using logging priorities
//
//	https://www.honeybadger.io/blog/golang-logging/
//	https://github.com/Sirupsen/logrus
//
// TODO: add parameter to delete table if exists
func (self Database) Create_table(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {

	var tblName string
	var val starlark.Value
	var auto_rename_table_name = false
	if err := starlark.UnpackArgs(
		b.Name(), args, kwargs,
		"name", &tblName,
		"columns", &val,
		"auto_rename_table_name?", &auto_rename_table_name); err != nil {
		return nil, err
	}
	DLf("create_table: tblName: %v\n", tblName)
	dct := val.(*starlark.Dict)
	colMap := StarlarkDictToMap(dct)
	colNames := make([]string, 0)

	//FIXME: this relies on dct.Keys returns same sequence of columns as declared in starlark
	// e.g. new_tbl_q = newDB.create_table(
	// name = "newTable",
	// columns = {
	//     "Col1": "TEXT",
	//     "Col23": "NUMERIC",
	//     "Col2": "NUMERIC",
	//     "Col0": "NUMERIC"
	// })
	//This can not be reliable. Should revise columns signature of the create_table function

	for _, colName := range dct.Keys() {
		colNames = append(colNames, RemoveQuotesFromString(colName.String()))
	}
	DLf("colNames: %v\n", colNames)
	DLf("colMap: %v\n", colMap)
	tblNameActual, colNamesArray, err := self.createTable(
		tblName, &colNames, &colMap, auto_rename_table_name)
	if err != nil {
		return starlark.None, err
	}
	DLf("tblNameActual: %v\n", tblNameActual)
	DLf("colNamesArray: %v\n", colNamesArray)
	ret := NewQuery(&self,
		fmt.Sprintf("SELECT * FROM %s", tblNameActual),
		tblNameActual)
	return ret, nil
} //func (self Database) create_table(thread *starlark.Thread,
