package internals

import (
	"database/sql"
	"io"
	"log"
	"os"
	"time"

	"github.com/mattn/go-sqlite3"
)

var debugMode = false
var initialized = false
var driverString = "sqlite3.1"

func Init(debug bool, extensionsSlice *[]string) {
	if !initialized {
		tm_str := time.Now().Format("2006-Jan-02_15H04m05.00000")
		main_sql_file_name := "./.test_dbs/DB_" + tm_str + ".sqlite3"
		debugMode = debug
		initLogging()
		//TODO: make this parameters to the command line utility or call to starlark function open_db
		if extensionsSlice != nil {
			driverString = "sqlite3_with_extensions"
			sql.Register(driverString,
				&sqlite3.SQLiteDriver{
					Extensions: *extensionsSlice,
				})
		} else {
			sql.Register(driverString,
				&sqlite3.SQLiteDriver{})
		}
		RegisterExportedFunctions(main_sql_file_name)
		initialized = true
	} //if !initialized {
} //func Init(debug bool, extensionsSlice *[]string) {

func initLogging() {
	file, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	if debugMode {
		DL = log.New(file, "DBG: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		DL = log.New(io.Discard, "DBG: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	DLf = DL.Printf
} //func InitLogging(devMode bool) {
