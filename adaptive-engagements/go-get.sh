#!/bin/bash
file="$1"
while IFS="=" read -r dep rev; do
  case "$dep" in
  '#'*) ;;
  *)
    go get "${dep}@${rev}"
    echo "go get ${dep}@${rev} done"
    ;;
  esac
done <'deps.txt'
