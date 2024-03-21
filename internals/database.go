package internals

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"go.starlark.net/starlark"
)

// HIGH: review/update logging with logging levels
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

var database_export_fields_array = []string{
	"Run_query", "Load_excel_sheet", "Get_tables", "Create_table",
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
	ret.exporter.RegisterBuiltIns(&iFacePointer,
		[]string{"Run_query", "Load_excel_sheet", "Get_tables", "Create_table"})
	return &ret
} //func NewDatabase(file_name string) *Database {

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
	// _, err := self.db_connection.Query(query_SQL)
	// res, err := self.db_connection.Exec("BEGIN IMMEDIATE TRANSACTION")
	// fmt.Printf("\"2\": %v\n", "2")
	//TODO: check need to add and commit transaction for each statement like below
	// tx, err := self.db_connection.Begin()
	// tx.Commit()

	// fmt.Printf("\"3\": %v\n", "3")
	_, err := self.db_connection.Exec(query_SQL)
	// fmt.Printf("\"4\": %v\n", "4")
	// fmt.Printf("\"5\": %v\n", "5")
	// res, err = self.db_connection.Exec("COMMIT TRANSACTION")
	// fmt.Printf(" self.db_connection.Exec res: %v\n", res)
	if err != nil {
		DLf("Database.Run_query. Error running query:%v\n", err)
		return nil, err
	}
	// fmt.Printf("\"6\": %v\n", "6")
	ret := NewQuery(&self, query_SQL, "")
	// fmt.Printf("\"7\": %v\n", "7")
	return ret, nil
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
