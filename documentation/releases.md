# Releases

## 2026Apr19
* Database object is fully documented in API docs
* Example (example_db_object.star) created to demonstrate database object methods in action
* Cleaning of internals/database.go file, also of starlark_exports.go, starlarkValueStub.go  

## 2026Apr18
* upgraded [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) v1.14.24 => v1.14.42
* the above updates SQLite version to [3.51.3](https://sqlite.org/releaselog/3_51_3.html)
* Added script /examples/sqlite_ver.star to check SQLite version after upgrade
* Added \documentation\build_instructions\upgrade_sqlite.md file with sqlite driver upgrade instructions
* New version of build script (build_executable_win.sh) which allows for more flexible time recording (in UTC) and also should be able to be run on Linux (tested only on cygwin)


## 2026Apr17
* open_db function accepts delete_db_if_exists optional parameter, which will delete db file (if such exists) before opening it for new.
    * Test coverage added
* Error message prints now print out all the errors encountered.
* Added check for .star file to exist before execution with appropriate error message.
* Small code cleans

## 28Mar2025 
* now development will be done in golang v1.23+
* Updated qax-os/excelize version to v2.9.0 (of 14-Oct-2024)
* Updated google/starlark-go package to v0.0.0-20250318223901-d9371fef63fe (18Mar2025)
* Tests are running well.



## Release steps
On the feature branch:
1. Run tests until executed successfully
2. Update this file (releases.md) with information on the release
3. Update release information (VER_NUMBER, VERSION_CODE_NAME, SQLITE_VER) in build_executable_win.sh
4. Test run build_executable_win.sh. If all good - commit changes to feature branch.
5. Squash the feature branch to the master branch.
6. Run build_executable_win.sh to capture git commit hash into executable
7. Push changes to the origin/master
8. Create tag for release
9. Create new release on github. Upload zipped executable.

