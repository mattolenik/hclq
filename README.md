[![Build status](https://ci.appveyor.com/api/projects/status/a013fsit6rh0nk93?svg=true)](https://ci.appveyor.com/project/MattOlenik/hclq)

# About

hclq is a command line tool for querying and manipulating [HashiCorp HCL](https://github.com/hashicorp/hcl) files, such as those used by [Terraform](https://terraform.io) and other. It's similar to [jq](https://github.com/stedolan/jq), but for HCL. It can also modify HCL, with the option of modifying files in-place.

Use cases include:

 * Performing custom inspection and validation of configuration
 * Enforce rules, naming conventions, etc
 * Preprocessing to compensate for HCL interpolation shortcomings
 * Custom manipulation for wrappers or other utilities
 * A robust alternative to parsing files with grep, sed, etc

hclq outputs JSON for easy processing with other tools.

Note: HCL2 support (as used in Terraform 0.12) is planned for this project, see issue [#11](/../../issues/11)

# Installation

## Binary Release

Latest releases are available here on Github. hclq is a single binary, installation is as simple as placing the binary in your PATH.


# Auto-Updating Script

Install with auto-updating script. You may be prompted with `sudo` if your user doesn't have permission to write to the default install directory. Either enter your password or re-run with the `-d` option to specify a different install directory.

```
# Download the script
curl -sSLo install.sh https://install.hclq.sh

# Inspect the script before executing, as is good practice for code from the internet, regardless of trust!
less install.sh

# Run it!
sh install.sh
```

Or just download `install.sh` right from the repo and run that.

## Install with Go

`go get -u github.com/mattolenik/hclq`

Keep in mind this will install whatever is on the `master` branch. A stable release is recommended unless you are doing development.

# Help Text
```
hclq is a tool for manipulating the config files used by HashiCorp tools.

hclq uses a "breadcrumb" or "path" style query. Given the HCL:
    data "foo" "bar" {
        id = "100"
        other = [1, 2, 3]
    }

A query for 'data.foo.bar.id' would return 100. Arrays/lists must be matched
with the [] suffix, e.g. 'data.foo.bar.other[]' or 'data.foo.bar.other[1]'.

Match types:
    literal     Match a literal value.
    list[]      Match a list and retrieve all items.
    list[1]     Match a list and retrieve a specific item.
    /regex/     Match anything according to the specified regex.
    /regex/[]   Match a list according to the regex and retrieve all items.
    /regex/[1]  Match a list according to the regex and retrieve a specific item.
    *           Match anything.

Queries can return either single or multiple values. If a query matches e.g.
multiple arrays across multiple objects, a list of arrays will be returned.
If this query is used with a set command, ALL of those matching arrays will be
set.

Usage:
  hclq [command]

Available Commands:
  get         retrieve values matching <query>
  help        Help about any command
  set         set matching value(s), specify a string, number, or JSON object or array

Flags:
  -h, --help         help for hclq
  -i, --in string    read input from this file, otherwise use stdin
  -o, --out string   write output to this file, otherwise use stdout
      --version      version for hclq

Use "hclq [command] --help" for more information about a command.
```

# Examples

## Getting Values

Let's try getting and setting some values from a simple document, we'll refer to it as `example.tf`. All output is JSON.

```sh
data "foo" {
  bin = [1, 2, 3]

  bar = "foo string"
}

data "baz" {
  bin = [4, 5, 6]

  bar = "baz string"
}
```

Let's say we want the value of `bar` from the `foo` object. The result is printed as a plain string.

```sh
$ cat example.tf | hclq get 'data.foo.bar'

"foo string"
```

### Getting Lists

For a list, append `[]` to tell hclq to look for lists, otherwise it'll try to find single items:

```sh
$ cat example.tf | hclq get 'data.foo.bin[]'

[1,2,3]
```

### Wildcard Matching

hclq can match across many objects and return combined results. Use the wildcard `*` to match any key:

```sh
$ cat example.tf | hclq get 'data.*.bar'

["foo string","baz string"]
```

## Setting Values

hclq can also edit HCL documents. Simply form a query just like with `get`, but also provide a new value. Anything that matches the query will get set to the new value. In other words, anything that would be returned by `get`, will also be affected by `set`. For example:

```sh
$ cat example.tf | hclq set 'data.*.bar' "new string"

data "foo" {
  bin = [1, 2, 3]

  bar = "new string"
}

data "baz" {
  bin = [4, 5, 6]

  bar = "new string"
}
```

#### Setting Lists

It works on lists, too:

```sh
$ cat example.tf | hclq set 'data.*.bin[]' "[10, 11]"

data "foo" {
  bin = [
    10,
    11,
  ]

  bar = "foo string"
}

data "baz" {
  bin = [
    10,
    11,
  ]

  bar = "baz string"
}
```

## Installation Script

The script `install.sh` (which is served at https://install.hclq.sh) can take the following options:

```
Install script for hclq – https://hclq.sh

Options:
         -d <dir>    specify install directory, defaults to /usr/local/bin
         -a <GOARCH> set specific architecture, values correspond to GOARCH values
         -o <GOOS>   set specific OS, values correspond to GOOS values
         -q          quiet mode, will not print output
         -h          show this help message

```

If the DEBUG environment variable is set, the script will enable `set -x` to aid in debugging.

## Project Status

hclq can be said to be beta software. All or almost all of the features are present, with the most common use cases being supported.
Changes to the CLI interface may occur but the general way the CLI works should remain the same. If you run hclq in a production or near-production environment, you should pin your version of hclq to one that you have tested within your system.

Whether or not hclq should be used "in production" depends on the criticality of its application and is up to your judgement.
As with any unfamiliar open source project, inspect its source and test cases when evaluating if it meets your production standards.

# Contributing

Pull requests are very welcome! If you provide a fix, please try to provide a test case to prevent regressions. Pull requests without tests won't necessarily be rejected (especially for a trivial fix), but tests are helpful and greatly appreciated!
