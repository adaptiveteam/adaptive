#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
echo Using backend $DIR/$ADAPTIVE_CLIENT_ID.tfbackend
terraform init -reconfigure -backend-config=$DIR/$ADAPTIVE_CLIENT_ID.tfbackend
