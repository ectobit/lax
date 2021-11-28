# lax

[![Build Status](https://github.com/ectobit/lax/workflows/build/badge.svg)](https://github.com/ectobit/lax/actions)
[![Go Reference](https://pkg.go.dev/badge/go.ectobit.com/lax.svg)](https://pkg.go.dev/go.ectobit.com/lax)
[![Go Report](https://goreportcard.com/badge/go.ectobit.com/lax)](https://goreportcard.com/report/go.ectobit.com/lax)
[![License](https://img.shields.io/badge/license-BSD--2--Clause--Patent-orange.svg)](https://github.com/ectobit/lax/blob/main/LICENSE)

Logger abstraction in Go.

Besides implementing own logger interface, all adapters also implement [pgx.Logger](https://github.com/jackc/pgx/blob/master/logger.go) interface and may implement more interfaces in the future. This eliminates usage of multiple adapters.

## Adapters

- [zap](https://github.com/uber-go/zap)

## Contribution

- `make lint` lints the project
- `make test` runs unit tests
- `make test-cov` displays test coverage (requires docker-stack to be up)
