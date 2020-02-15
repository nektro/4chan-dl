#!/usr/bin/env bash

set -e
set -x

go test
go build
./4chan-dl \
