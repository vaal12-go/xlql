# XLQL starlark API 

<!--toc-->
- [XLQL starlark API](#xlql-starlark-api)
- [Overview](#overview)
- [Global functions](#global-functions)
    * [open_db](#open-db)
    * [get_datetime_formatted](#get-datetime-formatted)
    * [list_worksheets](#list-worksheets)
- [Objects](#objects)
    * [Database object](#database-object)
        * [class database](#class-database)
            * [load_excel_sheet](#load-excel-sheet)
            * [get_tables](#get-tables)
    * [Query object](#query-object)
        * [class query](#class-query)
            * [print](#print)
            * [get_cursor](#get-cursor)
                * [save_to_excel](#save-to-excel)
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

*method* load_excel_sheet(file_name string, sheet_name string, drop_table bool = False) *returns query object*

parameters should be indicated by name, those are not positional.

drop_table parameter is optional and is false by default. If true load_excel_sheet will drop the sqlite table with the name as derived from sheet_name and will load data from XL file into newly created table.

#### get_tables

_method_ get_tables() _returns list of strings_
TBD


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