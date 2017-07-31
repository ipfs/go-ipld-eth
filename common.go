package ipldeth

import (
	"bytes"

	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"

	common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

// IPLD Codecs for Ethereum
const (
	MEthBlock           = 0x90
	MEthBlockList       = 0x91
	MEthTxTrie          = 0x92
	MEthTx              = 0x93
	MEthTxReceiptTrie   = 0x94
	MEthTxReceipt       = 0x95
	MEthStateTrie       = 0x96
	MEthAccountSnapshot = 0x97
	MEthStorageTrie     = 0x98
)

// rawdataToCid takes the desired coded and a slice of bytes
// and returns the proper cid of the object.
func rawdataToCid(codec uint64, rawdata []byte) *cid.Cid {
	c, err := cid.Prefix{
		Codec:    codec,
		Version:  1,
		MhType:   mh.KECCAK_256,
		MhLength: -1,
	}.Sum(rawdata)
	if err != nil {
		panic(err)
	}
	return c
}

// keccak256ToCid takes a keccak256 hash and returns its cid based on
// the codec given.
func keccak256ToCid(codec uint64, h []byte) *cid.Cid {
	buf, err := mh.Encode(h, mh.KECCAK_256)
	if err != nil {
		panic(err)
	}

	return cid.NewCidV1(codec, mh.Multihash(buf))
}

// TODO
// Add documentation
func castCommonHash(codec uint64, h common.Hash) *cid.Cid {
	mhash, err := mh.Encode(h[:], mh.KECCAK_256)
	if err != nil {
		panic(err)
	}

	return cid.NewCidV1(codec, mhash)
}

// TODO
// Add documentation
func getRLP(object interface{}) []byte {
	buf := new(bytes.Buffer)
	if err := rlp.Encode(buf, object); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

// localTrie wraps a go-ethereum trie and its underlying memory db.
// It contributes to the creation of the trie node objects.
type localTrie struct {
	db   *ethdb.MemDatabase
	trie *trie.Trie
}

// newLocalTrie initializes and returns a localTrie object
func newLocalTrie() *localTrie {
	var err error
	lt := &localTrie{}

	lt.db, err = ethdb.NewMemDatabase()
	if err != nil {
		panic(err)
	}

	lt.trie, err = trie.New(common.Hash{}, lt.db)
	if err != nil {
		panic(err)
	}

	return lt
}

// add receives the index of an object and its rawdata value
// and includes it into the localTrie
func (lt *localTrie) add(idx int, rawdata []byte) {
	key, err := rlp.EncodeToBytes(uint(idx))
	if err != nil {
		panic(err)
	}

	lt.trie.Update(key, rawdata)
}

// getKeys returns the stored keys of the memory database
// of the localTrie for further processing.
func (lt *localTrie) getKeys() [][]byte {
	var err error

	_, err = lt.trie.Commit()
	if err != nil {
		panic(err)
	}

	return lt.db.Keys()
}
