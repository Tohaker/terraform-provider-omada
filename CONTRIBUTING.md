# Contributing

Thanks for your interest in contributing to the Terraform Provider for Omada! This guide covers the contribution workflow, expectations for tests, the conventions used to structure a service, and how to run the provider locally.

## Contribution process

1. **Fork** the repository to your own GitHub account.
2. **Branch** from `main` in your fork using a short, descriptive name (e.g. `feat/site-tags`, `fix/site-import`).
3. **Implement** your change, following the code structure described below.
4. **Add or update tests** — every change must be covered by automated tests (see
   [Testing requirements](#testing-requirements)).
5. **Run the checks locally** before pushing:
   ```shell
   make fmt
   make lint
   make testacc
   ```
   If you have changed the schema or examples, also run `make generate` to refresh the docs under [docs/](docs/).
6. **Open a pull request** against `Tohaker/terraform-provider-omada:main`. Describe the motivation, the change, and any user-visible behaviour. Link related issues.
7. **Review.** A maintainer will review your PR. Please respond to feedback by pushing additional commits to the same branch (do not force-push during review unless asked). CI must be green before merge.
8. **Merge.** PRs are squash-merged into `main` once approved and CI is green.
9. **Release.** Releases are cut by maintainers from `main` by tagging a new semver version. The tag triggers the release workflow which publishes to the [Terraform Registry](https://registry.terraform.io/providers/Tohaker/omada/latest). Add a brief entry to [CHANGELOG.md](CHANGELOG.md) as part of your PR when the change is user-visible.

## Testing requirements

**All changes must be accompanied by automated tests.** PRs without tests will not be merged.

- New resources, data sources, or schema attributes require acceptance tests under the matching `internal/service/<name>/` package (`resource_test.go`, `data_source_*_test.go`).
- Bug fixes require a regression test that fails without the fix.
- Acceptance tests run against an in-process `httptest` Omada API stand-in provided by [internal/acctest/acctest.go](internal/acctest/acctest.go) — no live controller is needed. Use `acctest.NewTestServer(t)` and register handlers on its `Mux`.
- Run the full suite with `make testacc`. Tests are gated by `TF_ACC=1`, which the make target sets for you.

## Service structure

Each managed object lives in its own package under `internal/service/<name>/`. The site service ([internal/service/site/](internal/service/site/)) is the canonical reference. A service is split into the following files so that schema, API translation, and Terraform plumbing stay separated:

- **`model.go`** — Go structs that mirror the Terraform schema using `terraform-plugin-framework/types` (e.g. `siteResourceModel`, `siteCommonModel`). Also holds small package-internal types like the `siteClient` wrapper around the SDK.
- **`expand.go`** — Functions that **expand** Terraform plan/state models into the request payloads expected by the Omada SDK. One direction only: Terraform → API.
- **`flatten.go`** — Functions that **flatten** Omada API responses back into the Terraform model types. One direction only: API → Terraform.
- **`resource.go`** — The `resource.Resource` implementation: `Metadata`, `Schema`, `Configure`, `Create`, `Read`, `Update`, `Delete`, and (where supported) `ImportState`. This file should be thin — it wires the framework lifecycle to `expand`/`flatten` and the SDK client.
- **`data_source_<name>.go`** — A `datasource.DataSource` implementation following the same pattern as the resource. Use one file per data source (e.g. `data_source_list.go`).
- **`*_test.go`** — Acceptance tests for the corresponding file (`resource_test.go`, `data_source_list_test.go`). Shared test helpers go in `acctest_test.go`.

When adding a new service:

1. Create `internal/service/<name>/` with the files above.
2. Register the resource/data source in
   [internal/provider/provider.go](internal/provider/provider.go).
3. Add an example under `examples/resources/omada_<name>/` or
   `examples/data-sources/omada_<name>/` so `make generate` produces docs.

## Running the provider locally

The repository ships with a VS Code dev container that installs Go, Terraform, and a `~/.terraformrc` with a `dev_overrides` block pointing at `$GOPATH/bin`. If you use it, most of the setup below is already done for you.

### Prerequisites

- [Go](https://golang.org/doc/install) >= 1.25
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0

### Build and install

From the repo root:

```shell
go install .
```

This builds the provider and writes the binary to `$GOPATH/bin` (typically `/go/bin`
in the dev container, or `$(go env GOPATH)/bin` otherwise).

### Configure Terraform to use your local build

Create or edit `~/.terraformrc` so Terraform resolves `Tohaker/omada` from your local
build instead of the registry. Replace `/go/bin` with the directory containing the
binary you just installed:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/Tohaker/omada" = "/go/bin"
  }

  # All other providers continue to install from their origin registry.
  direct {}
}
```

The dev container writes this file automatically via
[.devcontainer/postCreateCommand.sh](.devcontainer/postCreateCommand.sh).

> Note: with `dev_overrides` set, **do not** run `terraform init` for the omada provider — Terraform will warn and use your local binary directly. `terraform plan` and `terraform apply` work as normal.

### Try it against a controller

Point the provider at your Omada controller and run a config from `examples/`:

```hcl
provider "omada" {
  host          = "https://omada.example.com"
  controller_id = "<your-controller-id>"
  client_id     = "<openapi-client-id>"
  client_secret = "<openapi-client-secret>"
}
```

```shell
cd examples/resources/omada_site
terraform plan
```

### Useful make targets

| Target          | What it does                                          |
| --------------- | ----------------------------------------------------- |
| `make build`    | `go build ./...`                                      |
| `make install`  | Build and `go install ./...`                          |
| `make fmt`      | `gofmt -s -w` over the repo                           |
| `make lint`     | Run `golangci-lint`                                   |
| `make test`     | Unit tests                                            |
| `make testacc`  | Full acceptance test suite (`TF_ACC=1`)               |
| `make generate` | Regenerate provider docs under `docs/`                |

### Adding dependencies

```shell
go get github.com/author/dependency
go mod tidy
```

Commit the resulting `go.mod` and `go.sum` changes.
