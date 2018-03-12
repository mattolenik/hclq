[![Build Status](https://travis-ci.org/mattolenik/hclq.svg?branch=master)](https://travis-ci.org/mattolenik/hclq)

# About

hclq is a command line tool for querying and manipulating [HashiCorp HCL](https://github.com/hashicorp/hcl) files, such as those used by [Terraform](https://terraform.io) and other. It's similar to [jq](https://github.com/stedolan/jq), but for HCL. It can also modify HCL, with the option of modifying files in-place.

Use cases include:

 * Performing custom inspection and validation of configuration
 * Enforce rules, naming conventions, etc
 * Preprocessing to compensate for HCL interpolation shortcomings
 * Custom manipulation for wrappers or other utilities
 * A robust alternative to parsing files with grep, sed, etc

hclq outputs JSON for easy processing with other tools.

Note: hclq should be considered alpha software at this time.

# Installation

Install with auto-updating script:

`curl -sL https://install.hclq.sh | sh`

Install with Go:

`go get -u github.com/mattolenik/hclq`

Or download a release from GitHub. More info on the script is down below.

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

## Installation Details

hclq is distributed as a single binary for Linux, macOS, and Windows. For now, an automated install and update script is provided. Homebrew support is planned for the future.

As mentioned above, the following will automatically install or update hclq:

`curl -sL https://install.hclq.sh | sh`

The script can take two flags: `-q` for quiet mode, and `-d` to set the installation directory (the default is /usr/local/bin).

The following quietly installs hclq into /usr/bin/

`curl -sL https://install.hclq.sh | sh -s -- -q -d /usr/bin`


## Project Status

hclq should be considered alpha software and not suitable for production. If you really feel the need for it, snap to a specific release and be sure to test your use cases thoroughly.

There's still trouble with some of the list processing, and full objects are not yet supported. More tests are needed, especially for the set commands.
