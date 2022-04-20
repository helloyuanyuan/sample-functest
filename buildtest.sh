#!/bin/bash

env=$1

echo "[Test Env]" ${env}

sed -i "s/\(env: \)\(.*\)/\1${env}/g" config/profile.yaml

export GOPROXY=https://goproxy.cn

go test -v ./functest