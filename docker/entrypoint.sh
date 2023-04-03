#!/bin/sh

set -e

# Start Maildrip SMTP server
/usr/local/bin/mailpit --smtp-ssl-cert /usr/local/share/ca-certificates/airmail.crt --smtp-ssl-key /root/privkey.pem
