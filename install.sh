#!/bin/sh
# Installs or upgrades hclq, by default installing into /usr/local/bin
# This can be overridden by passing a directory as the first parameter.
set -e

destination=${1:-"/usr/local/bin"}
[ ! -d "$destination" ] && printf "Install directory '%s' does not exist\n", "$destination" && exit 1

if ! command -v jq > /dev/null 2>&1; then
  printf "jq is required for this install script"
  [ "$(uname)" = "Darwin" ] && printf ", it can be installed with 'brew install jq\n'"
  exit 1
fi

if command -v hclq > /dev/null 2>&1; then
  msg="Upgrading hclq..."
  ver="$(hclq --version)"
else
  msg="Installing hclq..."
fi

latest="$(curl -s https://api.github.com/repos/mattolenik/hclq/releases/latest)"
tag="$(printf "$latest" | jq -r '.tag_name')"
[ "$tag" = "$ver" ] && printf "hclq is already at the latest version\n" && exit 0
printf "%s\n" "$msg"

hclq_url=$(curl -s https://api.github.com/repos/mattolenik/hclq/releases/latest | jq -r '.assets[] | .browser_download_url' | grep -i "$(uname)")
curl --progress-bar -JLo "$destination/hclq" "$hclq_url" && chmod +x "$destination/hclq"

printf "hclq now at version %s\n" "$(hclq --version)"
