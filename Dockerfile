FROM golang:1.20-alpine3.17 AS base
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

WORKDIR /root
COPY ./docker/entrypoint.sh ./docker/openssl.cnf ./
RUN chmod +x /root/entrypoint.sh

RUN apk add --no-cache curl direnv openssl\
    && openssl req -x509 -newkey rsa:4096 -nodes -keyout privkey.pem -out airmail.pem -sha256 -days 365 -config /root/openssl.cnf \
    && cp airmail.pem /usr/local/share/ca-certificates/airmail.crt \
    && update-ca-certificates

COPY --from=axllent/mailpit /mailpit /usr/local/bin/mailpit
RUN chmod +x /usr/local/bin/mailpit

# Install AIR
WORKDIR /
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

WORKDIR /app

ENTRYPOINT ["/bin/bash", "-c", "/bin/air"]


