#!/usr/bin/env bash

set -e
go test
go build
./4chan-dl \
