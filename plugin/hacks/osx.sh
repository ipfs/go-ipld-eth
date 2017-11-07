#!/bin/bash
set -e

################################################################################
# We get you, for some reason you want to make this work in OSX,
# without the VM hassles. Just run this command, it will:
#
# * Declare this plugin in the preloader at 
#   $GOPATH/src/github.com/ipfs/go-ipfs/plugin/loader/preload_list
# * Alter the `main` package declaration to `plugin` at
#   $GOPATH/src/github.com/ipfs/go-ipld-eth/plugin/eth.go
# * cd into `go-ipfs`, `make-build` and execute `ipfs daemon`.
# * Once the execution of `ipfs daemon` ends, revert the changes.
################################################################################

# Change
echo "ipldeth github.com/ipfs/go-ipld-eth/plugin 0" >> $GOPATH/src/github.com/ipfs/go-ipfs/plugin/loader/preload_list
sed -i '' 's/package main/package plugin/' $GOPATH/src/github.com/ipfs/go-ipld-eth/plugin/eth.go
gx-go rw

# Build and Execute
cd $GOPATH/src/github.com/ipfs/go-ipfs/
make build
./cmd/ipfs/ipfs daemon

# Revert
git checkout -- $GOPATH/src/github.com/ipfs/go-ipfs/plugin/loader/preload_list
git checkout -- $GOPATH/src/github.com/ipfs/go-ipfs/plugin/loader/preload.go
cd $GOPATH/src/github.com/ipfs/go-ipld-eth
gx-go uw
sed -i '' 's/package plugin/package main/' $GOPATH/src/github.com/ipfs/go-ipld-eth/plugin/eth.go
