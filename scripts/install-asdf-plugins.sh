#!/usr/bin/env bash

# This script will install the ASDF plugins required for this project

set -euo pipefail
IFS=$'\n\t'

plugin_list=$(asdf plugin list || echo "")

install_plugin() {
  plugin=$1
  if ! echo "$plugin_list" | grep -q "${plugin}" >/dev/null; then
    echo "# Installing plugin" "$@"
    asdf plugin add "$@" || {
      echo "Failed to perform plugin installation: " "$@"
      exit 1
    }
  fi

  echo "# Installing ${plugin} version"
  asdf install "${plugin}" || {
    echo "Failed to perform version installation: ${plugin}"
    exit 1
  }
}

# Install golang first as some of the other plugins require it
install_plugin golang
install_plugin goreleaser
install_plugin golangci-lint
