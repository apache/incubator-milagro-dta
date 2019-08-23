set -e

echo "Run go fmt"
test -z "$(gofmt -s -l . 2>&1 | grep -v vendor | tee /dev/stderr)"

echo "Run go lint"
golint -set_exit_status $(go list ./... | grep -v /vendor/)
