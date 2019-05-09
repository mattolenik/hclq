#!/bin/sh
# Installs or upgrades hclq, by default installing into /usr/local/bin
# This can be overridden with the -d parameter, see help()
set -e
[ -n "${DEBUG:-}" ] && set -x

REPO="mattolenik/hclq"

help() {
  cat <<EOF
Install script for hclq – https://hclq.sh

Options:
         -d <dir>    specify install directory, defaults to /usr/local/bin
         -a <GOARCH> set specific architecture, values correspond to GOARCH values
         -o <GOOS>   set specific OS, values correspond to GOOS values
         -q          quiet mode, will not print output
         -h          show this help message

EOF
}

println() {
  [ -z "$QUIET" ] && printf "%s\\n" "$*"
}

fail() {
  [ -z "$QUIET" ] && printf "%s\\n" "$*" 1>&2
  exit 1
}

platform_check() {
  # OS and ARCH variables also used to download binary in main function
  OS="$(uname | awk '{print tolower($0)}')"
  ARCH="$(uname -m)"
  case "$ARCH" in
    amd64) ARCH=amd64;;
    x86_64) ARCH=amd64;;
    i386) ARCH=386;;
    i686) ARCH=386;;
    arm) ARCH=arm;;
    *) fail "Unsupported or undetected platform: '$ARCH'";;
  esac
}

main() {
  platform_check

  while getopts ":qhdao:" opt; do
    case $opt in
      q) QUIET=true ;;
      d) DESTINATION="$OPTARG" ;;
      a) ARCH="$OPTARG" ;;
      o) OS="$OPTARG" ;;
      h) help && exit 0 ;;
      \?) println "Invalid option -$opt" ;;
      :) fail "Option -$opt requires an argument";;
    esac
  done

  DESTINATION="${DESTINATION:-/usr/local/bin}"
  [ ! -d "$DESTINATION" ] && fail "Install directory '$DESTINATION' does not exist"

  SUDO_CMD=""
  if touch "$DESTINATION" 2>&1 | grep -q "Permission denied"; then
    SUDO_CMD=sudo
  fi

  # Final binary location
  hclq_bin="$DESTINATION/hclq"

  if command -v "$hclq_bin" > /dev/null 2>&1; then
    msg="Upgrading $hclq_bin"
    ver="$($hclq_bin --version)"
  else
    msg="Installing $hclq_bin"
  fi

  # Get latest release info in JSON
  latest="$(curl -s https://api.github.com/repos/$REPO/releases/latest)"

  # Get the latest tag
  tag="$(printf '%s' "$latest" | grep tag_name | awk -F'"' '{print $4}')"

  # Bail if the tag isn't new
  [ "$tag" = "$ver" ] && println "$hclq_bin is already at the latest version" && exit 0

  println "$msg"

  # Extract URL for actual binary
  hclq_url=$(printf '%s' "$latest" | grep -i "browser_download_url.*$OS-$ARCH" | awk -F'"' '{print $4}')
  if [ -z "$hclq_url" ]; then
    fail "hclq is not available for OS '$OS' on architecture '$ARCH'"
  fi
  tmp_bin="$(mktemp)"
  trap 'rm -f $tmp_bin' EXIT
  # Only include --silent argument if QUIET is defined
  curl ${QUIET+--silent} --progress-bar -JLo "$tmp_bin" "$hclq_url"
  chmod +x "$tmp_bin"
  $SUDO_CMD mv -f "$tmp_bin" "$hclq_bin"

  println "$hclq_bin now at version $("$hclq_bin" --version)"
}

main "$@"
