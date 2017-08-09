# Ethereum IPLD plugin

Add ethereum support to ipfs!

## Building
Make sure to rewrite the gx dependencies in the directory above first, then
Either run `make` or `go build -buildmode=plugin -o=ethereum.so`.

* *NOTE*: As of [2017.08.09](https://golang.org/pkg/plugin/) the `plugins` lib
in Go only works in Linux.

## Installing
Move `ethereum.so` to `$IPFS_PATH/plugins/ethereum.so` and set it to be executable:

```sh
mkdir -p ~/.ipfs/plugins
mv ethereum.so ~/.ipfs/plugins/
chmod +x ~/.ipfs/plugins/ethereum.so
```

### I don't have linux but I want to do this somehow!

As stated above, the _plugin_ library only works in Linux. Bug the go team to
support your system!

* Or use a linux virtualbox, and mount this directory.

* Or hack your way via docker-fu [with this short, unsupported guide](hacks/docker.md)

* Or, if you are in OSX, [use this handy script](hacks/osx.sh)

## Usage and Examples

Make sure you have the right version of ipfs installed and start up the ipfs daemon!

### Add an ethereum block written in JSON

You may want to take a block given by your favorite client's JSON RPC API.
We have a couple of those in the `test-data` directory.

```
cat ./test_data/eth-block-body-json-997522 | ipfs dag put --input-enc json --format eth-block
```

And get the CID of the block header back

```
z43AaGEzuAXhWf9pWAm63QCERtFpqcc6gQX3QBBNaG1syxGGhg6
```

Now, you can get this block header

```
ipfs dag get z43AaGEzuAXhWf9pWAm63QCERtFpqcc6gQX3QBBNaG1syxGGhg6
```

Which will get you (with the right IPLD cids formatted for the other objects)

```
{"bloom":"0x00000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000","coinbase":"0x4bb96091ee9d802ed039
c4d1a5f6216f90f81b01","difficulty":11966502474733,"extra":"0xd783010400844765746
887676f312e352e31856c696e7578","gaslimit":3141592,"gasused":21000,"mixdigest":"0
x2565992ba4dbd7ab3bb08d1da34051ae1d90c79bc637a21aa2f51f6380bf5f6a","nonce":"0xf7
a14147c2320b2d","number":997522,"parent":{"/":"z43AaGF24mjRxbn7A13gec2PjF5XZ1WXX
CyhKCyxzYVBcxp3JuG"},"parentHash":"0x8ad6d5cbe7ec75ed71d5153dd58f2fd413b17c398ad
2a7d9309459ce884e6c9b","receiptHash":"0xa73a95d90de29c66220c8b8da825cf34ae969efc
7f9a878d8ed893565e4b4676","receipts":{"/":"z44vkPhjt2DpRokuesTzi6BKDriQKFEwe4Pvm
6HLAK3YWiHDzrR"},"root":{"/":"z45oqTRunK259j6Te1e3FsB27RJfDJop4XgbAbY39rwLmfoVWX
4"},"rootHash":"0x11e5ea49ecbee25a9b8f267492a5d296ac09cf6179b43bc334242d052bac59
63","time":1455362245,"tx":{"/":"z443fKyLvyDQBBQRGMNnPb8oPhPerbdwUX2QsQCUKqte1hy
4kwD"},"txHash":"0x7ab22cfcf6db5d1628ac888c25e6bc49aba2faaa200fc880f800f1db1e8bd
3cc","uncleHash":"0x08793b633d0b21b980107f3e3277c6693f2f3739e0c676a238cbe24d9ae6
e252","uncles":{"/":"z43c7o73GVAMgEbpaNnaruD3ZbF4T2bqHZgFfyWqCejibzvJk41"}}
```

You can read it better with some help. For example, use `python -m json.tool` to
get

```
{
    "bloom": "0x00000000.. edited ..0",
    "coinbase": "0x4bb96091ee9d802ed039c4d1a5f6216f90f81b01",
    "difficulty": 11966502474733,
    "extra": "0xd783010400844765746887676f312e352e31856c696e7578",
    "gaslimit": 3141592,
    "gasused": 21000,
    "mixdigest": "0x2565992ba4dbd7ab3bb08d1da34051ae1d90c79bc637a21aa2f51f6380bf5f6a",
    "nonce": "0xf7a14147c2320b2d",
    "number": 997522,
    "parent": {
        "/": "z43AaGF24mjRxbn7A13gec2PjF5XZ1WXXCyhKCyxzYVBcxp3JuG"
    },
    "parentHash": "0x8ad6d5cbe7ec75ed71d5153dd58f2fd413b17c398ad2a7d9309459ce884e6c9b",
    "receiptHash": "0xa73a95d90de29c66220c8b8da825cf34ae969efc7f9a878d8ed893565e4b4676",
    "receipts": {
        "/": "z44vkPhjt2DpRokuesTzi6BKDriQKFEwe4Pvm6HLAK3YWiHDzrR"
    },
    "root": {
        "/": "z45oqTRunK259j6Te1e3FsB27RJfDJop4XgbAbY39rwLmfoVWX4"
    },
    "rootHash": "0x11e5ea49ecbee25a9b8f267492a5d296ac09cf6179b43bc334242d052bac5963",
    "time": 1455362245,
    "tx": {
        "/": "z443fKyLvyDQBBQRGMNnPb8oPhPerbdwUX2QsQCUKqte1hy4kwD"
    },
    "txHash": "0x7ab22cfcf6db5d1628ac888c25e6bc49aba2faaa200fc880f800f1db1e8bd3cc",
    "uncleHash": "0x08793b633d0b21b980107f3e3277c6693f2f3739e0c676a238cbe24d9ae6e252",
    "uncles": {
        "/": "z43c7o73GVAMgEbpaNnaruD3ZbF4T2bqHZgFfyWqCejibzvJk41"
    }
}
```

Note that we are including both the hashes and `cid`s of the links.

#### Piping from the RPC

The astute reader will say "_Let's then pipe the output of my RPC directly
to IPFS!_"

```
curl -s -X POST \
	--data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x1b4", true],"id":1}' \
	https://mainnet.infura.io | ipfs dag put --input-enc json --format eth-block && echo
```

Will give you

```
z43AaGF7XiKhgVVcYxNJv3ZrebEkDE5yhna22N74AusBdMvi6pV
```


And then calling

```
ipfs dag get z43AaGF7XiKhgVVcYxNJv3ZrebEkDE5yhna22N74AusBdMvi6pV
```

Returns 

```
{"bloom":"0x00000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000","coinbase":"0xbb7b8287f3f0a933474a
79eae42cbca977791171","difficulty":21109876668,"extra":"0x476574682f4c5649562f76
312e302e302f6c696e75782f676f312e342e32","gaslimit":5000,"gasused":0,"mixdigest":
"0x4fffe9ae21f1c9e15207b1f472d5bbdd68c9595d461666602f2be20daf5e7843","nonce":"0x
689056015818adbe","number":436,"parent":{"/":"z43AaGF8SkCtKoht2v1e3yC9DWHi4iV2dy
nyi3BTCP7sPs7HR2T"},"parentHash":"0xe99e022112df268087ea7eafaf4790497fd21dbeeb6b
d7a1721df161a6657a54","receiptHash":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b9
96cadc001622fb5e363b421","receipts":{"/":"z44vkPheUUg5HBpxkq5sFFz5d9ckigtBBW7WCJ
XQSZA1gV233Ap"},"root":{"/":"z45oqTS9WCLjMeLnFvTbWiqxXRi1PdwYtDjnNQy6PyWKokGD8r8
"},"rootHash":"0xddc8b0234c2e0cad087c8b389aa7ef01f7d79b2570bccb77ce48648aa61c904
d","time":1438271100,"tx":{"/":"z443fKyJXGFJKPgzhha8eqpkEz3rHUL5M7cvcfJQVGzwt3Mw
cVn"},"txHash":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b4
21","uncleHash":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49
347","uncles":{"/":"z43c7o74hjCAqnyneWetkyXU2i5KuGQLbYfVWZMvJMG4VTYABtz"}}
```

Or go even more _extreme_ with a single pipe

```
ipfs dag get \
	$(curl -s -X POST \
		--data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x2DC6C0", true],"id":1}' \
	https://mainnet.infura.io | ipfs dag put --input-enc json --format eth-block)
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
{"bloom":"0x00000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000","coinbase":"0x52bc44d5378309ee2abf
1539bf71de1b7d7be3b5","difficulty":12555463106190,"extra":"0xd783010303844765746
887676f312e342e32856c696e7578","gaslimit":3141592,"gasused":231000,"mixdigest":"
0x5b10f4a08a6c209d426f6158bd24b574f4f7b7aa0099c67c14a1f693b4dd04d0","nonce":"0xf
491f46b60fe04b3","number":999999,"parent":{"/":"z43AaGF6wP6uoLFEauru5oLK5JS5MGfN
uGDK1xWEpQK4BqkJkL3"},"parentHash":"0xd33c9dde9fff0ebaa6e71e8b26d2bda15ccf111c7a
f1b633698ac847667f0fb4","receiptHash":"0x7fa0f6ca2a01823208d80801edad37e3e3a003b
55c89319b45eb1f97862ad229","receipts":{"/":"z44vkPhhDSTXPAswvC1rdDunzkgZ7FgAAnhG
QtNDNDk9m9N2BZA"},"root":{"/":"z45oqTSAZvPiiPV8hMZDH5fi4NkaAkMYTJC6PmaeWBmYUpbMp
oh"},"rootHash":"0xed98aa4b5b19c82fb35364f08508ae0a6dec665fa57663dca94c5d70554cd
e10","time":1455404037,"tx":{"/":"z443fKyHHMwVy13VXtD4fdRcUXSqkr79Q5E8hcmEravVBq
3Dc51"},"txHash":"0x447cbd8c48f498a6912b10831cdff59c7fbfcbbe735ca92883d4fa06dcd7
ae54","uncleHash":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d
49347","uncles":{"/":"z43c7o74hjCAqnyneWetkyXU2i5KuGQLbYfVWZMvJMG4VTYABtz"}}
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
{"bloom":"0x00000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000","coinbase":"0x4bb96091ee9d802ed039
c4d1a5f6216f90f81b01","difficulty":11966502474733,"extra":"0xd783010400844765746
887676f312e352e31856c696e7578","gaslimit":3141592,"gasused":21000,"mixdigest":"0
x2565992ba4dbd7ab3bb08d1da34051ae1d90c79bc637a21aa2f51f6380bf5f6a","nonce":"0xf7
a14147c2320b2d","number":997522,"parent":{"/":"z43AaGF24mjRxbn7A13gec2PjF5XZ1WXX
CyhKCyxzYVBcxp3JuG"},"parentHash":"0x8ad6d5cbe7ec75ed71d5153dd58f2fd413b17c398ad
2a7d9309459ce884e6c9b","receiptHash":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b
996cadc001622fb5e363b421","receipts":{"/":"z44vkPheUUg5HBpxkq5sFFz5d9ckigtBBW7WC
JXQSZA1gV233Ap"},"root":{"/":"z45oqTRunK259j6Te1e3FsB27RJfDJop4XgbAbY39rwLmfoVWX
4"},"rootHash":"0x11e5ea49ecbee25a9b8f267492a5d296ac09cf6179b43bc334242d052bac59
63","time":1455362245,"tx":{"/":"z443fKyLvyDQBBQRGMNnPb8oPhPerbdwUX2QsQCUKqte1hy
4kwD"},"txHash":"0x7ab22cfcf6db5d1628ac888c25e6bc49aba2faaa200fc880f800f1db1e8bd
3cc","uncleHash":"0x08793b633d0b21b980107f3e3277c6693f2f3739e0c676a238cbe24d9ae6
e252","uncles":{"/":"z43c7o73GVAMgEbpaNnaruD3ZbF4T2bqHZgFfyWqCejibzvJk41"}}
```

What's the difference you ask?

No difference in the output. You will always get a block header. **But**,
when you add a block body, (that is header, txs and ommers list), the
transactions get processed and imported into the IPLD merkle forest too.

## Navigate to a block's parent (and parent of a parent...)

If you have a chain of blocks available, you can easily navigate to a
block's parent and so on.

Import the following blocks

```
cat ./test_data/eth-block-body-json-999999 | ipfs dag put --input-enc json --format eth-block
cat ./test_data/eth-block-body-json-999998 | ipfs dag put --input-enc json --format eth-block
cat ./test_data/eth-block-header-rlp-999997 | ipfs dag put --input-enc raw --format eth-block
cat ./test_data/eth-block-header-rlp-999996pfs dag get z43AaGF4uHSY4waU68L3DLUKHZP7yfZoo6QbLmid5HomZ4WtbWw | ipfs dag put --input-enc raw --format eth-block
```

(Notice how we are using block headers and bodies in different encodings).

Now, let's see how this goes, so we have this block

```
ipfs dag get z43AaGF4uHSY4waU68L3DLUKHZP7yfZoo6QbLmid5HomZ4WtbWw
```

```
{"bloom":"0x00000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000","coinbase":"0x52bc44d5378309ee2abf
1539bf71de1b7d7be3b5","difficulty":12555463106190,"extra":"0xd783010303844765746
887676f312e342e32856c696e7578","gaslimit":3141592,"gasused":231000,"mixdigest":"
0x5b10f4a08a6c209d426f6158bd24b574f4f7b7aa0099c67c14a1f693b4dd04d0","nonce":"0xf
491f46b60fe04b3","number":999999,"parent":{"/":"z43AaGF6wP6uoLFEauru5oLK5JS5MGfN
uGDK1xWEpQK4BqkJkL3"},"parentHash":"0xd33c9dde9fff0ebaa6e71e8b26d2bda15ccf111c7a
f1b633698ac847667f0fb4","receiptHash":"0x7fa0f6ca2a01823208d80801edad37e3e3a003b
55c89319b45eb1f97862ad229","receipts":{"/":"z44vkPhhDSTXPAswvC1rdDunzkgZ7FgAAnhG
QtNDNDk9m9N2BZA"},"root":{"/":"z45oqTSAZvPiiPV8hMZDH5fi4NkaAkMYTJC6PmaeWBmYUpbMp
oh"},"rootHash":"0xed98aa4b5b19c82fb35364f08508ae0a6dec665fa57663dca94c5d70554cd
e10","time":1455404037,"tx":{"/":"z443fKyHHMwVy13VXtD4fdRcUXSqkr79Q5E8hcmEravVBq
3Dc51"},"txHash":"0x447cbd8c48f498a6912b10831cdff59c7fbfcbbe735ca92883d4fa06dcd7
ae54","uncleHash":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d
49347","uncles":{"/":"z43c7o74hjCAqnyneWetkyXU2i5KuGQLbYfVWZMvJMG4VTYABtz"}}
```

... and we call its parent

```
ipfs dag get z43AaGF4uHSY4waU68L3DLUKHZP7yfZoo6QbLmid5HomZ4WtbWw/parent
```

```
{"bloom":"0x00000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000","coinbase":"0xf8b483dba2c3b7176a3d
a549ad41a48bb3121069","difficulty":12561596698199,"extra":"0xd983010302844765746
887676f312e342e328777696e646f7773","gaslimit":3141592,"gasused":252000,"mixdiges
t":"0xcaf27314d80cb3e888d32646402d617d8f8379ca23a6b0255e974e407ffdd846","nonce":
"0xbc7609306a77d0a2","number":999998,"parent":{"/":"z43AaGF67aUUDzGGimXySbgNzJJi
tkTVUTvpaf9jrqxe8BKuJL2"},"parentHash":"0xc6fd988b2d086a7b6eee3d25bad45383039101
4ba268cf6cc5d139741cb51273","receiptHash":"0xb0310e47b0cc7d3bb24c65ec21ec0ddf8dc
f1672bc9866d6ba67e83d33215568","receipts":{"/":"z44vkPhkV1Tp7osq3p4yThA7EdE5ikvZ
UZTtDpvvpkMNGvxC9HZ"},"root":{"/":"z45oqTSAdVfS8g8n7NrSBKTeydujwoRgw52ZQehEZaVhC
d4QNx6"},"rootHash":"0xee8306f6cebba17153516cb6586de61d6294b49bc5534eb9378acb848
907b277","time":1455404013,"tx":{"/":"z443fKyKQh6b7HWVtYXuJi6shDUtUTsFaw4g3vToP5
n9eEvb3Jn"},"txHash":"0x6414d72a4c223bce7d1309869332b148670eb66af4e3b3ba6d1a55aa
0bb3fd4f","uncleHash":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142f
d40d49347","uncles":{"/":"z43c7o74hjCAqnyneWetkyXU2i5KuGQLbYfVWZMvJMG4VTYABtz"}}
```

Why not calling its "grandparent" (parent of a parent)?

```
ipfs dag get z43AaGF4uHSY4waU68L3DLUKHZP7yfZoo6QbLmid5HomZ4WtbWw/parent/parent
```

```
{"bloom":"0x00000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000","coinbase":"0x52bc44d5378309ee2abf
1539bf71de1b7d7be3b5","difficulty":12567733286589,"extra":"0xd783010303844765746
887676f312e342e32856c696e7578","gaslimit":3141592,"gasused":189000,"mixdigest":"
0xedd380b8b600469c89d763fbb73c1ed4128164c2b8ccc41ed73d3e16f8d2a8de","nonce":"0x7
b9a013e3da652ca","number":999997,"parent":{"/":"z43AaGF26webNZ5MTwHkhjcQZEqGkBfx
65gDXALsa8171tUf5tU"},"parentHash":"0x8b6535a0e3e346ee87e0194456d95971988d398163
9bf7065f602d11b7adeab9","receiptHash":"0x85c15ea267eda062e4470a875f6fe3135d8d63f
561e409f5c0732c5539c35d1b","receipts":{"/":"z44vkPhhdMZfdtKqVtLTzXcppkFsT3W4PgZt
K5SmTgPbHZsTUKC"},"root":{"/":"z45oqTS1N7W29LJtqzpcFZie3cPgwAyf6BT73MetEePYD9RGF
n3"},"rootHash":"0x64d912e03889ea4754dd1039bd38a19677335aacd3399a3c3a3a74314588d
584","time":1455403990,"tx":{"/":"z443fKyFeSt9z2MYAdG8GU26oT882qxkwCLE7C61UNSZJ4
RpYEQ"},"txHash":"0x2c2c26e1629b431ad5fa033d90f4ec5c2b59d437cf1a34082195f5f771b3
735d","uncleHash":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d
49347","uncles":{"/":"z43c7o74hjCAqnyneWetkyXU2i5KuGQLbYfVWZMvJMG4VTYABtz"}}
```

Since we are there, let's see what happens with their parent in turn

```
ipfs dag get z43AaGF4uHSY4waU68L3DLUKHZP7yfZoo6QbLmid5HomZ4WtbWw/parent/parent/parent
```

```
{"bloom":"0x00000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000000000000000000000000
00000000000000000000000000000000000000000000","coinbase":"0xf8b483dba2c3b7176a3d
a549ad41a48bb3121069","difficulty":12573872872824,"extra":"0xd983010302844765746
887676f312e342e328777696e646f7773","gaslimit":3141592,"gasused":42000,"mixdigest
":"0xb388de31f2a59ec8cacf363866101a6904545d4dffa5d69537427e0df6f3aa2f","nonce":"
0x5645faf4502c64d9","number":999996,"parent":{"/":"z43AaGF5oDbc3A3yMSMAjGaWeVCGv
A2Pgrbb1C4mqfpZSCBhQiC"},"parentHash":"0xc249891eb893a583be09d904b7d952988098fd8
bdf5de09003f7a4811fd0c591","receiptHash":"0xc7ce189fbc688fd45b844288a4d6016ca600
2d77b1fa9e741716622608fb9312","receipts":{"/":"z44vkPhn5Bj9VT3BVsuDaMZ87gBQfYGtz
SuJnYiEzLc91jggZbP"},"root":{"/":"z45oqTS3Sk3JQMTG7W3nXFF5o8RVGvbNotPLpGP1tbNJ22
RNRBZ"},"rootHash":"0x83c016016b084a6074ea327e9ede376501965a8a18141f1bb3aef7a7c7
32bfec","time":1455403968,"tx":{"/":"z443fKyPdU3PSLUXTNtUhQ8BefUipk3CQf3n4xvRMf8
LLRbiRKN"},"txHash":"0xa2c9608d6d1083b677012732bf149d232f02d32d465423b1ccb630693
8bad451","uncleHash":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd
40d49347","uncles":{"/":"z43c7o74hjCAqnyneWetkyXU2i5KuGQLbYfVWZMvJMG4VTYABtz"}}
```

... And so on...

## Navigate through the transactions of a block

(WIP)

## TODO

* `[0x90]` - `eth-block` input:
  * Can we get the `eth-tx` (`[0x91]`) pointed by the `eth-tx-trie` leaf?

* `[0x92]` - `eth-tx-receipt`:
  * Propose a script to get all receipts from a block and make a JSON array of them.
  * Support the input of this JSON array to form the `eth-tx-receipt-trie` (`[0x96]`) leaves, and the `eth-tx-receipt` objects.

* `[0x97]` - `eth-state-trie`. Support input for RLP encoded state trie elements.
  * HINT: We get them from the Parity IPFS API.

* The rest of the IPLD ETH Types:
  * `[0x93]` - `eth-account-snapshot`
  * `[0x94]` - `eth-block-list`
  * `[0x98]` - `eth-storage-trie`
