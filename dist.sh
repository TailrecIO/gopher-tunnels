#!/usr/bin/env bash


if [ -d ./resources ]
then
    make
    DIST_DIR="./dist"
    if [ "$1" != "" ]; then
        DIST_DIR="./dist.$1"
    fi
    echo "Creating a distribution package in $DIST_DIR ..."
    mkdir $DIST_DIR
    cp ./bin/client/gopher $DIST_DIR
    cp ./bin/tools/config $DIST_DIR
    cp ./bin/tools/echo_server $DIST_DIR
    if [ "$1" != "" ]; then
        $DIST_DIR/config -input ./application.yml -output $DIST_DIR/gopher.$1.yml -stage $1
    else
        $DIST_DIR/config -input ./application.yml -output $DIST_DIR/gopher.yml
    fi
else
    echo "this script must be executed at the project root"
    exit 1
fi



