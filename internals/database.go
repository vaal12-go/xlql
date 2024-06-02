package internals

import (
	"database/sql"
	"fmt"
	"slices"

	_ "github.com/mattn/go-sqlite3"
	"go.starlark.net/starlark"
)

//TODO: clean file

// TODO: review/update logging with logging levels
// LOW: to consider if this id_counter is needed at all
var id_counter = 0

// HIGH: create documentation to classes and functions. Also create diagramms: https://d2lang.com/tour/uml-classes
type Database struct {
	db_id                           int64
	db_file_name                    string
	db_name                         string
	db_connection                   *sql.DB
	database_builtinsMap            map[string](*starlark.Builtin)
	runQueryMethodBuiltin           *starlark.Builtin
	loadExcelWorksheetMethodBuiltin *starlark.Builtin
	getTablesMethodBuiltin          *starlark.Builtin
	createTableMethodBuiltin        *starlark.Builtin
	exporter                        MethodExporter
	StarlarkValueImplementationStub
}

// TODO: add new_inmem_db function to create temp DB quickly and this will replace default DB
// TODO: add copying of inmem db to disk:
// https://pkg.go.dev/github.com/mattn/go-sqlite3#SQLiteConn.Backup
var database_export_fields_array = []string{
	"Run_query", "Load_excel_sheet", "Get_tables", "Create_table",
	"Exec_sql",
}

// TODO: check if buiiltins can be implemented with a map - e.g. map of starlark name as key and builtin itself as value
//
//	this way Attr will also be simpler
func NewDatabase(file_name string) *Database {
	id_counter++
	ret := Database{
		db_id:                int64(id_counter),
		db_file_name:         file_name,
		db_connection:        startSQLite(file_name),
		database_builtinsMap: make(map[string]*starlark.Builtin),
	}
	//below var is necessary to get proper pointer to ret
	var iFacePointer interface{} = ret
	// ret.exporter.RegisterBuiltIns(&iFacePointer,
	// 	[]string{"Run_query", "Load_excel_sheet", "Get_tables", "Create_table"})
	ret.exporter.RegisterBuiltIns(&iFacePointer,
		database_export_fields_array)
	return &ret
} //func NewDatabase(file_name string) *Database {

func (self Database) dropTableIfExists(tblName string) error {
	tbls := self.get_tables_actual()
	// fmt.Printf("dropTableIfExists tbls: %v\n", tbls)
	var err error = nil
	if slices.Contains(tbls, tblName) {
		// fmt.Printf("\"dropping table\": %v\n", "dropping table")
		err = self.execSQLInternal(
			`DROP TABLE ` + tblName)
	}
	return err
}

// TODO: add execSQL function for starlark
func (self Database) execSQLInternal(sql_str string) error {
	_, err := self.db_connection.Exec(sql_str)
	if err != nil {
		fmt.Printf("execSQLInternal. Error executing query:%s \n\terr: %v\n", sql_str, err)
		return err
	}
	return err
} //func (self Database) execSQLInternal(sql_str string) {

func (self Database) Exec_sql(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {

	sql_str := ""
	if err := starlark.UnpackArgs(
		b.Name(), args, kwargs,
		"sql?", &sql_str); err != nil {
		return nil, err
	}

	if sql_str == "" {
		// fmt.Printf("Fatal: sql parameter of exec_sql cannot be empty.\n")
		return starlark.None, fmt.Errorf("sql parameter of exec_sql cannot be empty.")
	}

	// fmt.Printf("\"Exec SQL is running\": %v\n", "Exec SQL is running")
	// fmt.Printf("sql_str: %v\n", sql_str)

	err := self.execSQLInternal(sql_str)

	return starlark.None, err
} //func (self Database) Exec_sql(thread *starlark.Thread,

func (self Database) Run_query(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	//TODO add parameter that this is query without result or add other function something like run_sql
	query_SQL := ""
	if err := starlark.UnpackArgs(
		b.Name(), args, kwargs,
		"query", &query_SQL); err != nil {
		fmt.Printf("run_query:%v\n", err)
		return nil, err
	}
	// fmt.Printf("Executing SQL:%s\n", query_SQL)

	//TODO: check need to add and commit transaction for each statement like below
	// tx, err := self.db_connection.Begin()
	// tx.Commit()
	err := self.execSQLInternal(query_SQL)
	if err != nil {
		return nil, err
	}
	// _, err := self.db_connection.Exec(query_SQL)

	// res, err = self.db_connection.Exec("COMMIT TRANSACTION")

	// if err != nil {
	// 	DLf("Database.Run_query. Error running query:%v\n", err)
	// 	return nil, err
	// }
	// fmt.Printf("\"6\": %v\n", "6")
	ret, err := NewQuery(&self, query_SQL, "")
	// fmt.Printf("\"7\": %v\n", "7")
	return ret, err
} //func (self Database) run_query(thread *starlark.Thread,

func (self Database) Attr(name string) (starlark.Value, error) {
	return self.exporter.GetMethod(name)
}

func (self Database) Close() {
	self.db_connection.Close()
}

func (self Database) AttrNames() []string {
	ret := []string{"run_query", "load_excel_sheet",
		"get_tables", "create_table"}
	return ret
}

func (self Database) String() string {
	return fmt.Sprintf("Database. FileName:%s", self.db_file_name)
}
