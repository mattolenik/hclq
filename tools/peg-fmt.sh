#!/usr/bin/env bash
[[ -n ${DEBUG:-} ]] && set -x
set -euo pipefail

fail() { echo "$*" >&2; exit 1; }

NL=$'\n'
EXTENSION=peg

process() {
  local before_code=""
  local after_code=""
  local code=""
  local code_start=false
  local code_end=false
  while IFS='' read -r line; do
    if [[ $line == "{//code"* ]]; then
      if [[ $code_start == true ]]; then
        fail "{//code block declared twice"
      fi
      code_start=true
      before_code+="$line$NL"
      continue
    elif [[ $line == "}//code"* ]]; then
      if [[ $code_start != true ]]; then
        fail "Found ending }//code block before beginning {//code block"
      fi
      code_end=true
      code_start=false
      after_code="$line$NL"
      continue
    fi
    if [[ $code_end == true ]]; then
      after_code+="$line$NL"
    elif [[ $code_start == false ]]; then
      before_code+="$line$NL"
    elif [[ $code_start == true ]]; then
      code+="$line$NL"
    fi
  done
  printf '%s\n' "$before_code"
  printf '%s\n' "$code" | gofmt
  printf '%s\n' "$after_code"
}

main() {
  readarray -t files < <(find . -type d \( -path ./vendor -o -path ./.git \) -prune -o -name "*.$EXTENSION" -print | cut -c3-)
  for file in "${files[@]}"; do
    contents=$(<"$file")
    processed="$(process <<< "$contents")"
    if [[ $contents != "$processed" ]]; then
      printf '%s' "$processed" > "$file"
      echo "$file"
    fi
  done
}

main "$@"