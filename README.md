# XLQL (excel SQL)

## Overview

xlql is a command line utility which aims to enable SQL queries and python like language access to data in Excel files. This is done by exporting Excel data from excel files to sqlite database. SQLite SQL can be used against such data. Export to Excel of resulting data is also supported at the moment.

<!--toc-->
<!-- tocstop -->

## Installation

Download an executable for your platform from the releases section.

## Basic use

``` xlql.exe -f file_to_execute.star ```

## Examples
Super basic [Example 1](/examples/example.star)
will just load countries sheet from the Sample HR Database.xlsx and print it's content to console. 

> *Run this example from the directory, which has this Sample HR Database.xlsx file.*

[Example 2](/examples/example2.star) shows how to run SQL join on Excel data.
> *Run this example from the directory, which has this Sample HR Database.xlsx file.*

[Example 3](/examples/example3.star) shows hot data generated by starlark script (Taylor approximation for sin function in range -3 to 4) can be exported to Excel, where a graph can be generates as below

![Sin graph](/README-imgs/sine_approximation_data%20and%20graph.png)



## Command line parameters
TBD



## Starlark language

This application is based on excellent [starlark implementation in golang](https://github.com/google/starlark-go/tree/master). 

To quote from the site above
 
> Starlark is a dialect of Python intended for use as a configuration language. 

For complete syntax and specification see [Starlark in Go: Language definition](https://github.com/google/starlark-go/blob/master/doc/spec.md#function-and-method-calls)


## API

### Overview
TBD

### Global functions
TBD


### Database object
TBD

### Query object
TBD

### Cursor object
TBD

## Thanks (libraries used)

* github.com/jedib0t/go-pretty/v6 v6.4.8
* github.com/joho/godotenv v1.5.1
* github.com/mattn/go-sqlite3 v1.14.17
* github.com/xuri/excelize/v2 v2.8.0
* go.starlark.net v0.0.0-20231013162135-47c85baa7a64









