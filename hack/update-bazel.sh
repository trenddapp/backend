#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

cd $(git rev-parse --show-toplevel)

if ! command -v bazelisk &> /dev/null; then
    echo "Install bazelisk at https://github.com/bazelbuild/bazelisk"
    exit 1
fi


bazelisk run //:gazelle
bazelisk run //:gazelle-update-repos
bazelisk run //:gazelle
bazelisk build //...
