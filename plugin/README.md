# Ethereum IPLD plugin

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

### I don't have linux but I want to do this somehow!

As stated above, the _plugin_ library only works in Linux. Bug the go team to support your system!

* ... Or use a linux virtualbox, and mount this directory.

* ... Or hack your way via docker-fu [with this short, unsupported guide](hacks/docker.md)

* ... Or, if you are in OSX, [use this handy script](hacks/osx.sh)

## Usage and Examples

(TODO)
(Note that there are two kind of inputs for this problem: json and raw)

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
{"bloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000400080000000000000000000000000000000280000004000000000000000000000000000000000080000000000000000000000000000000000000204000200000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020400000000000200000000000000000000000000010000000000000000000000000100000000000000000000000000000000000000000000400000000000080000000000000000000","coinbase":"0xea674fdde714fd979de3edf0f56aa9716b898ec8","difficulty":1283310643319628,"extra":"ZXRoZXJtaW5lLWV1NQ==","gaslimit":6719052,"gasused":562490,"mixdigest":"0x1872a178541e5e263a1a68d797a72baa0bcbac1500b29d51491c586f17a0fab2","nonce":"0x1ac84bc00f34b563","number":4052384,"parent":{"/":"z43AaGEs3RrhkLLQkRuHt3e8Ms9NGY7Matigoxj45Hs1BzDHBp3"},"receipts":{"/":"z44vkPhh97phomDczCHJk4CWqwDLL6YHQUPdgz8uEi8R3LLj6XF"},"root":{"/":"z45oqTRvsWb3UyAwG3o1uJ9zg2v76DQuc3YJEyMPSq3fJ5fZdxe"},"time":1500627255,"tx":{"/":"z443fKyTGWGtvy3M9nyRWM3zwhEwBazCWYjLpYNV5HYKJUspSFH"},"uncles":{"/":"z43c7o74hjCAqnyneWetkyXU2i5KuGQLbYfVWZMvJMG4VTYABtz"}}
```
