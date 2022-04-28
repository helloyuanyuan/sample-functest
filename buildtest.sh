#!/bin/bash

env=$1
# token=$2

echo "[Test Env]" ${env}

sed -i "s/\(env: \)\(.*\)/\1${env}/g" config/main.yaml
# sed -i "s/\(token: \)\(.*\)/\1${token}/g" config/main.yaml

go test -v ./functest
go clean -testcache