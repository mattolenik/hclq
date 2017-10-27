#!/bin/sh
# Installs or upgrades hclq, by default installing into /usr/local/bin
# This can be overridden by passing a directory as the first parameter.
set -e

destination=${1:-"/usr/local/bin"}
[ ! -d "$destination" ] && printf "Install directory '%s' does not exist\n", "$destination" && exit 1

if command -v hclq > /dev/null 2>&1; then
  msg="Upgrading hclq..."
  ver="$(hclq --version)"
else
  msg="Installing hclq..."
fi

latest="$(curl -s https://api.github.com/repos/mattolenik/hclq/releases/latest)"
tag="$(printf "$latest" | grep tag_name | awk -F'"' '{print $4}')"
[ "$tag" = "$ver" ] && printf "hclq is already at the latest version\n" && exit 0
printf "%s\n" "$msg"

hclq_url=$(curl -s https://api.github.com/repos/mattolenik/hclq/releases/latest | grep -i "browser_download_url.*$(uname)" | awk -F'"' '{print $4}')
curl --progress-bar -JLo "$destination/hclq" "$hclq_url" && chmod +x "$destination/hclq"

printf "hclq now at version %s\n" "$(hclq --version)"
