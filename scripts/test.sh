#!/bin/bash

go clean -testcache
go test -p 1 ./...