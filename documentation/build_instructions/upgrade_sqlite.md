# SQLite driver and version upgrade

XLQL depends on the [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) driver to operate.
This is a compiled driver, which brings along the corresponding version of the sqlite C bindings.

At the moment (2026Apr18) driver version v1.14.42 is dependent on the sqlite bindings version [3.51.3](https://sqlite.org/releaselog/3_51_3.html)

To upgrade the driver version a gcc toolchain must be installed for your platform (see [Installation](https://github.com/mattn/go-sqlite3#installation) section of the driver documentation)

Golang version currently used is 1.23.6.

## Upgrade procedure 
1. Update the driver code:
```console
go get github.com/mattn/go-sqlite3
```

2. Tidy the module
```console
go mod tidy
```

3. Build the driver version:
```console
go build ./cmd
```


### To check current version of the driver package:
```console
go list -m github.com/mattn/go-sqlite3
```

### To check SQLite version 
run the examples\sqlite_ver.star script
```console
go run ./cmd/main.go -f ./examples/sqlite_ver.star
```