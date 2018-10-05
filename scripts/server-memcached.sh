#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

my_dir="$(dirname "$0")"

GO_FILES=`find ${my_dir}/../cmd/util/. -type f \( -iname "*.go" ! -iname "*_test.go" \)`
MEMCACHED_URL=localhost:11211 CACHE=memcached go run ${GO_FILES}
