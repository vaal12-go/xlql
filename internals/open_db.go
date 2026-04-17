package internals

import (
	"errors"
	"fmt"
	"log"
	"os"

	"go.starlark.net/starlark"
)

func open_db(thread *starlark.Thread,
	b *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	var fName string
	var del_if_exists bool
	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "file_name", &fName,
		"delete_db_if_exists?", &del_if_exists); err != nil {
		return nil, err
	}
	// fmt.Println("open_db:26 del_if_exists::", del_if_exists)
	if del_if_exists {
		if _, err := os.Stat(fName); err == nil {
			fmt.Printf("DB File '%s' exists. Will attempt deletion.", fName)
			if err := os.Remove(fName); err != nil {
				fmt.Printf("Cannot delete DB file '%s'. Possibly locked by other application. Exiting.")
				log.Fatal(err)
			}
			fmt.Printf("File '%s' is deleted. Proceeding with creation of new DB.", fName)
		} else if errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does *not* exist
			fmt.Printf("DB File '%s' DOES NOT exist. Exiting.\n", fName)
		} else {
			fmt.Printf("DB File may or may not exist. Will proceed with processing file, though it may fail.")
			// Schrodinger: file may or may not exist. See err for details.
		}
	}

	ret := NewDatabase(fName)
	return ret, nil
}
