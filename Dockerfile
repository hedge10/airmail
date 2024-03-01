FROM golang:1.21-alpine3.17 AS base
ARG VERSION=SNAPSHOT
ENV CGO_ENABLED=0

WORKDIR /app
COPY . /app/

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'github.com/hedge10/airmail.Version=$VERSION'" -o airmail ./cmd/airmail

FROM gcr.io/distroless/static AS prod

COPY --from=base /app/airmail .

EXPOSE 9900

ENTRYPOINT [ "/airmail" ]
