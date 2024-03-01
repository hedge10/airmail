set dotenv-load

IMAGE := "hedge10/airmail"

build:
    docker build --no-cache -t {{ IMAGE }} .

check: style
    docker exec airmail go vet ./...
    docker exec airmail golangci-lint run

clean:
    docker exec airmail go mod tidy

down:
    docker-compose down -v

release:
    goreleaser release --snapshot --clean

style:
    docker exec airmail gofmt -w cmd pkg

run:
    docker run --rm -it -p 9900:9900 -e AM_SMTP_HOST="$AM_SMTP_HOST" -e AM_SMTP_PORT="$AM_SMTP_PORT" -e AM_SMTP_USER="$AM_SMTP_USER" -e AM_SMTP_PASS="$AM_SMTP_PASS" -e AM_SMTP_AUTH_MECHANISM="$AM_SMTP_AUTH_MECHANISM" -e AM_DEBUG="$AM_DEBUG" --entrypoint /bin/sh {{ IMAGE }}

test:
    docker exec airmail go test -cover -coverprofile=coverage.out ./...

test-coverage: test
    docker exec airmail go tool cover -func coverage.out | grep total | awk '{print $3}'

up:
    docker-compose up

update:
    docker exec airmail go get -u ./...
