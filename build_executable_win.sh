# This script is created to set building parameters for xlql executable build
# It was tested on bash (cygwin) on windows

# GENERAL BUILD PARAMETERS
VER_NUMBER=0.5.01
VERSION_CODE_NAME="2026Apr18_RELEASE"
SQLITE_VER="3.51.3"

echo "STARTING xlql executable build (Windows): $(date -u)"
echo "     VER_NUMBER=$VER_NUMBER"
echo "     VERSION_CODE_NAME=$VERSION_CODE_NAME"
echo "     SQLITE_VER=$SQLITE_VER"


# REM GETTING CURRENT GIT HASH
LAST_COMMIT_HASH=$(git show -s --format=%h)
UTC_DATETIME=$(date -u +"%Y-%b-%e_%k-%MZ")
# echo $UTC_DATETIME
# set /p LAST_COMMIT_HASH=<build\temp.txt
# del build\temp.txt

# echo $LAST_COMMIT_HASH
LINKER_PARAMETERS="-w -s -X main.version=$VER_NUMBER  -X main.ver_codename=$VERSION_CODE_NAME -X main.ver_hash=$LAST_COMMIT_HASH -X main.ver_sqlite=$SQLITE_VER -X main.build_time=$UTC_DATETIME"
# echo $LINKER_PARAMETERS

# REM ACTUAL BUILD COMMAND
# Q1="go build -o ./build/xlql.exe -ldflags "$LINKER_PARAMETERS" ./cmd/main.go"

# echo $Q1
set -o xtrace
go build -o ./build/xlql.exe -ldflags "$LINKER_PARAMETERS" ./cmd/main.go
set +o xtrace

echo "BUILD SUCCESSFULL: $(date -u)"
# @ECHO OFF
# @ECHO --------------------------------   
# @ECHO BUILD SUCCESSFULL
# @ECHO --------------------------------


