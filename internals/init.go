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

func Init(debug bool) {
	tm_str := time.Now().Format("2006-Jan-02_15H04m05.00000")
	main_sql_file_name := "./.test_dbs/DB_" + tm_str + ".sqlite3"
	debugMode = debug
	initLogging()
	//This is needed to register sqlite with extensions
	//TODO: make this parameters to the command line utility or call to starlark function open_db
	sql.Register("sqlite3_with_extensions",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				//TODO: check relative path
				// `c:\Users\may13\AGVDocs\Dev\2. Go\10.LUA-test\sqliteextensions\pivotvtab.dll`,
				// `./sqliteextensions/pivotvtab.dll`,
			},
		})
	RegisterExportedFunctions(main_sql_file_name)
}

func initLogging() {
	file, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	// fmt.Printf("\"Initializing init\": %v\n", "Initializing init")
	if debugMode {
		DL = log.New(file, "DBG: ", log.Ldate|log.Ltime|log.Lshortfile)
		// DLf = DL.Printf
		// fmt.Printf("DL.Printf: %v\n", DL)
	} else {
		DL = log.New(io.Discard, "DBG: ", log.Ldate|log.Ltime|log.Lshortfile)
		// DLf = DiscardPrintf
		// fmt.Printf("DL.Printf: %v\n", DL)
	}
	DLf = DL.Printf
} //func InitLogging(devMode bool) {
