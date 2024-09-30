@ECHO OFF
REM GENERAL BUILD PARAMETERS
SET VER_NUMBER=0.4.01
SET VERSION_CODE_NAME="30Sep2024_RELEASE"
SET SQLITE_VER="3.46.1"

REM GETTING CURRENT GIT HASH
git show -s --format=%%h > build\temp.txt
set /p LAST_COMMIT_HASH=<build\temp.txt
del build\temp.txt

REM -X main.build_date="%date%"
SET LINKER_PARAMETERS="-w -s -X main.version=%VER_NUMBER%  -X main.ver_codename=%VERSION_CODE_NAME% -X main.ver_hash=%LAST_COMMIT_HASH% -X main.ver_sqlite=%SQLITE_VER% -X main.build_time=%TIME% "

REM ACTUAL BUILD COMMAND
@ECHO ON
go build -o ./build/xlql.exe -ldflags %LINKER_PARAMETERS% ./cmd/main.go

@ECHO OFF
@ECHO --------------------------------   
@ECHO BUILD SUCCESSFULL
@ECHO --------------------------------


