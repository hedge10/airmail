FROM golang:1.21-alpine3.17 AS base
ARG VERSION=SNAPSHOT
ENV CGO_ENABLED=0

WORKDIR /app

COPY . /app/

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'github.com/hedge10/airmail.Version=$VERSION'" -o airmail ./cmd/airmail

FROM gcr.io/distroless/static AS prod
LABEL org.opencontainers.image.source="https://github.com/hedge10/airmail"

COPY --from=base /app/airmail .

EXPOSE 9900

ENTRYPOINT [ "/airmail" ]

FROM base AS dev

WORKDIR /
COPY ./docker/entrypoint.sh ./docker/tls.sh ./docker/openssl.cnf ./
RUN chmod +x ./entrypoint.sh ./tls.sh

RUN apk add --no-cache curl direnv openssl

RUN /bin/sh tls.sh /openssl.cnf

COPY --from=axllent/mailpit /mailpit /usr/local/bin/mailpit
RUN chmod +x /usr/local/bin/mailpit

COPY --from=golangci/golangci-lint /bin/golangci-lint /usr/local/bin/golangci-lint
RUN chmod +x /usr/local/bin/golangci-lint

# Install AIR
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

WORKDIR /app

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/bin/bash", "-c", "/bin/air"]

