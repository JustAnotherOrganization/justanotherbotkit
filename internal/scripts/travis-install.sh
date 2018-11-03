#!/bin/sh

# Disable install for module builds.
if [ "$GO111MODULE" != on ]; then
    go get -d -t -v ./...
fi
