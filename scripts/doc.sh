#!/usr/bin/env bash

set -euo pipefail

godoc2md github.com/andy2046/parallel \
    > $GOPATH/src/github.com/andy2046/parallel/parallel.md
