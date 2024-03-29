# Comand line parameters

## -f 

file name of starlark file which is to be executed. Only one file supported at the moment.

    xlql.exe -f test.star

## -e

SQLIte Extensions which are to be loaded to the driver. 

A list of path(s) to .dll files of extensions. Relative paths to current directory are supported. Multiple paths are to be separated by question mark "?".

    xlql.exe -f test.star -e "./sqliteextensions/pivotvtab.dll?c:\Users\may13\Desktop\dir with spaces\uuid.dll"

TODO: No wildcards are supported yet.

Double quotes on Windows are necessary if file path(s) contain spaces. Otherwise - not necessary.

At the moment this is only test on Windows, but sqlite .dll extensions seems to be supported. See "sqliteextensions" folder for sample extensions.

See examples\example_extension_UUID_simple.star. *sqliteextensions\uuid.dll is needed to run this example.*


Sources of sqlite extensions:
- [github.com/nalgeon](https://github.com/nalgeon/sqlean?tab=readme-ov-file#main-set)
- [sqlpkg.org](https://sqlpkg.org/)



## -v

Will print version and date of the xlql.exe file.


## -h 

Will print short help message on command line parameters.
