#!/usr/bin/env bash


if [ -d ./resources ]
then
    make
    mkdir ./dist
    cp ./bin/client/gopher ./dist
    cp ./bin/tools/config ./dist
    cp ./bin/tools/echo_server ./dist
    if [ "$1" != "" ]; then
        ./dist/config -input ./application.yml -output ./dist/gopher.$1.yml -stage $1
    else
        ./dist/config -input ./application.yml -output ./dist/gopher.dev.yml -stage dev
    fi
else
    echo "this script must be executed at the project root"
    exit 1
fi



