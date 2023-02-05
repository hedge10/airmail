# Airmail

⚠️ This project is under heavy development. Please open an issue for questions and recommendations.


Airmail is a small "form service" to be used for static pages contact forms. It takes your form submission and sends it to a SMTP server.

## Configuration

### Airmail server

You can configure to which IP and Port Airmail should bind the webserver with these environment variables:

- `AM_HOST`, default: `127.0.0.1`
- `AM_PORT`, default: `8081`

### SMTP connection

Set these environment variables to let Airmail connect to your SMTP:

- `AM_SMTP_HOST`, default: `localhost`
- `AM_SMTP_PORT`, default: `25`
- `AM_SMTP_USER`, default: `<empty>`
- `AM_SMTP_PASS`, default: `<empty>`
- `AM_SMTP_AUTH_MECHANISM`, default: `none` (choose between `none`, `plain`, `login`, `cram-md5`, `ntlm`)

