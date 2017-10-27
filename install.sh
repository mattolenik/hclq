#!/usr/bin/env dash
# Installs or upgrades to the latest version of hclq
set -e

main() {
  destination=${1:-"/usr/local/bin"}

  if [ ! -d "$destination" ]; then
    echo "Install directory '$destination' does not exist"
    printUsage
    exit 1
  fi

  if ! command -v jq > /dev/null 2>&1; then
    printf "jq is required for this install script"
    [ "$(uname)" = "Darwin" ] && printf ", it can be installed with 'brew install jq'"
    echo
    exit 1
  fi

  if command -v hclq > /dev/null 2>&1; then
    msg="Upgrading hclq..."
    ver=$(hclq --version)
  else
    msg="Installing hclq..."
  fi

  latest=$(curl -s https://api.github.com/repos/mattolenik/hclq/releases/latest)
  tag=$(echo $latest | jq -r '.tag_name')
  if [ "$tag" = "$ver" ]; then
    echo "hclq is already at the latest version"
    exit 0
  fi

  echo $msg
  hclq_url=$(curl -s https://api.github.com/repos/mattolenik/hclq/releases/latest | jq -r '.assets[] | .browser_download_url' | grep -i "$(uname)")
  curl -sSJLo "$destination/hclq" "$hclq_url" && chmod +x "$destination/hclq"

  echo "hclq now at version $(hclq --version)"
}

printUsage() {
  echo "Usage: install.sh [installDir]"
  echo
  echo "Default installDir is /usr/local/bin"
}

main "$@"
