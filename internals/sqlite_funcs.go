package internals

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

var sqliteDBRegistry []*sql.DB = make([]*sql.DB, 0)

func startSQLite(sqlite_file_name string) *sql.DB {
	sFName := sqlite_file_name
	if !(strings.Contains(sFName, ":") || strings.Contains(sFName, "?")) {
		sFName = "file:" + sqlite_file_name + "?cache=shared&mode=rwc&_busy_timeout=5000"
	}

	// fmt.Printf("sFName: %v\n", sFName)
	db, err := sql.Open(driverString, sFName)
	if err != nil {
		fmt.Printf("Error opening DB:%s\n", err)
		log.Fatal(err)
	}
	//TODO: check if this connection already present. Possibly create DB connection right here for the each file separately (this will allow adding plugins)
	//TODO: check what happens if two connections work on the same file
	//TODO: add remove connection
	sqliteDBRegistry = append(sqliteDBRegistry, db)
	return db
} //func startSQLite(sqlite_file_name string) *sql.DB {

// TODO: rethink below as now col names are quoted, so maybe just quotes are to be removed (or masked)
var colNameSymbolsToBeReplaced = map[string]string{
	"\n":     "",
	"/":      "",
	"'":      "",
	" ":      "_",
	"\u00a0": "",
}

// TODO: make a map of symbols to be replaced
func cleanSQLColName(xl_col_name string) string {
	ret := xl_col_name
	for k, v := range colNameSymbolsToBeReplaced {
		ret = strings.ReplaceAll(ret, k, v)
	}
	return ret
}

// TODO: consider this answer: https://stackoverflow.com/a/43049720
func cleanSQLTableName(xl_table_name string) string {
	ret := xl_table_name
	for k, v := range colNameSymbolsToBeReplaced {
		ret = strings.ReplaceAll(ret, k, v)
	}
	return ret
}

func prepareTableCreateStatement(table_name string,
	colArray *[]*XLsqlColumn) (sqlStmt string, err error) {

	// colNamesArr *[]string,
	// colNamesTypes *map[string]string

	if len(*colArray) == 0 {
		return "", fmt.Errorf("prepareTableCreateStatement. Empty colNamesArr is supplied.")
	}

	colNamesSQL := ""
	// colNamesSQLArr = make([]string, 0)
	for _, column := range *colArray {
		// colType := (*colNamesTypes)[colName]
		colType := column.sqlColType
		// colNameClean := cleanSQLColName(column.xlColName)
		column.sqlColName = cleanSQLColName(column.xlColName)

		colNamesSQL = fmt.Sprintf(
			"%s '%s' %s,", colNamesSQL, column.sqlColName, colType)
		// colNamesSQLArr = append(colNamesSQLArr, column.sqlColName)
	}
	colNamesSQL = colNamesSQL[:len(colNamesSQL)-1] //To remove trailing comma

	tableName := cleanSQLTableName(table_name)
	sqlStmt = fmt.Sprintf(
		`	create table %s (%s);
			delete from "%s"; `,
		tableName, colNamesSQL, tableName)

	return sqlStmt, nil
} //func prepareTableCreateStatement(table_name string, colNamesArr *[]string,

func (self Database) createTable(sheetName string,
	columnsArr *[]*XLsqlColumn,
	auto_rename_table_name bool) (tblName string, err error) {
	//TODO: check the list of tables instead of db calls to check if database exists

	// colNamesArr *[]string,
	// colNamesTypes *map[string]string,

	tableName := sheetName
	i := 0
	// tx, err := self.db_connection.Begin()
	// defer tx.Commit()
	for {
		sqlStmt, err := prepareTableCreateStatement(tableName, columnsArr)
		if err != nil {
			return "", err
		}

		// fmt.Printf("sqlStmt: %v\n", sqlStmt)

		err = self.execSQLInternal(sqlStmt)
		if err != nil { //Table exists most of the times. Should create new one
			if err.Error() != fmt.Sprintf("table %s already exists", tableName) {
				fmt.Printf("createTable: Unknown table create error:: %v\n", err)
				return "", err
			}
			if !auto_rename_table_name {
				return "", fmt.Errorf("table %s already exists and auto_rename_table_name is false. Doing nothing. Table is not created.", tableName)
			}
			tableName = cleanSQLTableName(
				fmt.Sprintf("%s_%d", sheetName, i))
			i++
		} else {
			return tableName, nil
		} //if err != nil {//Table exists most of the times. Should create new one
	} //for {
} //func createTable(sheetName string, colNames []string) {

func prepareInsertStatementFromArray2(
	tbl_name string,
	colNamesValuesSlice []string) (insert_sql string) {
	colNamesSQL := ""
	colValsSQL := ""
	for _, key := range colNamesValuesSlice {
		if len(colNamesSQL) == 0 {
			colNamesSQL = fmt.Sprintf("\"%s\"", RemoveQuotesFromString(key))
			colValsSQL = fmt.Sprintf("?")
		} else {
			colNamesSQL = fmt.Sprintf("%s, \"%s\"", colNamesSQL, RemoveQuotesFromString(key))
			colValsSQL = fmt.Sprintf("%s, ?", colValsSQL)
		}
	} //for key, val := range colNamesValuesMap {
	insert_sql = fmt.Sprintf("INSERT INTO '%s' \n\t(%s) \n\tVALUES (%s)",
		tbl_name, colNamesSQL, colValsSQL)
	return insert_sql
} //func prepareInsertStatementFromMap(

func prepareInsertStatement(tx *sql.Tx, tableName string, columnsArr *([]*XLsqlColumn)) *sql.Stmt {
	//colNames []string

	qMarksSQL := ""
	//Reduce function is from: https://gosamples.dev/generics-reduce/
	// colNamesSQL := reduce(columnsArr, func(accum string, value *XLsqlColumn) string {
	// 	qMarksSQL = qMarksSQL + "?,"
	// 	return fmt.Sprintf("%s '%s',", accum, (*value).sqlColName)
	// }, "")

	colNamesSQL := ""
	for _, column := range *columnsArr {
		qMarksSQL = qMarksSQL + "?,"
		colNamesSQL = fmt.Sprintf("%s '%s',", colNamesSQL, column.sqlColName)
	}

	colNamesSQL = colNamesSQL[:len(colNamesSQL)-1]
	qMarksSQL = qMarksSQL[:len(qMarksSQL)-1]
	insString := fmt.Sprintf(
		"insert into %s (%s) values(%s)",
		tableName,
		colNamesSQL,
		qMarksSQL)

	// fmt.Printf("insString: %v\n", insString)

	res, err := tx.Prepare(insString)
	if err != nil {
		log.Fatal(err)
	}
	return res
} //func prepareInsertStatement(tx *sql.Tx, tableName string, colNames []string) *sql.Stmt {
