# Excel-SQL




## Overview

Excel-SQL is a command line utility which aims to provide programmatical access to Excel files by exporting those to sqlite database and exeuting programs in python like language (starlark) against such databases.

![The San Juan Mountains are beautiful!](/TOTALCMD64_vFs0BvTqcc.gif "San Juan Mountains")

<!--toc-->
<!-- tocstop -->

## Installation

Download an executable for your platform from the releases section.

## Basic use

``` excel-sql.exe -f file_to_execute.star ```


## Command line parameters



## Starlark language

This application is based on excellent [starlark implementation in golang](https://github.com/google/starlark-go/tree/master). 

To quote from the site above
 
> Starlark is a dialect of Python intended for use as a configuration language. 

For complete syntax and specification see [Starlark in Go: Language definition](https://github.com/google/starlark-go/blob/master/doc/spec.md#function-and-method-calls)


## API

### Overview

### Global functions


### Database

### Query

### Cursor


## Thanks (libraries used)

* github.com/jedib0t/go-pretty/v6 v6.4.8
* github.com/joho/godotenv v1.5.1
* github.com/mattn/go-sqlite3 v1.14.17
* github.com/xuri/excelize/v2 v2.8.0
* go.starlark.net v0.0.0-20231013162135-47c85baa7a64









