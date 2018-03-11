[![Build Status](https://travis-ci.org/mattolenik/hclq.svg?branch=master)](https://travis-ci.org/mattolenik/hclq)

# hclq

hclq is a command line tool for querying and manipulating [HashiCorp HCL](https://github.com/hashicorp/hcl) files, such as those used by [Terraform](https://terraform.io). It's similar to [jq](https://github.com/stedolan/jq), but for HCL. In addition to retrieving values, hclq can also modify values. By default, hclq returns `get` results as JSON, which can be further processed by jq or some other means. It also has a raw mode for unformatted output.

## Examples



## Installation

hclq is distributed as a single binary for Linux, macOS, and Windows. For now, an automated install and update script is provided. Homebrew support is planned for the future.

`curl -sL https://install.hclq.sh | sh`

The script can take two flags: `-q` for quiet mode, and `-d` to set the installation directory (the default is /usr/local/bin).

The following quietly installs hclq into /usr/bin/
`curl -sL https://install.hclq.sh | sh -s -- -q -d /usr/bin`

And as always, binaries for all platforms are available as GitHub releases.

### Getting Values

### Setting Values
