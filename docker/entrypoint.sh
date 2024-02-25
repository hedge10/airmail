#!/bin/sh

set -e

# Start Mailpit SMTP server
/usr/local/bin/mailpit --smtp-tls-cert /usr/local/share/ca-certificates/airmail.crt --smtp-tls-key /certs/privkey.pem --smtp-auth-accept-any &

exec "$@"
