package main

import (
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

	// excelizeFile, err := excelize.OpenFile("GoodBlueTableStyle_29Apr2024.xlsx")
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	panic(1)
	// }

	// tbls, err := excelizeFile.GetTables("Sheet1")
	// if err != nil {
	// 	fmt.Printf("Getting tables err: %v\n", err)
	// 	panic(1)
	// }

	// for _, tbl := range tbls {
	// 	fmt.Printf("tbl: %v\n", tbl)
	// }

	// panic(1)

	var extSlice *[]string = nil
	if sqlite_extensions != "" {
		// fmt.Printf("sqlite_extensions: %v\n", sqlite_extensions)
		extSlc := strings.Split(sqlite_extensions, "?")
		extSlice = &extSlc
	}

	internals.Init(dbgMode, extSlice) //Logging is also initiated here
	//See internals/utils.go for available loggin methods
	// internals.InfoLogger.Println("\n\n\n\n\n ******************************************* ")
	internals.ExecStarlarkFile(fName)
	internals.Close()
} //func main() {

//SHA256-9cd7c8b6e0e7cf266accf920c4ec53d133e0e5f5eb3635506140ee8ef7a514d0
