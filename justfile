style:
    gofmt -w cmd pkg

check: style
    go vet ./...
    golangci-lint run

clean:
    go mod tidy

test:
    go test -v -cover ./...

update:
    go get -u ./...
