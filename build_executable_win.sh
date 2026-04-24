# This script is created to set building parameters for xlql executable build
# It was tested on bash (cygwin) on windows

# GENERAL BUILD PARAMETERS
VER_NUMBER=0.5.02
VERSION_CODE_NAME="2026Apr19_RELEASE"
SQLITE_VER="3.51.3"

echo "STARTING xlql executable build (Windows): $(date -u)"
echo "     VER_NUMBER=$VER_NUMBER"
echo "     VERSION_CODE_NAME=$VERSION_CODE_NAME"
echo "     SQLITE_VER=$SQLITE_VER"

# GETTING CURRENT GIT HASH
LAST_COMMIT_HASH=$(git show -s --format=%h)
UTC_DATETIME=$(date -u +"%Y-%b-%e_%k-%MZ")

LINKER_PARAMETERS="-w -s -X main.version=$VER_NUMBER  -X main.ver_codename=$VERSION_CODE_NAME -X main.ver_hash=$LAST_COMMIT_HASH -X main.ver_sqlite=$SQLITE_VER -X main.build_time=$UTC_DATETIME"

# ACTUAL BUILD COMMAND
set -o xtrace
go build -o ./build/xlql.exe -ldflags "$LINKER_PARAMETERS" ./cmd/main.go
set +o xtrace

echo "BUILD SUCCESSFULL: $(date -u)"

