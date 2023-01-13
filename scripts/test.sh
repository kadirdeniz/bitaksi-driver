#!/bin/bash

go clean -testcache
go test -p 1 $(go list ./... | grep -v /test/)
