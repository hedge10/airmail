build:
    goreleaser build --single-target --snapshot --clean

check: style
    go vet ./...
    golangci-lint run

clean:
    go mod tidy

style:
    gofmt -w cmd pkg

test:
    go test -v -cover ./...

update:
    go get -u ./...
