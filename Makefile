# This is how we want to name the binary output
BINARY=multiline-test

# These are the values we want to pass for Version and BuildTime
VERSION=0.1.0
BUILD_TIME=`date +%Y-%m-%d`
GIT_HASH=`git rev-parse --verify HEAD`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS="-X main.buildDate=`date +%Y-%m-%d` -X main.version=${VERSION} -X main.commitHash=${GIT_HASH}"

all:
	go build -ldflags ${LDFLAGS} -o ${BINARY} main.go
