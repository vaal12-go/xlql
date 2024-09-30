package internals

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
	"go.starlark.net/starlark"
)

// [x]: load_excel_sheet - file_name and sheet_name should be used without param names
// TODO: load_excel_sheet - sheet_name if not indicated - first sheet of the file should be used

type XLsqlColumn struct {
	xlColName    string
	xlColType    string
	sqlColType   string
	sqlColName   string
	colParameter *columnParameter
}

func printColArr(colArr *[]*XLsqlColumn) {
	for _, column := range *colArr {
		fmt.Printf("\t column: %#v\n", column)
	}
}

// HIGH: add test routines for datetime parsing
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
	err = calcExcelLoadableRange(params)
	if err != nil {
		return starlark.None, err
	}
	defer func() {
		if err := params.excelizeFile.Close(); err != nil {
			log.Fatalf("Error closing Excel file:%v", err)
		}
	}()

	if params.drop_table {
		err = self.dropTableIfExists(params.sheet_name)
		if err != nil {
			// fmt.Printf("Load_excel_shee. err: %v\n", err)
			return nil, err
		}
	}

	colArray, err := getColumnNamesAndTypes(params)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("colArray: %v\n", colArray)
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
		// fmt.Printf("\"load_excel committed\": %v\n", "load_excel committed")
		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var stmt *sql.Stmt
	sqlTblName := ""
	xlColNames := make([]string, 0)
	//TODO: replace xlColNames with colArray
	for _, column := range *colArray {
		xlColNames = append(xlColNames, column.xlColName)
	}

	//TODO: review create table to accept corrected Table name and column types
	sqlTblName, err = self.createTable(
		params.sheet_name, colArray, params.auto_rename_table_name)
	if err != nil {
		return starlark.None, err
	}

	stmt = prepareInsertStatement(tx, sqlTblName, colArray)
	defer stmt.Close()

	ret, err := NewQuery(&self,
		fmt.Sprintf("SELECT * FROM %s", sqlTblName),
		sqlTblName)

	if err != nil {
		return starlark.None, err
	}
	err = self.IterateDataRows(params, colArray, stmt, ret.query_sql)
	if err != nil {
		return starlark.None, err
	}
	return ret, nil
} //func (self *Database) load_excel_sheet(thread *starlark.Thread,

func (self Database) IterateDataRows(params *loadExceSheetParams,
	columnsArr *[]*XLsqlColumn,
	insertStatement *sql.Stmt,
	query_sql string) error {

	if params.calc_data_start_row < 0 { //No data rows to iterate
		return nil
	}
	rows, err := params.excelizeFile.Rows(params.sheet_name)
	if err != nil {
		return err
	}
	temp_cursor, err := NewCursorInternal(self.db_connection, query_sql)
	sqlColTypes := temp_cursor.GetColumnTypes()
	temp_cursor.Close()
	//TODO: check if getColumnNamesCoords is ever used
	for rowNo := 0; rows.Next(); rowNo++ {
		if (rowNo + 1) < int(params.calc_data_start_row) {
			continue
		}
		if int64(rowNo) > params.calc_data_end_row {
			break
		}
		//TODO: move processing of row to SQL to new function
		row, err := rows.Columns()
		if err != nil {
			log.Fatal(err)
		}
		anySlice := make([]any, 0)
		// Why copying below should be done: https://blog.merovius.de/posts/2018-06-03-why-doesnt-go-have-variance-in/
		tableColNo := 0
		for colNo, colCell := range row { //Iterating all columns
			if (colNo+1)-int(params.calc_data_start_col) >= len(sqlColTypes) {
				continue
			}
			if (colNo + 1) >= int(params.calc_data_start_col) { //We are in actual table - need to insert to DB
				//TODO: move cell type handling to separate function
				//TODO: move type conversion/handling to separate function

				if len(colCell) == 0 { //Empty cell - NULL value
					anySlice = append(anySlice, nil)
				} else { //Non empty cell - processing further
					// fmt.Printf("sqlColTypes[tableColNo].DatabaseTypeName(): %v\n", sqlColTypes[tableColNo].DatabaseTypeName())
					// fmt.Printf("colCell: %v\n", colCell)
					switch sqlColTypes[tableColNo].DatabaseTypeName() {
					case "TEXT":
						anySlice = append(anySlice, colCell)
					//case "TEXT":

					case "NUMERIC":
						//TODO: Add code if cell starts with zero, then it can only be string (or date). Leading zeros should mean string type in Excel
						if colCell[0] == '0' && len(colCell) > 1 {
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
								// (Date) 	VALUES  (date('2023-05-23'))
							} //if err == nil {//Can convert to Float
						} //if err == nil {//Can convert to INT
					//case "NUMERIC":

					case "DATE":
						// fmt.Printf("\"Have date. Will convert\": %v\n", "Have date. Will convert")
						// fmt.Printf("(*columnsArr)[colNo]: %v\n", (*columnsArr)[colNo])
						// fmt.Printf("(*columnsArr)[colNo].colParameter.Format: %v\n", (*columnsArr)[colNo].colParameter.Format)
						tm, err := time.Parse(
							(*columnsArr)[colNo].colParameter.Format,
							colCell)
						if err != nil {
							errStr := fmt.Sprintf("Cannot parse date:%v with format:%s\n",
								colCell, (*columnsArr)[colNo].colParameter.Format)
							fmt.Print(errStr)
							anySlice = append(anySlice, errStr)
						} else {
							anySlice = append(anySlice, tm)
						}
					//case "DATE":

					default:
						anySlice = append(anySlice, colCell)
					} //switch sqlColTypes[tableColNo].DatabaseTypeName() {
				} //if len(colCell)==0 {//Empty cell - NULL value
				tableColNo++
			} //if (colNo + 1) >= int(valStartX) {//We are in actual table - need to insert to DB

		} //for colNo, colCell := range row {//Iterating all columns
		for len(anySlice) < len(*columnsArr) {
			anySlice = append(anySlice, nil)
		}
		_, err = insertStatement.Exec(anySlice...)
		if err != nil {
			return err
		}
		// } //if rowNo >= int(valStartY) { //we are in the value row
	} //for rowNo := 1; rows.Next(); rowNo++ {
	return nil
} //func (self Database) IterateDataRows(params *loadExceSheetParams,
