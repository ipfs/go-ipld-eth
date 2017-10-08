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

Which will give you (with the right IPLD cids formatted for the other objects).

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
CyhKCyxzYVBcxp3JuG"},"receipts":{"/":"z44vkPhjt2DpRokuesTzi6BKDriQKFEwe4Pvm6HLAK
3YWiHDzrR"},"root":{"/":"z45oqTRunK259j6Te1e3FsB27RJfDJop4XgbAbY39rwLmfoVWX4"},"
time":1455362245,"tx":{"/":"z443fKyLvyDQBBQRGMNnPb8oPhPerbdwUX2QsQCUKqte1hy4kwD"
},"uncles":{"/":"z43c7o73GVAMgEbpaNnaruD3ZbF4T2bqHZgFfyWqCejibzvJk41"}}
```

You can read it better with some help. For example, use `python -m json.tool` to
get

```
{
    "bloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
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
    "receipts": {
        "/": "z44vkPhjt2DpRokuesTzi6BKDriQKFEwe4Pvm6HLAK3YWiHDzrR"
    },
    "root": {
        "/": "z45oqTRunK259j6Te1e3FsB27RJfDJop4XgbAbY39rwLmfoVWX4"
    },
    "time": 1455362245,
    "tx": {
        "/": "z443fKyLvyDQBBQRGMNnPb8oPhPerbdwUX2QsQCUKqte1hy4kwD"
    },
    "uncles": {
        "/": "z43c7o73GVAMgEbpaNnaruD3ZbF4T2bqHZgFfyWqCejibzvJk41"
    }
}
```

**NOTE:** From now on, in this tutorial, we will be applying this tool
to improve readability of the output.

#### Piping from the RPC

The astute reader will say "_Let's then pipe the output of my RPC directly
to IPFS!_". OK then,

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
{
    "bloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "coinbase": "0xbb7b8287f3f0a933474a79eae42cbca977791171",
    "difficulty": 21109876668,
    "extra": "0x476574682f4c5649562f76312e302e302f6c696e75782f676f312e342e32",
    "gaslimit": 5000,
    "gasused": 0,
    "mixdigest": "0x4fffe9ae21f1c9e15207b1f472d5bbdd68c9595d461666602f2be20daf5e7843",
    "nonce": "0x689056015818adbe",
    "number": 436,
    "parent": {
        "/": "z43AaGF8SkCtKoht2v1e3yC9DWHi4iV2dynyi3BTCP7sPs7HR2T"
    },
    "receipts": {
        "/": "z44vkPheUUg5HBpxkq5sFFz5d9ckigtBBW7WCJXQSZA1gV233Ap"
    },
    "root": {
        "/": "z45oqTS9WCLjMeLnFvTbWiqxXRi1PdwYtDjnNQy6PyWKokGD8r8"
    },
    "time": 1438271100,
    "tx": {
        "/": "z443fKyJXGFJKPgzhha8eqpkEz3rHUL5M7cvcfJQVGzwt3MwcVn"
    },
    "uncles": {
        "/": "z43c7o74hjCAqnyneWetkyXU2i5KuGQLbYfVWZMvJMG4VTYABtz"
    }
}
```

Moreover, you can go even more _extreme_ with a single pipe...

```
ipfs dag get \
	$(curl -s -X POST \
		--data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x2DC6C0", true],"id":1}' \
	https://mainnet.infura.io | ipfs dag put --input-enc json --format eth-block)
```

Which retrieves from the remote RPC in INFURA, imports into IPFS, and then retrieves the very result.

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

We get our header back, in JSON.

```
{
    "bloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "coinbase": "0x52bc44d5378309ee2abf1539bf71de1b7d7be3b5",
    "difficulty": 12555463106190,
    "extra": "0xd783010303844765746887676f312e342e32856c696e7578",
    "gaslimit": 3141592,
    "gasused": 231000,
    "mixdigest": "0x5b10f4a08a6c209d426f6158bd24b574f4f7b7aa0099c67c14a1f693b4dd04d0",
    "nonce": "0xf491f46b60fe04b3",
    "number": 999999,
    "parent": {
        "/": "z43AaGF6wP6uoLFEauru5oLK5JS5MGfNuGDK1xWEpQK4BqkJkL3"
    },
    "receipts": {
        "/": "z44vkPhhDSTXPAswvC1rdDunzkgZ7FgAAnhGQtNDNDk9m9N2BZA"
    },
    "root": {
        "/": "z45oqTSAZvPiiPV8hMZDH5fi4NkaAkMYTJC6PmaeWBmYUpbMpoh"
    },
    "time": 1455404037,
    "tx": {
        "/": "z443fKyHHMwVy13VXtD4fdRcUXSqkr79Q5E8hcmEravVBq3Dc51"
    },
    "uncles": {
        "/": "z43c7o74hjCAqnyneWetkyXU2i5KuGQLbYfVWZMvJMG4VTYABtz"
    }
}
```

#### Adding an RLP encoded block body (header, txs and uncle list)

We should get similars result trying to parse an RLP encoded block body.
Add `./test-data/eth-block-body-rlp-997522`

```
cat ./test_data/eth-block-body-rlp-997522 | ipfs dag put --input-enc raw --format eth-block
```

```
z43AaGExMLxj6ujVVbx3j4LRc6QGMBiqYCrgot5hG8Vnxm7Tf9M
```

```
ipfs dag get z43AaGExMLxj6ujVVbx3j4LRc6QGMBiqYCrgot5hG8Vnxm7Tf9M
```

```
{
    "bloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
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
    "receipts": {
        "/": "z44vkPheUUg5HBpxkq5sFFz5d9ckigtBBW7WCJXQSZA1gV233Ap"
    },
    "root": {
        "/": "z45oqTRunK259j6Te1e3FsB27RJfDJop4XgbAbY39rwLmfoVWX4"
    },
    "time": 1455362245,
    "tx": {
        "/": "z443fKyLvyDQBBQRGMNnPb8oPhPerbdwUX2QsQCUKqte1hy4kwD"
    },
    "uncles": {
        "/": "z43c7o73GVAMgEbpaNnaruD3ZbF4T2bqHZgFfyWqCejibzvJk41"
    }
}
```

##### What's the difference you ask?

No difference in the output. You will always get a block header. **But**,
when you add a block body, (that is header, txs and ommers list),**the
transactions get processed and imported into the IPLD merkle forest too**.

## Navigate to a block's parent (and parent of a parent...)

If you have a chain of blocks available, you can easily navigate to a
block's parent and so on.

Import the following blocks

```
cat ./test_data/eth-block-body-json-999999 | ipfs dag put --input-enc json --format eth-block
cat ./test_data/eth-block-body-json-999998 | ipfs dag put --input-enc json --format eth-block
cat ./test_data/eth-block-header-rlp-999997 | ipfs dag put --input-enc raw --format eth-block
cat ./test_data/eth-block-header-rlp-999997 | ipfs dag put --input-enc raw --format eth-block
```

(Notice how we are using block headers and bodies in different encodings).

Now, let's see how this goes, so we have this block (first cid returned)

```
ipfs dag get z43AaGF4uHSY4waU68L3DLUKHZP7yfZoo6QbLmid5HomZ4WtbWw | python -m json.tool | grep number
```

```
"number": 999999,
```

and we call its parent

```
ipfs dag get z43AaGF4uHSY4waU68L3DLUKHZP7yfZoo6QbLmid5HomZ4WtbWw/parent | python -m json.tool | grep number
```

Unsurprisingly we get


```
"number": 999998,
```

Why not calling its "grandparent" (parent of a parent)?

```
ipfs dag get z43AaGF4uHSY4waU68L3DLUKHZP7yfZoo6QbLmid5HomZ4WtbWw/parent/parent | python -m json.tool | grep number
```

to get...

```
"number": 999997,
```

Since we are there, let's see what happens with their parent in turn

```
ipfs dag get z43AaGF4uHSY4waU68L3DLUKHZP7yfZoo6QbLmid5HomZ4WtbWw/parent/parent/parent | python -m json.tool | grep number
```

```
"number": 999996,
```

... And so on...

## Navigate through the transactions of a block

Cool. So let's say that the IPLD merkle forest has the transactions belonging
to the block 4,139,497.

We can import them from its block body json

```
cat ./test_data/eth-block-body-json-4139497 | ipfs dag put --input-enc json --format eth-block
```

Getting back the cid `z43AaGEtGPmuXQpwmknmt7hcQRRuoX6SjgDaMTfkxYcXJMn4VPx`.

We can navigate the merkle tree of the transactions in this block
resolving the link `/tx` and referencing with their indices with their RLP
equivalent. For example to get to the transaction `0x01` (in RLP), we just

```
ipfs dag get z43AaGEtGPmuXQpwmknmt7hcQRRuoX6SjgDaMTfkxYcXJMn4VPx/tx/01
```

Returning

```
{
    "": {
        "gas": 186844,
        "gasPrice": 51000000000,
        "input": "a9059cbb000000000000000000000000744346c50253300694aea6d7e03f55a3ea91f8a30000000000000000000000000000000000000000000000000000013061e0a9ab",
        "nonce": 790605,
        "r": "0xe925321edf5dc905fa0ebf9a08d8915e0ce90463d55c19e8bdf0dc8e5e6ddc73",
        "s": "0x328a5099139ae2e3f3be2736dec30fd2b3240892b77575e588b8f84a0e11307b",
        "toAddress": "0x41e5560054824ea6b0732e656e3ad64e20e94e45",
        "v": "0x25",
        "value": "0x0"
    },
    "type": "leaf"
}
```

There, we have a leaf of the trie, to access individual fields, we just resolve
them

```
ipfs dag get z43AaGEtGPmuXQpwmknmt7hcQRRuoX6SjgDaMTfkxYcXJMn4VPx/tx/01/nonce
```

Obtaining `790605`.

Now, Let's do some manual traversing

```
ipfs dag get ipfs dag get z43AaGEtGPmuXQpwmknmt7hcQRRuoX6SjgDaMTfkxYcXJMn4VPx/tx
```

Returns a branch node

```
{
...


    "8": {
        "/": "z443fKyRJvB8PQEdWTL44qqoo2DeZr8QwkasSAfEcWJ6uDUWyh6"
    },
    "9": null,
    "a": null,
    "b": null,
    "c": null,
    "d": null,
    "e": null,
    "f": null,
    "type": "branch"
}
```

What happens if we follow to the `8` child?

```
ipfs dag get z43AaGEtGPmuXQpwmknmt7hcQRRuoX6SjgDaMTfkxYcXJMn4VPx/tx/8
```

Mmm, another branch

```
{
    "0": {
        "/": "z443fKyQhyzood9hQHXyYzbGAJeJMxMWDpbrUTXGm55WxoGGWhn"
    },
    "1": {
        "/": "z443fKyMsFsxojbxvSCpJApyCvWKE9jCgrGc98cKRJjMgVBptvN"
    },
    "2": {
        "/": "z443fKyR2PNJ3gNLTrPEmkHJh4YJ2mNMU9QX4HuBFNfBGnkb444"
    },
...
}
```

Try again, with `2`

```
ipfs dag get z43AaGEtGPmuXQpwmknmt7hcQRRuoX6SjgDaMTfkxYcXJMn4VPx/tx/82 | jsontool
```

```
{
    "01": {
        "/": "z443fKyJaFfaE7Hsozvv7HGEHqNWPEhkNgzgnXjVKdxqCE74PgF"
    },
    "type": "extension"
}
```

OK, an extension. It has a key of `01`, so it's telling us that the only way
to follow into this rabbit hole (i.e. be able to catch the next value), is by
entering the next two nibbles (`01`)

```
ipfs dag get z43AaGEtGPmuXQpwmknmt7hcQRRuoX6SjgDaMTfkxYcXJMn4VPx/tx/820
Error: not enough nibbles to traverse this extension

 ipfs dag get z43AaGEtGPmuXQpwmknmt7hcQRRuoX6SjgDaMTfkxYcXJMn4VPx/tx/8201
{"0":{"/":"z443fKyNhpksFN3ixSGZr2QMD1YrZoSQFz5zFt2ZwyTvaZPtWWw"},"1":{"/":"z443fKySNAgfM3gM5R2W6aEtEzgekfY4QAx2sfTeVFp3uJiAQzd"},"2":{"/":"z443fKySgQc9JHXeNyYCzgxN7358eW5wvM6yRm9MVbhd6gofbB7"},"3":{"/":"z443fKyFPKZUHbZF9Q3hPHrQvC3wX4A1BFrrXdwJmTQZaAx6rwN"},"4":null,"5":null,"6":null,"7":null,"8":null,"9":null,"a":null,"b":null,"c":null,"d":null,"e":null,"f":null,"type":"branch"}
```

We eventually reach a leaf at `820100`, which is the RLP equivalent of `255`.

```
ipfs dag get z43AaGEtGPmuXQpwmknmt7hcQRRuoX6SjgDaMTfkxYcXJMn4VPx/tx/820100
```

```
{
    "": {
        "gas": 90000,
        "gasPrice": 4000000000,
        "input": "",
        "nonce": 40243,
        "r": "0x981b6223c9d3c319716da3cf057da84acf0fef897f4003d8a362d7bda42247db",
        "s": "0x66be134c4bc432125209b5056ef274b7423bcac7cc398cf60b83aaff7b95469f",
        "toAddress": "0xe0e6c781b8cba08bc8407eac0101b668d1fa6f49",
        "v": "0x26",
        "value": "0xc495a958603400"
    },
    "type": "leaf"
}
```

And getting their values just referencing them

```
ipfs dag get z43AaGEtGPmuXQpwmknmt7hcQRRuoX6SjgDaMTfkxYcXJMn4VPx/tx/820100/gasPrice
4000000000
```
### Traversing the state trie

#### Inserting test data

##### Block 0

In case you haven't downloaded the json of the block 0 yourself yet from an RPC
ethereum client, you can find it in the `test_data` dirctory, let's add it to
your datastore

```
cat test_data/eth-block-body-json-0 | ipfs dag put --input-enc json --format eth-block
```

You should get its cid, write it down somewhere. As we are going to use it later.

```
z43AaGF73rnZ14vjAkMQ8xoNfBShmq8qaiqFuELAx1vxSTzfGY2
```

##### A traversal of state tries

In the `test_data` directory, you can find several RLP encoded raw files including
the state trie nodes needed for this experiment. We will import them using the command

```
find test_data/. -name "eth-state-trie-rlp-*" -exec sh -c "cat {} | ipfs dag put --input-enc raw --format eth-state-trie; echo" \;
```

Getting the result

```
z45oqTRuZDa8Kvo9MWtEmGZJfnsg39ngDa4CyHY77oRD5Yd8PiX
z45oqTRzQCbfigmMLqYdCoe2qGfo8pimaPTFGDwcsjT7WeYNQJG
z45oqTS26iKhYHbe5shTQhyW6DZE1Ffy24ENRVKQUvvZZwAyFuV
z45oqTS2HJkhG9adyTo1oAyZX9gArKH9mYunD1B3cY6Daowu3BC
z45oqTS87AwCRMhezXeH45bHKELbVeAh4WaDnC4eLa45XAKxpUE
z45oqTS8xp8GjTp5mZM4KS8trCfiK2hbZ2RoKwZhpnMr9pLQsTr
z45oqTS97WG4WsMjquajJ8PB9Ubt3ks7rGmo14P5XWjnPL7LHDM
z45oqTSAQWTPvpx4JsgqCzkpsn64J4rfqXvy25nSxvVdWbHd1ue
```

#### Let's do this

##### Getting the traversal path

Suppose we want to know at the block `0` (genesis block), what is the `balance`
of the account `0x5abfec25f74cd88437631a7731906932776356f9` (HINT: The EF account).

The first thing we need to do is getting the `keccak256()` hash of this address.

Sounds cumbersome? Maybe. In the future, we may add a system to transform the address
into their secure hash, and _then_ initiate the resolving... For now, we can use a
script in our favorite language to get this hash.

In go it'd be

```go
package main

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	var ethAddress string

	flag.String(&ethAddress, "eth-address", "5abfec25f74cd88437631a7731906932776356f9", "Address which keccak-256 hash we want")
	flag.Parse()

	kb, err := hex.DecodeString(ethAddress)
	if err != nil {
		panic(err)
	}
	secureKey := crypto.Keccak256(kb)

	fmt.Printf("%x\n", secureKey)
}
```

Copy it into the `/tmp` directory, name it `getKeccak256.go` and do

```
cd /tmp
go build getKeccak256.go

## Remove the "0x" part
./getKeccak256 --eth-address 5abfec25f74cd88437631a7731906932776356f9

## You get this
## cdd3e25edec0a536a05f5e5ab90a5603624c0ed77453b2e8f955cf8b43d4d0fb
```

OK. This last result _is_ the very traversal path!
Let's move though this path from the genesis block

##### Traversing step by step

As we got the cid of the block 0, we can get it performing

```
ipfs dag get z43AaGF73rnZ14vjAkMQ8xoNfBShmq8qaiqFuELAx1vxSTzfGY2
```

To access the state trie root, let's call `/root` in this element

```
ipfs dag get z43AaGF73rnZ14vjAkMQ8xoNfBShmq8qaiqFuELAx1vxSTzfGY2/root
```

Getting a `branch` kind of node

```
{
...
    "c": {
        "/": "z45oqTS26iKhYHbe5shTQhyW6DZE1Ffy24ENRVKQUvvZZwAyFuV"
    },
    "d": {
        "/": "z45oqTS2s1xH8sZBRmqBw5nKC3HCVAn3kNEoBFiESe6Fh1zybLG"
    },
    "e": {
        "/": "z45oqTRz6Skjh42RUJrMTv2Jdcmf6Z8bnRcrXw6RFRdJpUgXKkT"
    },
    "f": {
        "/": "z45oqTS3rVwGVvYm6gnJH1MtuY3bFyReLMq6XoNgVqs5naK9wEr"
    },
    "type": "branch"
}
```

Now, as the path is `cdd3e25edec0a536a05f5e5ab90a5603624c0ed77453b2e8f955cf8b43d4d0fb`,
we want the `c` link in this branch, to start with.

Let's do a couple iterations

```
ipfs dag get z43AaGF73rnZ14vjAkMQ8xoNfBShmq8qaiqFuELAx1vxSTzfGY2/root/c
```

```
{
...
    },
    "d": {
        "/": "z45oqTS8xp8GjTp5mZM4KS8trCfiK2hbZ2RoKwZhpnMr9pLQsTr"
    },
    "e": {
        "/": "z45oqTSBFKHowXCEU6GV7o3ZHhRDCRgfbU1JHHLM1abFkwx2ZSm"
    },
    "f": {
        "/": "z45oqTRyYKpFd3zSAkh8udFdvJ2VsCyZ6p3daG4cqLUnJUkJwCt"
    },
    "type": "branch"
}
```

```
ipfs dag get z43AaGF73rnZ14vjAkMQ8xoNfBShmq8qaiqFuELAx1vxSTzfGY2/root/c/d
```

```
{
...
    },
    "d": {
        "/": "z45oqTS2HJkhG9adyTo1oAyZX9gArKH9mYunD1B3cY6Daowu3BC"
    },
    "e": {
        "/": "z45oqTRynfnn9xUv6nyyEkx8JDLhSBxFbhcHLdqJtq261yMWHYN"
    },
    "f": {
        "/": "z45oqTRtpRLZUE4yoq4AjB9JfYGvvhPNJejuNGieB7kc6ECh6Ze"
    },
    "type": "branch"
}
```

```
ipfs dag get z43AaGF73rnZ14vjAkMQ8xoNfBShmq8qaiqFuELAx1vxSTzfGY2/root/c/d/d
```

```
{
    "3e25edec0a536a05f5e5ab90a5603624c0ed77453b2e8f955cf8b43d4d0fb": {
        "balance": 11901484239480000000000000,
        "codeHash": "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
        "nonce": 0,
        "root": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"
    },
    "type": "leaf"
}
```

Boom! We hit a `leaf` pretty quickly! This one contains the rest of the path.

To access the contains of the account, we need to complete the leaf, in other words, do

```
ipfs dag get z45oqTS97WG4WsMjquajJ8PB9Ubt3ks7rGmo14P5XWjnPL7LHDM/c/d/d/3/e/2/5/e/d/e/c/0/a/5/3/6/a/0/5/f/5/e/5/a/b/9/0/a/5/6/0/3/6/2/4/c/0/e/d/7/7/4/5/3/b/2/e/8/f/9/5/5/c/f/8/b/4/3/d/4/d/0/f/b/
```

(We expect you to use a tool "in real life" to do this... I did it by hand and its boooring).

Getting the object

```
{
    "balance": 11901484239480000000000000,
    "codeHash": "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
    "nonce": 0,
    "root": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"
}
```

From there is a matter of taking the right element, in this case, `balance`

```
ipfs dag get z45oqTS97WG4WsMjquajJ8PB9Ubt3ks7rGmo14P5XWjnPL7LHDM/c/d/d/3/e/2/5/e/d/e/c/0/a/5/3/6/a/0/5/f/5/e/5/a/b/9/0/a/5/6/0/3/6/2/4/c/0/e/d/7/7/4/5/3/b/2/e/8/f/9/5/5/c/f/8/b/4/3/d/4/d/0/f/b/balance
```

Getting this result. A LOT OF CRYPTO-MONEY.

```
11901484239480000000000000
```

Check this result here [in etherscan](https://etherscan.io/address/0x5abfec25f74cd88437631a7731906932776356f9).

## TODO

This is a _Work in Progress_. There are a number of ethereum elements to add.
Stay tuned!

* `[0x95]` - `eth-tx-receipt`:
  * Propose a script to get all receipts from a block and make a JSON array of them.
  * Support the input of this JSON array to form the `eth-tx-receipt-trie` (`[0x96]`) leaves, and the `eth-tx-receipt` objects.

* The rest of the IPLD ETH Types:
  * `[0x91]` - `eth-block-list`
  * `[0x98]` - `eth-storage-trie`
