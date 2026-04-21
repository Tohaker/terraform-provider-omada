# Terraform Provider Omada

Terraform Provider for managing a [TP-Link Omada Software Controller](https://www.tp-link.com/omada-sdn/). Supports v5.x and v6.x.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.25

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install .
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

View the usage documentation in the [Terraform Registry](https://registry.terraform.io/providers/Tohaker/omada/latest/docs).

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

If you use VSCode as your IDE, a `devcontainer` config is provided. This will setup;
1. A fresh Linux container with Go and Terraform installed
2. Terraform configured to override installation of the `Tohaker/omada` provider with the local installation.

To compile the provider, run `go install .`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

```shell
make testacc
```
