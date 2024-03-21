package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"test.com/excel-ark/internals"
)

const (
	VER_STRING = "excel_scripter version"
)

var (
	fName        string
	printVersion bool
	printHelp    bool

	version      string
	build_date   string
	commit       string
	ver_codename string
	ver_hash     string
)

func main() {
	flag.StringVar(&fName, "f", "", "Starlark file to execute.")
	flag.BoolVar(&printHelp, "h", false, "Print help message")
	flag.BoolVar(&printVersion, "v", false, "Print version message")
	flag.Parse()
	if printHelp {
		flag.PrintDefaults()
		return
	}
	if printVersion {
		ver_str := fmt.Sprintf("%s \"%s\" %s \n",
			VER_STRING, version, build_date)
		if ver_codename != "" {
			ver_str = ver_str + "name:" + ver_codename
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
			fmt.Printf("DEV mode")
			dbgMode = true
		}
	} else {
		//Logging is not initiated yet
		// fmt.Printf("Cannot load .env file. Considering this a non debug mode Err: %s", err)
	}
	if fName == "" {
		if DEVELOPER_MODE == "" {
			log.Fatal("No .star file specified with -f paramater. Exiting.")
		} else {
			fName = "test.star"
		}
	}
	internals.Init(dbgMode) //Logging is also initiated here
	//See internals/utils.go for available loggin methods
	internals.InfoLogger.Println("\n\n\n\n\n ******************************************* ")
	internals.ExecStarlarkFile(fName)
	internals.Close()
} //func main() {

//SHA256-9cd7c8b6e0e7cf266accf920c4ec53d133e0e5f5eb3635506140ee8ef7a514d0
