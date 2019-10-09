#!/usr/bin/env bash
[[ -n ${DEBUG:-} ]] && set -x
set -euo pipefail

fail() { echo "$*" >&2; exit 1; }

NL=$'\n'
EXTENSION=peg
START_BLOCK='{#code'
END_BLOCK='}#code'

process() {
  local before_code=""
  local after_code=""
  local code=""
  local code_start=false
  local code_end=false
  while IFS='' read -r line; do
    if [[ $line == "$START_BLOCK"* ]]; then
      if [[ $code_start == true ]]; then
        fail "$START_BLOCK block declared twice"
      fi
      code_start=true
      before_code+="$line$NL"
      continue
    elif [[ $line == "$END_BLOCK"* ]]; then
      if [[ $code_start != true ]]; then
        fail "Found ending $END_BLOCK block before beginning $START_BLOCK block"
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
  if ! command -v gofmt &> /dev/null; then
    fail "gofmt not found, Go must be installed for $(basename $0) to work"
  fi
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