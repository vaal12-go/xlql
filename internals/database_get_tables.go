package internals

import (
	"log"

	"go.starlark.net/starlark"
)

func (self Database) get_tables_actual() []string {
	retArr := make([]string, 0)
	sql := `
	    SELECT 
			name
		FROM 
			sqlite_schema
		WHERE 
			type ='table' AND 
			name NOT LIKE 'sqlite_%';`
	curs, err := NewCursorInternal(self.db_connection, sql)
	if err != nil {
		log.Fatalf("get_tables_actual. NewCursorInternal failed with error:%v\n", err)
	}
	for curs.Next() {
		// currTblArr, err := curs.GetRow()
		currTblArr, err := curs.GetRow()
		if err != nil {
			log.Fatalf("get_tables_actual. GetRow failed with error:%v\n", err)
		}
		retArr = append(retArr, ((*currTblArr)[0].(string)))
	}
	return retArr
} //func (self Database) get_tables_actual() []string {

func (self Database) Get_tables(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	retArr := self.get_tables_actual()
	valArr := make([]starlark.Value, 0)
	for _, st := range retArr {
		valArr = append(valArr, starlark.String(st))
	}
	return starlark.NewList(valArr), nil
} //func (self Database) get_tables(thread *starlark.Thread,
