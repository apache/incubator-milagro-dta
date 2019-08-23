#! /bin/bash

GO111MODULE=on go test -race -cover `go list ./... | grep -v disabled`

status=$?

exit $status
