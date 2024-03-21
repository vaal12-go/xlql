package internals

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strings"

	"go.starlark.net/starlark"
)

var QUERY_TYPE_METHODS = []string{"Get_column_names", "Print", "Get_cursor",
	"Add_row", "Save_to_excel", "Table_name", "Get_sha"}

type Query struct {
	connected_db             *Database
	query_sql                string
	table_name               string //For cases when query is created from load_excel_sheet
	get_column_names_builtin *starlark.Builtin
	get_column_types_builtin *starlark.Builtin
	get_cursor_builtin       *starlark.Builtin
	print_builtin            *starlark.Builtin
	add_row_builtin          *starlark.Builtin
	save_to_excel_builtin    *starlark.Builtin
	exporter                 MethodExporter
	StarlarkIterableImplmentationStub
}

func (self Query) Get_sha_internal() (string, error) {
	hash := sha256.New() // Use the Hash in crypto/sha256
	// colArr := self.get_col_names_internal()
	curs, err := NewCursorInternal(self.connected_db.db_connection, self.query_sql)
	if err != nil {
		log.Fatalf("query.Get_sha_internal. NewCursorInternal failed with error:%v\n", err)
	}
	colNames := curs.GetColumnNames()
	colTypes := curs.GetColumnTypes()
	for idx, colName := range colNames {
		hash.Write([]byte(colName))
		hash.Write([]byte(colTypes[idx].DatabaseTypeName()))
	}
	for curs.Next() {
		currTblArr, err := curs.GetRow()
		if err != nil {
			log.Fatalf("query.print. GetRow failed with error:%v\n", err)
		}
		for idx := range *currTblArr {
			hash.Write([]byte(getStringValueFromInterface((*currTblArr)[idx])))
		} //for idx := range *currTblArr {
	} //for curs.Next() {

	sha_str := fmt.Sprintf("%x", hash.Sum(nil)) // Get encoded hash sum
	return sha_str, nil
} //func (self Query) getSHA() (string, error) {

func (self Query) Get_sha(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	sha, err := self.Get_sha_internal()
	if err != nil {
		return starlark.None, fmt.Errorf("Error in getSHA:%v\n", err)
	}
	return starlark.String(sha), nil
} //func (self Query) get_cursor(thread *starlark.Thread,

func (self Query) Get_cursor(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	return NewCursor(&self), nil
} //func (self Query) get_cursor(thread *starlark.Thread,

func (self Query) Table_name(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.String(self.table_name), nil
} //func (self Query) Table_name(thread *starlark.Thread,

func NewQuery(db *Database, sql string,
	table_name string) *Query {
	ret := Query{
		query_sql:    sql,
		connected_db: db,
		table_name:   table_name,
	}
	//below var is necessary to get proper pointer to ret
	var iFacePointer interface{} = ret
	ret.exporter.RegisterBuiltIns(&iFacePointer, QUERY_TYPE_METHODS)
	return &ret
} //func NewQuery(db *Database, sql string,

func (self Query) Attr(name string) (starlark.Value, error) {
	switch name {
	// case "table_name":
	// 	//[x]: make this a function call so it cannot be set from starlark
	// 	// fmt.Printf("self.table_name: %v\n", self.table_name)
	// 	return starlark.String(self.table_name), nil
	default:
		return self.exporter.GetMethod(name)
	}
}

func (self Query) AttrNames() []string {
	ret := reduce(QUERY_TYPE_METHODS, func(acc []string, val string) []string {
		return append(acc, strings.ToLower(val))
	}, make([]string, 0))
	return ret
}

func (self Query) Type() string {
	return "Query"
}

func (self Query) String() string {
	tbl_str := ""
	if self.table_name != "" {
		tbl_str = fmt.Sprintf("\nTable: %s", self.table_name)
	}
	return fmt.Sprintf("Query%s\nSQL:%s", tbl_str, self.query_sql)
}
