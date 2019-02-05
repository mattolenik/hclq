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
