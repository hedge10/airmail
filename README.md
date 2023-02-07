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



## Contributing

Thank you for considering contributing. Please have a look at the following sections to help you setting up a local environment.

### Development environment

Run `docker-compose up -d` in the root folder to build a local SMTP server and providing a small container running Airmail.

The SMTP server is based on [Mailpit](https://github.com/axllent/mailpit), running the server on port `1025` and the web interface on `http://localhost:8025`. Make sure these ports are not already in use by your host system.

After the container is built and up and running, run `docker exec -it airmail sh` to give you a shell into the

