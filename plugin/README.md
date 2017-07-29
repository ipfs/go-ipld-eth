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

Make sure you have the right version of ipfs installed and start up the ipfs daemon!

### Add an ethereum block written in JSON

You may want to take a block given by your favorite client's JSON RPC API.
We have a couple of those in the `test-data` directory.

```
cat ./test_data/eth-block-body-json-997522 | ipfs dag put --input-enc json --format eth-block
```

And get the CID of the block header back!

```
z43AaGEzuAXhWf9pWAm63QCERtFpqcc6gQX3QBBNaG1syxGGhg6
```

Now, you can get this block header

```
ipfs dag get z43AaGEzuAXhWf9pWAm63QCERtFpqcc6gQX3QBBNaG1syxGGhg6
```

Which will get you (with the right IPLD formatted for the other objects)

```
{"bloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","coinbase":"0x4bb96091ee9d802ed039c4d1a5f6216f90f81b01","difficulty":11966502474733,"extra":"14MBBACER2V0aIdnbzEuNS4xhWxpbnV4","gaslimit":3141592,"gasused":21000,"mixdigest":"0x2565992ba4dbd7ab3bb08d1da34051ae1d90c79bc637a21aa2f51f6380bf5f6a","nonce":"0xf7a14147c2320b2d","number":997522,"parent":{"/":"z43AaGF24mjRxbn7A13gec2PjF5XZ1WXXCyhKCyxzYVBcxp3JuG"},"receipts":{"/":"z44vkPhjt2DpRokuesTzi6BKDriQKFEwe4Pvm6HLAK3YWiHDzrR"},"root":{"/":"z45oqTRunK259j6Te1e3FsB27RJfDJop4XgbAbY39rwLmfoVWX4"},"time":1455362245,"tx":{"/":"z443fKyLvyDQBBQRGMNnPb8oPhPerbdwUX2QsQCUKqte1hy4kwD"},"uncles":{"/":"z43c7o73GVAMgEbpaNnaruD3ZbF4T2bqHZgFfyWqCejibzvJk41"}}
```

#### Piping from the RPC

The astute reader will say "_Let's then pipe the output of my RPC directly to IPFS!_"

```
curl -s -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x1b4", true],"id":1}' https://mainnet.infura.io | ipfs dag put --input-enc json --format eth-block && echo
```

Will give you

```
z43AaGF7XiKhgVVcYxNJv3ZrebEkDE5yhna22N74AusBdMvi6pV
```


And call

```
ipfs dag get z43AaGF7XiKhgVVcYxNJv3ZrebEkDE5yhna22N74AusBdMvi6pV
```

To get

```
{"bloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","coinbase":"0xbb7b8287f3f0a933474a79eae42cbca977791171","difficulty":21109876668,"extra":"R2V0aC9MVklWL3YxLjAuMC9saW51eC9nbzEuNC4y","gaslimit":5000,"gasused":0,"mixdigest":"0x4fffe9ae21f1c9e15207b1f472d5bbdd68c9595d461666602f2be20daf5e7843","nonce":"0x689056015818adbe","number":436,"parent":{"/":"z43AaGF8SkCtKoht2v1e3yC9DWHi4iV2dynyi3BTCP7sPs7HR2T"},"receipts":{"/":"z44vkPheUUg5HBpxkq5sFFz5d9ckigtBBW7WCJXQSZA1gV233Ap"},"root":{"/":"z45oqTS9WCLjMeLnFvTbWiqxXRi1PdwYtDjnNQy6PyWKokGD8r8"},"time":1438271100,"tx":{"/":"z443fKyJXGFJKPgzhha8eqpkEz3rHUL5M7cvcfJQVGzwt3MwcVn"},"uncles":{"/":"z43c7o74hjCAqnyneWetkyXU2i5KuGQLbYfVWZMvJMG4VTYABtz"}}
```

Or go even more _extreme_ with a single pipe

```
ipfs dag get $(curl -s -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x2DC6C0", true],"id":1}' https://mainnet.infura.io | ipfs dag put --input-enc json --format eth-block)

{"bloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","coinbase":"0xea674fdde714fd979de3edf0f56aa9716b898ec8","difficulty":103975266902792,"extra":"ZXRoZXJtaW5lIC0gRVUy","gaslimit":3996095,"gasused":269381,"mixdigest":"0xb1907b36cfc58e666ec1f2d2b60422fc222b0994739bfe0a4b10ba68960cf2ab","nonce":"0xe9d09233833686d4","number":3000000,"parent":{"/":"z43AaGEzw2GGeV87BKohV7ykpq6FssJjBagzmGufRvVzq2Zvz5o"},"receipts":{"/":"z44vkPhdBNQZZd4HQGaDT34oRmkLQYputo7hqnYHKPdMNroFe7R"},"root":{"/":"z45oqTS4Ad9gDmfoQL9zY1g35AVWdbGbap1cQ37wbBaDZ5EASCs"},"time":1484475035,"tx":{"/":"z443fKyHcB6hA7QNg1XEFXsyKyat4EkUgVinTnvnhZLtJJgu2sH"},"uncles":{"/":"z43c7o74hjCAqnyneWetkyXU2i5KuGQLbYfVWZMvJMG4VTYABtz"}}
```

### Add an ethereum block encoded in RLP

This plugin also supports whether your block is an RLP encoded block header or
a block body (that is: its header, transactions and uncle list).

Let's test it out

#### Adding an RLP encoded block header

Just

```
cat ./test-data/eth-block-header-rlp-999999 | ipfs dag put --input-enc raw --format eth-block
```

You will get your cid `z43AaGF4uHSY4waU68L3DLUKHZP7yfZoo6QbLmid5HomZ4WtbWw`. Checking it,

```
ipfs dag get z43AaGF4uHSY4waU68L3DLUKHZP7yfZoo6QbLmid5HomZ4WtbWw
```

And we get our header back

```
{"bloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","coinbase":"0x52bc44d5378309ee2abf1539bf71de1b7d7be3b5","difficulty":12555463106190,"extra":"14MBAwOER2V0aIdnbzEuNC4yhWxpbnV4","gaslimit":3141592,"gasused":231000,"mixdigest":"0x5b10f4a08a6c209d426f6158bd24b574f4f7b7aa0099c67c14a1f693b4dd04d0","nonce":"0xf491f46b60fe04b3","number":999999,"parent":{"/":"z43AaGF6wP6uoLFEauru5oLK5JS5MGfNuGDK1xWEpQK4BqkJkL3"},"receipts":{"/":"z44vkPhhDSTXPAswvC1rdDunzkgZ7FgAAnhGQtNDNDk9m9N2BZA"},"root":{"/":"z45oqTSAZvPiiPV8hMZDH5fi4NkaAkMYTJC6PmaeWBmYUpbMpoh"},"time":1455404037,"tx":{"/":"z443fKyHHMwVy13VXtD4fdRcUXSqkr79Q5E8hcmEravVBq3Dc51"},"uncles":{"/":"z43c7o74hjCAqnyneWetkyXU2i5KuGQLbYfVWZMvJMG4VTYABtz"}}
```

#### Adding an RLP encoded block body (header, txs and uncle list)

We should get similars result trying to parse an RLP encoded block body.
Add `./test-data/eth-block-body-rlp-997522`

```
cat eth-block-body-rlp-997522 | ipfs dag put --input-enc raw --format eth-block
```

```
z43AaGExMLxj6ujVVbx3j4LRc6QGMBiqYCrgot5hG8Vnxm7Tf9M
```

```
ipfs dag get z43AaGExMLxj6ujVVbx3j4LRc6QGMBiqYCrgot5hG8Vnxm7Tf9M
```

```
{"bloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","coinbase":"0x4bb96091ee9d802ed039c4d1a5f6216f90f81b01","difficulty":11966502474733,"extra":"14MBBACER2V0aIdnbzEuNS4xhWxpbnV4","gaslimit":3141592,"gasused":21000,"mixdigest":"0x2565992ba4dbd7ab3bb08d1da34051ae1d90c79bc637a21aa2f51f6380bf5f6a","nonce":"0xf7a14147c2320b2d","number":997522,"parent":{"/":"z43AaGF24mjRxbn7A13gec2PjF5XZ1WXXCyhKCyxzYVBcxp3JuG"},"receipts":{"/":"z44vkPheUUg5HBpxkq5sFFz5d9ckigtBBW7WCJXQSZA1gV233Ap"},"root":{"/":"z45oqTRunK259j6Te1e3FsB27RJfDJop4XgbAbY39rwLmfoVWX4"},"time":1455362245,"tx":{"/":"z443fKyLvyDQBBQRGMNnPb8oPhPerbdwUX2QsQCUKqte1hy4kwD"},"uncles":{"/":"z43c7o73GVAMgEbpaNnaruD3ZbF4T2bqHZgFfyWqCejibzvJk41"}}

```

### TODO

* Support input for RLP encoded state trie elements (`eth-state-trie`).
* Support processing of txs to make the user able to retrieve `eth-tx` and `eth-tx-trie`.
* Supoort input for JSON encoded Transaction Receipts (`eth-tx-receipt`).
* The rest of the IPLD types in diverse inputs!
