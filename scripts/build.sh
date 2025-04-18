#!/bin/bash

# Needed env vars:
# TAG - SemVer tag
# GITHUB_SHA - Reversion number

set -eu

# Ensure Go flags are unset to avoid usage of go mod vendor
export GOFLAGS=""
go mod download

IFS=" "
GOARCH=amd64
CGO_ENABLED=0
OSLIST="linux" # So far no other OS is supported ;)

mkdir -p build

version=$TAG

if [ "$version" = "null" ] || [ "$version" = "" ]; then
  echo "WARN: No tags found on current commit: $GITHUB_SHA. Using commit hash instead."
  version="$GITHUB_SHA"
fi

go_ldflags="-X github.com/Koenigseder/badm/cmd.version=$version -X github.com/Koenigseder/badm/cmd.revision=$GITHUB_SHA"

for GOOS in $OSLIST; do
  export GOARCH
  export GOOS
  export CGO_ENABLED

  echo "Build BADM ($GOOS)"
  output_file="build/badm_${GOOS}_${GOARCH}"
  go build -o "badm" -ldflags "$go_ldflags" .
  tar cvfz "${output_file}.tar.gz" "badm"
done
