#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

cd $(git rev-parse --show-toplevel)

if ! command -v bazelisk &> /dev/null; then
    echo "Install bazelisk at https://github.com/bazelbuild/bazelisk"
    exit 1
fi

bazelisk run //service/currency:currency_image -- --help
bazelisk run //service/nft:nft_image -- --help
bazelisk run //service/user:user_image -- --help
bazelisk run //service/wordle:wordle_image -- --help
