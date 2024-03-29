package internals

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
)

func (self Query) Add_row(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	//HIGH: review checks for arguments and parameters and move those below after parsing of kwargs
	if self.table_name == "" {
		return starlark.None,
			fmt.Errorf("query.add_row. This query \n'%s'\n is not a table. Rows can only be added to actual table in the database.", self.query_sql)
	}
	colNamesValuesMap := make(map[string]string, 0)
	colNamesArr := self.get_col_names_internal()
	valArr := make([]string, 0)
	//TODO: when number of arguments is less than number of columns - populate the rest with nulls
	for i := 0; i < args.Len(); i++ {
		valArr = append(valArr,
			fmt.Sprintf("%s", RemoveQuotesFromString(args.Index(i).String())))
		colNamesValuesMap[colNamesArr[i]] = args.Index(i).String()
	}
	if len(valArr) < len(colNamesArr) && len(valArr) != 0 {
		return starlark.None,
			fmt.Errorf("query.add_row. Number of values passed [%d] is less than number of columns [%d]",
				len(valArr), len(colNamesArr))
	}

	//TODO: accept first non-named argument as 'values'
	for _, arg_tuple := range kwargs {
		arg_name := strings.ReplaceAll(arg_tuple[0].String(), "\"", "")
		if len(colNamesValuesMap) > 0 {
			return starlark.None,
				fmt.Errorf("query.add_row keyword argument [%s] after list of values.", arg_name)
		}
		if arg_name == "values" {
			// DLf("arg_name: %v\n", arg_name)
			// DLf("arg type: %T\n", arg_tuple[1])
			switch arg_tuple[1].(type) {
			case *starlark.List:
				val_list := arg_tuple[1].(*starlark.List)
				if len(colNamesArr) < val_list.Len() {
					return starlark.None,
						fmt.Errorf("add_row: number of values provided in 'values' argument is larger than number of columns in the table.")
				}
				for i := 0; i < val_list.Len(); i++ {
					// DLf("val_list.Index(i): %v\n", val_list.Index(i))
					colNamesValuesMap[colNamesArr[i]] = RemoveQuotesFromString(val_list.Index(i).String())
				}
				//TODO: check if None value will result in null DB row
			case *starlark.Dict: //TODO: implement dict argument in add_row
				val_list := arg_tuple[1].(*starlark.Dict)
				DLf("Dict val_list: %v\n", val_list)
				return starlark.None,
					fmt.Errorf("query.add_row dict values for values is not implemented")
			default:
				return starlark.None,
					fmt.Errorf("query.add_row Unsupported values argument type. Need List or Dictionary:%s", arg_name)
			} //switch arg_tuple[1].(type) {
		} else { //if arg_name == "values" {
			return starlark.None, fmt.Errorf("Argument name '%s' is not supported.", arg_name)
		}
	} //for _, arg_tuple := range kwargs {
	if len(colNamesValuesMap) == 0 {
		return starlark.None, nil
	}

	insertSQL := prepareInsertStatementFromArray2(self.table_name, colNamesArr)
	anyValArr := toAnyList(valArr)

	// fmt.Printf("self.connected_db.db_connection.Stats(): %#v\n", self.connected_db.db_connection.Stats())

	var err error
	// tx, err := self.connected_db.db_connection.Begin()
	// if err != nil {
	// 	return starlark.None, err
	// }
	_, err = self.connected_db.db_connection.Exec(insertSQL, anyValArr...)
	if err != nil {
		return starlark.None,
			fmt.Errorf("query.add_row Error executing INSERT statement:%s", err)
	}
	// tx.Commit()
	return starlark.None, nil
} //func (self Query) add_row(thread *starlark.Thread,
