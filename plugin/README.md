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
cat ../test_data/block-with-txs.bin | ipfs dag put --input-enc=raw --format=eth
```

Then take that hash and explore it with the `ipfs dag get` command:

```
> ipfs dag get z43AaGF4uHSY4waU68L3DLUKHZP7yfZoo6QbLmid5HomZ4WtbWw
{"bloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","coinbase":"0x52bc44d5378309ee2abf1539bf71de1b7d7be3b5","difficulty":12555463106190,"extra":"14MBAwOER2V0aIdnbzEuNC4yhWxpbnV4","gaslimit":3141592,"gasused":231000,"mixdigest":"0x5b10f4a08a6c209d426f6158bd24b574f4f7b7aa0099c67c14a1f693b4dd04d0","nonce":"0xf491f46b60fe04b3","number":999999,"parent":{"/":"z43AaGF6wP6uoLFEauru5oLK5JS5MGfNuGDK1xWEpQK4BqkJkL3"},"receipts":{"/":"z44vkPhhDSTXPAswvC1rdDunzkgZ7FgAAnhGQtNDNDk9m9N2BZA"},"root":{"/":"z45oqTSAZvPiiPV8hMZDH5fi4NkaAkMYTJC6PmaeWBmYUpbMpoh"},"time":1455404037,"tx":{"/":"z443fKyHHMwVy13VXtD4fdRcUXSqkr79Q5E8hcmEravVBq3Dc51"},"uncles":{"/":"z43c7o74hjCAqnyneWetkyXU2i5KuGQLbYfVWZMvJMG4VTYABtz"}}
```
