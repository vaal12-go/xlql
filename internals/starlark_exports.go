// This is a high level doc of the package internals
// Should be visible in docs
package internals

import (
	"log"

	"go.starlark.net/starlark"
	syntax "go.starlark.net/syntax"
)

func ExecStarlarkFile(fName string) *starlark.StringDict {
	thread := &starlark.Thread{Name: "my thread"}
	//THis enables while loops outside functions and use of global variables, which can change
	//https://github.com/google/starlark-go/blob/47c85baa7a64/syntax/options.go#L23
	execOptions := syntax.FileOptions{
		GlobalReassign:  true,
		Recursion:       true,
		While:           true,
		TopLevelControl: true,
	}
	DLf("execOptions: %v\n", execOptions)
	DLf("thread: %v\n", thread)
	DLf("fName: %v\n", fName)
	DLf("PredeclaredDict: %v\n", PredeclaredDict)
	globals, err := starlark.ExecFileOptions(
		&execOptions, thread, fName, nil, PredeclaredDict)
	if err != nil {
		log.Fatalf("ExecStarlarkFile. Error:%v\n", err.Error())
	}
	DLf("Script finished. Globals returned:", globals)
	return &globals
} //func ExecStarlarkFile(fName string) *starlark.StringDict {

func Close() {
	for _, db := range sqliteDBRegistry {
		db.Close()
	}
}

var PredeclaredDict = starlark.StringDict{}

var DefaultDB *Database

func RegisterExportedFunctions(DefaultDBName string) {
	PredeclaredDict["open_db"] =
		starlark.NewBuiltin("open_db",
			open_db)
	PredeclaredDict["list_worksheets"] =
		starlark.NewBuiltin("list_worksheets",
			list_worksheets)

	PredeclaredDict["get_datetime_formatted"] =
		starlark.NewBuiltin("get_datetime_formatted",
			get_datetime_formatted)

		//TODO: consider if 'defaultDB' is needed at all
	DefaultDB = NewDatabase(DefaultDBName) //[x]: remove default DB?
	PredeclaredDict["defaultDB"] =
		DefaultDB
} //func RegisterExportedFunctions(DefaultDBName string) {
