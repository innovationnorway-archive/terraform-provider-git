# Terraform Provider for Git

[2022.06.10] Fork of `innovationnorway/terraform-provider-git` to have a proper Darwin/arm build. Must be switched back to upstream after [Support for darwin_arm64](https://github.com/innovationnorway/terraform-provider-git/issues/69) would resolved

[2022.06.16] Add fetchure: searching git repo in CWD and recursively above it

![](https://github.com/innovationnorway/terraform-provider-git/workflows/test/badge.svg)

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.12.x
-	[Go](https://golang.org/doc/install) >= 1.14

## Usage

```hcl
provider "git" {}

data "git_repository" "example" {
  path = path.root
}

resource "azurerm_resource_group" "example" {
  ...

  tags = {
    branch = data.git_repository.example.branch
    commit = substr(data.git_repository.example.commit_sha, 0, 7)
    tag    = data.git_repository.example.tag
  }
}
```

## Contributing

To build the provider:

```sh
$ go build
```

To test the provider:

```sh
$ go test -v ./...
```

To run all acceptance tests:

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ TF_ACC=1 go test -v ./...
```

To run a subset of acceptance tests:

```sh
$ TF_ACC=1 go test -v ./... -run=TestAccDataSourceGitRepository
```
