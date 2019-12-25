#!/bin/bash

query_term="$1"


BIN="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPOS=$BIN/..

pushd $REPOS

grep -R -i --include=*.go $query_term ./*

popd

