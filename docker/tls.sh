#!/bin/sh

set -e

echo "Generate airmail TLS certificate start... "
mkdir certs

openssl req -x509 -newkey rsa:4096 -nodes -keyout certs/privkey.pem -out certs/airmail.pem -sha256 -days 365 -config "$1" 2>/dev/null

if [ -z "$GITHUB_WORKSPACE" ]; then
    echo "Copy certificate..."
    mv certs/airmail.pem /usr/local/share/ca-certificates/airmail.crt
    update-ca-certificates
fi

echo "Generate airmail TLS certificate end"
