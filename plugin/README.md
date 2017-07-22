# Ethereum ipld plugin

Add ethereum support to ipfs!

## Building
Make sure to rewrite the gx dependencies in the directory above first, then
Either run `make` or `go build -buildmode=plugin -o=ethereum.so`.

* *NOTE*: As of 2017.07.17 Plugins in Go only works in Linux.

## Installing
Move `ethereum.so` to `$IPFS_PATH/plugins/ethereum.so` and set it to be executable:

```sh
mkdir -p ~/.ipfs/plugins
mv ethereum.so ~/.ipfs/plugins/
chmod +x ~/.ipfs/plugins/ethereum.so
```

## I don't have linux but I want to do this somehow!

As stated above, the _plugin_ library only works in Linux. Bug the go team to support your system.

... Or use a linux virtualbox, and mount this directory.

... Or hack your way via docker-fu [with this short, unsupported guide](hacks/docker.md)

... Or, if you are in OSX, [use this handy script](hacks/osx.sh)

## Usage

Make sure you have the right version of ipfs installed (plugin support was
merged a few hours before this was written) and start up the ipfs daemon!

There are a couple test files in this repo, so try out:
```
cat ../test_data/eth-block-4052384 | ipfs dag put --input-enc=raw --format=eth-block

# Should return this hash
z43AaGEwLuiGeeRYszh2ZVtAe92HK796zn8Qz5REq7ztM1ZBz7d
```

Then take that hash and explore it with the `ipfs dag get` command:

```
ipfs dag get z43AaGEwLuiGeeRYszh2ZVtAe92HK796zn8Qz5REq7ztM1ZBz7d
```

Should get

```
<TODO>
```
