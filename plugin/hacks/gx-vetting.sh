#!/bin/bash

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd $CURRENT_DIR/../..

gx-go hook pre-test

echo unused
unused ./...

echo staticcheck
staticcheck ./...

echo gosimple
gosimple ./...

echo golint
golint ./...

gx-go hook post-test
