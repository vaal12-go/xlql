# XLQL starlark API 

<!--toc-->
- [XLQL starlark API](#xlql-starlark-api)
- [Overview](#overview)
- [Global functions](#global-functions)
    * [open_db](#open_db)
    * [get_datetime_formatted](#get_datetime_formatted)
    * [list_worksheets](#list_worksheets)
- [Objects](#objects)
    * [Database object](#database-object)
        * [class database](#class-database)
            * [load_excel_sheet](#load_excel_sheet)
            * [get_tables](#get_tables)
            * [exec_sql](#exec_sql)
    * [Query object](#query-object)
        * [class query](#class-query)
            * [print](#print)
            * [get_cursor](#get_cursor)
                * [save_to_excel](#save_to_excel)
    * [Cursor object](#cursor-object)
        * [class cursor](#class-cursor)

<!-- tocstop -->


# Overview
XLQL exports for use in starlark scripts run with the application following global functions and objects


# Global functions

## open_db

*function* open_db(file_name string) _returns Database object_
    
Will open/create sqlite database [file_name]. Can be absolute or relative path. If file does not exists - will create one.

As for most parameters this one is positional - parameter name can be omitted e.g.: 
    
    open_db("qwe1.sqlite")

## get_datetime_formatted

_function_ get_datetime_formatted(format string) _returns string_

Official doc: https://go.dev/src/time/format.go
an article: https://golang.cafe/blog/golang-time-format-example.html

Default format (no format string passed): "2006-01-02[15.04.05]"

Example format: "2006-Jan-02_15H04m05.00000" 
    
    print(get_datetime_formatted("Mon, 02 Jan 2006 15:04:05 MST"))

## list_worksheets

_function_ list_worksheets(xl_file_name string) _returns list of strings_


# Objects
## Database object
### class database

#### load_excel_sheet

*method* load_excel_sheet(file_name string, sheet_name string, drop_table bool = False, cols dict) *returns query object*

parameters should be indicated by name, those are not positional.

drop_table parameter is optional and is false by default. If true load_excel_sheet will drop the sqlite table with the name as derived from sheet_name and will load data from XL file into newly created table.

cols is optional and a dictionary of the form 
{
        "ColName1" : {
            "type" : "numeric"  #Types can be: "numeric" and "date"
        },
        "ColName2" : {
            "type" : "date",
            "format" : "2006-01-02" #Format is golang datetime format string as listed here: https://pkg.go.dev/time#pkg-constants - only year, month, day, hour, minute, second and AM/PM marks should be used (see below for convenience)
        }
}

Date formatting which can be used in __"format"__ field of __"cols"__ parameter:
* Year: "2006" "06"
* Month: "Jan" "January" "01" "1"
* Day of the month: "2" "_2" "02"
* Day of the week: "Mon" "Monday"
* Day of the year: "__2" "002"
* Hour: "15" "3" "03" (PM or AM)
* Minute: "4" "04"
* Second: "5" "05"
* AM/PM mark: "PM"

   Please note that golang is pretty strict about time formatting, so although "2023-08-12" date string will perfectly match "2006-01-02" format string, "2023-08-12 00:00:00" date string will produce an error if parsed with "2006-01-02" format.

   Parsing errors will be visible in SQLite file with the message (instead of date) like *"Cannot parse date:2023-08-12 00:00:00 with format:2006-01-02"*. If saved to excel, such message will produce incorrect date (*0001-01-01 00:00:00 +0000 UTC*, so called zero date in golang). This is due to the driver attempt to convert the date (which is an error string). To see an error in Excel use cast SQLite function. E.g. column "datetime_with_errors" in Sheet1 table contains errors: 

        SELECT cast(datetime_with_errors As text) FROM Sheet1 
   
   will show actual errors in excel like below:

        cast(datetime_with_errors As text)
        Cannot parse date:2023-08-12 00:00:00 with format:2006-01-02 
        Cannot parse date:2023-12-09 12:38:00 with format:2006-01-02 
        Cannot parse date:2023-12-10 00:00:00 with format:2006-01-02 
   
In the example above ("2023-08-12 00:00:00" datetime value) will be perfectly parsed with "2006-01-02 03:04:05" format string.

In future releases I will check for possibility to provide a custom parser function (user-written) to parse dates more reliably.

#### get_tables

_method_ get_tables() _returns list of strings_
TBD

#### exec_sql

_method_ exec_sql(sql_statement string)

Will execut SQL statement (e.g. "DROP TABLE sometable") which will not return any result.



## Query object
### class query
#### print

*method* print() *prints content of the query to console*
Will print content of the query to console. See examples/example4.star. 
Will try to provide colorized representation. 

![query.print() example to console](/documentation/img/query_print_example.png)

Not optimized for big tables (large number of columns or large cell content) though.

#### get_cursor

*method* get_cursor() *returns cursor object for iterating query rows*

From examples/example4.star:
    mem_q = inMemDB.run_query("SELECT * FROM qwe1")
    for row in mem_q.get_cursor():
        print(row)

This is only needed for iteration through starlark 'for' cycles.
Each iteration provides an array of fields stored in each row which is being iterated.


##### save_to_excel

method save_to_excel(file_name string, sheet_name string) returns none

Will save query to the excel worksheet with style 

## Cursor object
### class cursor
Is needed to iterate through result of the query (see examples/example4.star).
Only needed for iteration with 'for' loops


<!-- 
    "github.com/ycd/toc/pkg/toc" replaces "_" with "-" in TOC links. 
    TODO: To check if this can be amended. In the meantime should be replaced manually

 -->