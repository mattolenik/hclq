[![Build status](https://ci.appveyor.com/api/projects/status/a013fsit6rh0nk93?svg=true)](https://ci.appveyor.com/project/MattOlenik/hclq)

# About

hclq is a command line tool for querying and manipulating [HashiCorp HCL](https://github.com/hashicorp/hcl) files, such as those used by [Terraform](https://terraform.io), [Consul](https://consul.io), [Nomad](https://nomadproject.io), and [Vault](https://vaultproject.io). It's similar to [jq](https://github.com/stedolan/jq), but for HCL.

Use cases include:

 * Performing custom inspection and validation of configuration
 * Enforce rules, naming conventions, etc
 * Preprocessing to compensate for HCL interpolation shortcomings
 * Custom manipulation for wrappers or other utilities
 * A robust alternative to parsing files with grep, sed, etc


# Installation

## Binary Release

Latest releases are available here on Github. hclq is a single binary, installation is as simple as placing the binary in your PATH.

## Install with Go

`go get -u github.com/mattolenik/hclq`

# Help Text

# Contributing

Pull requests are very welcome! If you provide a fix, please try to provide a test case to prevent regressions. Pull requests without tests won't necessarily be rejected (especially for a trivial fix), but tests are helpful and greatly appreciated!

