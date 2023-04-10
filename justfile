set dotenv-load

IMAGE := "hedge10/airmail"

build:
    docker build --no-cache -t {{ IMAGE }} .

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

run:
    docker run --rm -it -p 9900:9900 -e AM_SMTP_HOST="$AM_SMTP_HOST" -e AM_SMTP_PORT="$AM_SMTP_PORT" -e AM_SMTP_USER="$AM_SMTP_USER" -e AM_SMTP_PASS="$AM_SMTP_PASS" -e AM_SMTP_AUTH_MECHANISM="$AM_SMTP_AUTH_MECHANISM" -e AM_DEBUG="$AM_DEBUG" --entrypoint /bin/sh {{ IMAGE }}

test:
    go test -v -cover ./...

up:
    docker-compose up

update:
    go get -u ./...
