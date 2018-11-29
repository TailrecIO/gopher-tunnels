#!/usr/bin/env bash

if [ -d ./resources ]
then
    make
    mkdir ./dist
    cp ./bin/client/gopher ./dist
    cp ./bin/tools/config ./dist
    cp ./bin/tools/http_echo ./dist
    ./dist/config -input ./application.yml -output ./dist/gopher.yml
else
    echo "this script must be executed at the project root"
    exit 1
fi



