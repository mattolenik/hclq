#!/bin/sh
# Installs or upgrades hclq, by default installing into /usr/local/bin
# This can be overridden with the -d parameter
# Use -q for quiet output
set -e

[ -n "${DEBUG:-}" ] && set -x

E_MISSING_ARG=4
E_MISSING_DIR=5
E_NO_ARCH=6

help() {
  cat <<EOF
Install script for hclq – https://hclq.sh

Options:
         -d <dir>   specify install directory, defaults to /usr/local/bin
         -q         quiet mode, will not print output
         -h         show this help message

EOF
exit 1
}

println() {
  [ -z "$quiet" ] && printf "%s\\n" "$1"
}

platform_check() {
  # OS variable also used to download binary
  OS="$(uname | awk '{print tolower($0)}')"
  ARCH="$(uname -m)"
  case "$ARCH" in
    amd64) ARCH=amd64;;
    x86_64) ARCH=amd64;;
    i386) ARCH=386;;
    i686) ARCH=386;;
    arm) ARCH=arm;;
    *) println "Unsupported or undetected platform: '$ARCH'" && exit $E_NO_ARCH;;
  esac
}

main() {
  platform_check

  while getopts ":qhd:" opt; do
    case $opt in
      q) quiet=true ;;
      d) destination="$OPTARG" ;;
      h) help ;;
      \?) println "Invalid option -$OPTARG" ;;
      :) println "Option -$OPTARG requires an argument" && exit $E_MISSING_ARG
    esac
  done

  destination="${destination:-/usr/local/bin}"
  [ ! -d "$destination" ] && println "Install directory '$destination' does not exist" && exit $E_MISSING_DIR

  if touch "$destination" | grep -q "Permission denied"; then
    println "Permission denied for installing into $destination"
    exit 1
  fi

  # Final binary location
  hclq_bin="$destination/hclq"

  if command -v "$hclq_bin" > /dev/null 2>&1; then
    msg="Upgrading $hclq_bin"
    ver="$($hclq_bin --version)"
  else
    msg="Installing $hclq_bin"
  fi

  # Get latest release info in JSON
  latest="$(curl -s https://api.github.com/repos/mattolenik/hclq/releases/latest)"

  # Get the latest tag
  tag="$(printf '%s' "$latest" | grep tag_name | awk -F'"' '{print $4}')"

  # Bail if the tag isn't new
  [ "$tag" = "$ver" ] && println "$hclq_bin is already at the latest version" && exit 0

  println "$msg"

  # Extract URL for actual binary
  hclq_url=$(printf '%s' "$latest" | grep -i "browser_download_url.*$OS-$ARCH" | awk -F'"' '{print $4}')
  if [ -z "$hclq_url" ]; then
    println "hclq is not available for OS '$OS' on architecture '$ARCH'" && exit $E_NO_ARCH
  fi
  # Only include --silent argument if quiet is defined
  curl ${quiet+--silent} --progress-bar -JLo "$hclq_bin" "$hclq_url"
  chmod +x "$hclq_bin"

  println "$hclq_bin now at version $("$hclq_bin" --version)"
}

main "$@"
