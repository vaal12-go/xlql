package internals

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/xuri/excelize/v2"
	"go.starlark.net/starlark"
)

// HIGH: load_excel_sheet - file_name and sheet_name should be used without param names
// TODO: load_excel_sheet - sheet_name if not indicated - first sheet of the file should be used

func (self Database) Load_excel_sheet(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {

	params, err := getParameters(args, kwargs)
	if err != nil {
		ErrorLogger.Printf("Error getting parameters for load_excel_sheet: %v\n", err)
		return starlark.None, err
	}
	params.excelizeFile, err = excelize.OpenFile(params.excel_file_name)
	if err != nil {
		ErrorLogger.Printf("Error opening excel file load_excel_sheet: %v\n", err)
		return starlark.None, err
	}
	defer func() {
		if err := params.excelizeFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	xlColNames, xlColTypes, err := getColumnNamesAndTypes(params)
	if err != nil {
		log.Fatal(err)
	}
	//TODO: check if col names are SQL compatible and fail
	//    https://pkg.go.dev/unicode#IsLetter
	//TODO: check if colNames are unique and if not - fail
	//TODO: add col_auto_rename parameter, so when col name is not SQL compatible or not unique
	//     those are autorenamed
	//TODO: remove leading and trailing spaces from col names and print a message about it
	tx, err := self.db_connection.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { //TODO: maybe combine with close excel defer function
		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var sqlColNames []string
	var stmt *sql.Stmt
	sqlTblName := ""
	colNamesTypesMap := make(map[string]string, 0)
	for idx, colName := range xlColNames {
		colNamesTypesMap[colName] = xlColTypes[idx]
	}
	//TODO: review create table to accept corrected Table name and column types
	sqlTblName, sqlColNames, err = self.createTable(
		params.sheet_name, &xlColNames, &colNamesTypesMap, params.auto_rename_table_name)
	if err != nil {
		return starlark.None, err
	}
	stmt = prepareInsertStatement(tx, sqlTblName, sqlColNames)
	defer stmt.Close()
	rows, err := params.excelizeFile.Rows(params.sheet_name)
	if err != nil {
		log.Fatal(err)
	}
	valStartX, colNamesStartY, _ := getColumnNamesCoords(params)
	valStartY := colNamesStartY + 1
	ret := NewQuery(&self,
		fmt.Sprintf("SELECT * FROM %s", sqlTblName),
		sqlTblName)
	temp_cursor, err := NewCursorInternal(self.db_connection, ret.query_sql)
	sqlColTypes := temp_cursor.GetColumnTypes()
	temp_cursor.Close()
	rowsInserted := 0
	for rowNo := 1; rows.Next(); rowNo++ {
		//TODO: move processing of row to SQL to new function
		if rowNo >= int(valStartY) { //we are in the value row
			row, err := rows.Columns()
			if err != nil {
				log.Fatal(err)
			}
			anySlice := make([]any, 0)
			// Why copying below should be done: https://blog.merovius.de/posts/2018-06-03-why-doesnt-go-have-variance-in/
			tableColNo := 0
			for colNo, colCell := range row { //Iterating all columns
				if tableColNo >= len(sqlColTypes) {
					continue
				}
				if (colNo + 1) >= int(valStartX) { //We are in actual table - need to insert to DB
					//TODO: move cell type handling to separate function
					//TODO: move type conversion/handling to separate function
					if len(colCell) == 0 { //Empty cell - NULL value
						anySlice = append(anySlice, nil)
					} else { //Non empty cell - processing further
						switch sqlColTypes[tableColNo].DatabaseTypeName() {
						case "TEXT":
							anySlice = append(anySlice, colCell)
						//case "TEXT":
						case "NUMERIC":
							//TODO: Add code if cell starts with zero, then it can only be string (or date). Leading zeros should mean string type in Excel
							if colCell[0] == '0' {
								anySlice = append(anySlice,
									fmt.Sprintf("'%s'", colCell))
								break
							}
							intVal, err := strconv.Atoi(colCell)
							if err == nil { //Can convert to INT
								anySlice = append(anySlice, intVal)
							} else { //Cannot convert to INT
								floatVal, err := strconv.ParseFloat(colCell, 64)
								if err == nil { //Can convert to Float
									anySlice = append(anySlice, floatVal)
								} else { //Cannot convert to float
									//TODO: add attempt to convert to date
									anySlice = append(anySlice, colCell)
									//Query which inserts date
									// INSERT INTO LessSimpleSheet_4
									// 	(Date)
									// 	VALUES
									// 	(date('2023-05-23'))
								} //if err == nil {//Can convert to Float
							} //if err == nil {//Can convert to INT
						//case "NUMERIC":
						default:
							anySlice = append(anySlice, colCell)
						} //switch sqlColTypes[tableColNo].DatabaseTypeName() {
					} //if len(colCell)==0 {//Empty cell - NULL value
				} //if (colNo + 1) >= int(valStartX) {//We are in actual table - need to insert to DB
				tableColNo++
			} //for colNo, colCell := range row {//Iterating all columns
			for len(anySlice) < len(sqlColNames) {
				anySlice = append(anySlice, nil)
			}
			_, err = stmt.Exec(anySlice...)
			if err != nil {
				log.Fatal(err)
			}
			rowsInserted++
		} //if rowNo >= int(valStartY) { //we are in the value row
	} //for rowNo := 1; rows.Next(); rowNo++ {
	return ret, nil
} //func (self *Database) load_excel_sheet(thread *starlark.Thread,
