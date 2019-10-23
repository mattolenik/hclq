#!/usr/bin/env bash
[[ -n ${DEBUG:-} ]] && set -x
set -euo pipefail

include_regex='^[[:blank:]]*include[[:blank:]]+([[:graph:]]+)[[:space:]]*$'
package_regex='^[[:blank:]]*package'

print_go_file() {
  while IFS='' read -r line; do
    if [[ $line =~ $package_regex ]]; then
      continue
    fi
    printf '%s\n' "$line"
  done < "$1"
}

main() {
  while IFS='' read -r line; do
    if [[ $line =~ $include_regex ]]; then
      print_go_file "${BASH_REMATCH[1]}"
      continue
    fi
    printf '%s\n' "$line"
  done
}

main "$@"