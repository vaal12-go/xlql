package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"test.com/excel-ark/internals"
)

const (
	VER_STRING = "xlql (Excel query language)"
)

var (
	fName             string
	printVersion      bool
	printHelp         bool
	version           string
	build_date        string
	commit            string
	ver_codename      string
	ver_hash          string
	sqlite_extensions string
	ver_sqlite        string
	build_time        string
)

func main() {
	flag.StringVar(&fName, "f", "", "Starlark file to execute.")
	flag.BoolVar(&printHelp, "h", false, "Print help message")
	flag.BoolVar(&printVersion, "v", false, "Print version message")
	flag.StringVar(&sqlite_extensions, "e", "",
		"list of sqlite extensions (separated by question mark to load. No wildcards in file names yet.")
	//TODO: explore if possible to add wildcards for extension names
	flag.Parse()
	if printHelp {
		flag.PrintDefaults()
		return
	}
	if printVersion {
		ver_str := fmt.Sprintf("%s \n\tversion:%s (build time:%s UTC+3) %s \n\tsqlite version:%s\n",
			VER_STRING, version, build_time, build_date, ver_sqlite)
		if ver_codename != "" {
			ver_str = ver_str + "\tcodename:" + ver_codename
		}
		if ver_hash != "" {
			ver_str = fmt.Sprintf("%s [git hash:%s]", ver_str, ver_hash)
		}
		fmt.Println(ver_str)
		return
	}
	var DEVELOPER_MODE = ""
	err := godotenv.Load(".env")
	dbgMode := false
	if err == nil {
		DEVELOPER_MODE = os.Getenv("DEVELOPER_MODE")
		if DEVELOPER_MODE == "" {
			dbgMode = false
		} else {
			fmt.Printf("DEV mode\n")
			dbgMode = true
		}
	} else {
		//Logging is not initiated yet
		// fmt.Printf("Cannot load .env file. Considering this a non debug mode Err: %s", err)
	}
	if fName == "" {
		if DEVELOPER_MODE == "" {
			// log.Fatal("No .star file specified with -f paramater. Exiting.")
			flag.PrintDefaults()
			return
		} else {
			fName = "test.star"
		}
	}

	// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
	if _, err := os.Stat(fName); err == nil {
		// fmt.Printf("File exists")
	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		fmt.Printf("File '%s' DOES NOT exist. Exiting.\n", fName)
		os.Exit(1)

	} else {
		fmt.Printf("File may or may not exist. Will proceed with processing file, though it may fail.")
		// Schrodinger: file may or may not exist. See err for details.
	}
	var extSlice *[]string = nil
	if sqlite_extensions != "" {
		extSlc := strings.Split(sqlite_extensions, "?")
		extSlice = &extSlc
	}

	internals.Init(dbgMode, extSlice) //Logging is also initiated here
	//See internals/utils.go for available loggin methods
	// internals.InfoLogger.Println("\n\n\n\n\n ******************************************* ")
	internals.ExecStarlarkFile(fName)
	internals.Close()
} //func main() {

//SHA-512: c46c48ddaccc53835237b551df240c1dc51ca78911fdec845481c04bb2bb438b71ed15431131850dd2654a136047dbe3129c728d4148f971d3bc0d7568124a86
