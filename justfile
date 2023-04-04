check: style
    go vet ./...
    golangci-lint run

clean:
    go mod tidy

down:
    docker-compose down -v

release:
    goreleaser build --single-target --snapshot --clean

style:
    gofmt -w cmd pkg

test:
    go test -v -cover ./...

up:
    docker-compose up

update:
    go get -u ./...
