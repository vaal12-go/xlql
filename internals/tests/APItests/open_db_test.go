package APItests

import (
	"testing"

	"test.com/excel-ark/internals"
)

const DB_REMOVE_FILE = "./test_scripts/open_db_remove_file.star"

func TestQuery(t *testing.T) {
	internals.Init(true, nil)
	t.Logf("Runnin test file:%s\n", DB_REMOVE_FILE)
	internals.ExecStarlarkFile(DB_REMOVE_FILE)
	internals.Close()
	t.Logf("Running the file 2nd time to test if existing DB file is deleted.")
	internals.Init(true, nil)
	internals.ExecStarlarkFile(DB_REMOVE_FILE)
	internals.Close()
}
