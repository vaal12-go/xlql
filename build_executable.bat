SET VER_NUMBER=0.2.0
SET VERSION_CODE_NAME="MAR24_RELEASE"

git show -s --format=%%h > build\temp.txt
set /p LAST_COMMIT_HASH=<build\temp.txt
del build\temp.txt


go build -o ./build/xlql.exe -ldflags "-w -s -X main.version=%VER_NUMBER% -X main.build_date=%date% -X main.ver_codename=%VERSION_CODE_NAME% -X main.ver_hash=%LAST_COMMIT_HASH%" ./cmd/main.go

