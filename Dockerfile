FROM axllent/mailpit:latest as mailpit

FROM golang:1.19-alpine

WORKDIR /root

COPY ./dev/docker/openssl.cnf /root/openssl.cnf

# Create self-signed certificate to run SMTP server (mailpit) with TLS
RUN apk add --no-cache direnv openssl \
    && rm -rf /var/cache/apk/*
RUN openssl req -x509 -newkey rsa:4096 -nodes -keyout privkey.pem -out airmail.pem -sha256 -days 365 -config /root/openssl.cnf

RUN cp airmail.pem /usr/local/share/ca-certificates/airmail.crt \
    && update-ca-certificates

COPY --from=mailpit /mailpit /usr/local/bin/mailpit

RUN chmod +x /usr/local/bin/mailpit

COPY ./dev/docker/entrypoint.sh /root/entrypoint.sh
RUN chmod +x /root/entrypoint.sh

WORKDIR /app

ENTRYPOINT ["/root/entrypoint.sh"]

