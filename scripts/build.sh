#!/usr/bin/env bash

APPNAME="authorizer"
export GO111MODULE=on

set -o errexit
set -o nounset

if [ -z $APPNAME ]; then
    echo "APPNAME must be set"
    exit 1
fi

export CGO_ENABLED=0

echo "Go building app"
go build -v -o build/$APPNAME cmd/$APPNAME/main.go
echo "Successfully built, exiting build script"