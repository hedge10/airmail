<a href="https://github.com/hedge10/airmail/actions/workflows/release.yml" target="_blank" rel="noopener"><img src="https://github.com/hedge10/airmail/actions/workflows/release.yml/badge.svg" alt="Release" /></a>
<a href="https://github.com/hedge10/airmail/actions/workflows/tests.yml" target="_blank" rel="noopener"><img src="https://github.com/hedge10/airmail/actions/workflows/tests.yml/badge.svg" alt="Tests" /></a>
[![Go Report Card](https://goreportcard.com/badge/github.com/hedge10/airmail)](https://goreportcard.com/report/github.com/hedge10/airmail)

# Airmail

![Airmail](./art/airmail-scheme.svg "Title")

## Development

Prerequite(s):

-   [just](https://github.com/casey/just) command runner

🗒️ You can use the development setup without using `just`. Have a look into the [`justfile`](./justfile) for individual commands.

Run `just up` in the root folder to spin up the following Docker container:

-   Airmail: the code is mirrored inside and live-reloaded via [Air](https://github.com/cosmtrek/air). Reach _Airmail_ via `http://localhost:9900`. Additionally this container also provides the local SMTP server [Mailpit](https://github.com/axllent/mailpit) reachable via `http://localhost:8025`

Happy coding! 🕺🏻💃🏻

## Contributing

Thank you for considering contributing. Please have a look at the [contributing guide](CONTRIBUTING.md).
