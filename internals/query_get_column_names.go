package internals

import (
	"log"

	"go.starlark.net/starlark"
)

// [x]: clean file
// [x]: review/update logging with logging levels
func (self Query) get_col_names_internal() []string {
	rows, err := self.connected_db.db_connection.Query(self.query_sql)
	defer rows.Close()
	if err != nil {
		log.Fatalf("Rows not received. with error: %v\n", err)
	}
	colNames, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	retArr := make([]string, len(colNames))
	copy(retArr, colNames)
	// for _, colName := range colNames {
	// 	retArr = append(retArr, colName)
	// }
	return retArr
} //func (self Query) get_col_names_internal() []string {

func (self Query) Get_column_names(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	strArr := self.get_col_names_internal()
	retArr := make([]starlark.Value, 0)
	for _, col := range strArr {
		retArr = append(retArr, starlark.String(col))
	}
	retList := starlark.NewList(retArr)
	return retList, nil
} //func (self Query) get_column_names(thread *starlark.Thread,
