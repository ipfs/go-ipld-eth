package ipldeth

import (
	"bytes"
	"fmt"
	"strings"

	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"
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

// rawdataToCid takes the desired codec and a slice of bytes
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

// commonHashToCid takes a go-ethereum common.Hash and returns its
// cid based on the codec given,
func commonHashToCid(codec uint64, h common.Hash) *cid.Cid {
	mhash, err := mh.Encode(h[:], mh.KECCAK_256)
	if err != nil {
		panic(err)
	}

	return cid.NewCidV1(codec, mhash)
}

// getRLP encodes the given object to RLP returning its bytes.
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

// rootHash returns the computed trie root.
// Useful for sanity checks on parsed data.
func (lt *localTrie) rootHash() []byte {
	return lt.trie.Hash().Bytes()
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

// decodeTrieNode takes a raw RLP-encoded trie node, returning its kind
// (branch, extension or leaf) and parsed data for further processing.
func decodeTrieNode(b []byte) (string, []interface{}, error) {
	var i []interface{}
	err := rlp.DecodeBytes(b, &i)
	if err != nil {
		return "", nil, err
	}

	switch len(i) {
	case 2:
		return decodeTrieLeafExtensionNode(i)
	case 17:
		return "branch", i, nil

	default:
		return "", nil, fmt.Errorf("unknown trie node type")
	}
}

// decodeTrieLeafExtensionNode takes a compact key, and returns its nodeKind and value.
func decodeTrieLeafExtensionNode(i []interface{}) (string, []interface{}, error) {
	first := i[0].([]byte)
	last := i[1].([]byte)

	switch first[0] / 16 {
	case '\x00':
		return "extension", []interface{}{
			nibbleToByte(first)[2:],
			last,
		}, nil
	case '\x01':
		return "extension", []interface{}{
			nibbleToByte(first)[1:],
			last,
		}, nil
	case '\x02':
		return "leaf", []interface{}{
			nibbleToByte(first)[2:],
			last,
		}, nil
	case '\x03':
		return "leaf", []interface{}{
			nibbleToByte(first)[1:],
			last,
		}, nil
	default:
		return "", nil, fmt.Errorf("unknown hex prefix")
	}
}

// nibbleToByte expands the nibbles of a byte slice into their own bytes.
func nibbleToByte(k []byte) []byte {
	var out []byte

	for _, b := range k {
		out = append(out, b/16)
		out = append(out, b%16)
	}

	return out
}

// resolveTriePath takes a trie path, and the nodeKind and elements of a trie node.
// After validating and normalizing the received path, it will resolve in the node
// whether there is a node link (with or without a rest of the path) or an error.
func resolveTriePath(p []string, nodeKind string, elements []interface{}) (interface{}, []string, error) {
	p, err := validateTriePath(p, getTxFields())
	if err != nil {
		return nil, nil, err
	}

	switch nodeKind {
	case "extension":
		nibblesCount := checkPathNibbles(elements[0].([]byte), p)
		if nibblesCount == -1 {
			return nil, nil, fmt.Errorf("no such link in this extension")
		}
		return &node.Link{Cid: elements[1].(*cid.Cid)}, p[nibblesCount:], nil
	case "branch":
		child := elements[getHexIndex(p[0])]
		if child != nil {
			return &node.Link{Cid: child.(*cid.Cid)}, p[1:], nil
		}
		return nil, nil, fmt.Errorf("no such link in this branch")
	default:
		return nil, nil, fmt.Errorf("nodeKind case not implemented")
	}
}

// validateTriePath takes a trie path, checking whether each element represents
// an hexadecimal character, and returns a slice of one hex character elements,
// allowing the input of paths such as /b/0d010/1 /0/1/1/b /cc001d4 possible.
func validateTriePath(p []string, specialFields map[string]interface{}) ([]string, error) {
	var (
		testString string
		output     []string
	)

	//
	lastValue := p[len(p)-1]
	if _, ok := specialFields[lastValue]; ok {
		// Remove this lastValue and add it after the validation.
		// Examples of lastValue: nonce, gasPrice for txs. balance for states.
		p = p[:len(p)-1]
	} else {
		lastValue = ""
	}

	for _, v := range p {
		if v == "" {
			return nil, fmt.Errorf("Unexpected blank element in path")
		}
		testString += v
	}

	testString = strings.ToLower(testString)

	for _, v := range testString {
		c := byte(v)

		switch {
		case '0' <= c && c <= '9':
			fallthrough
		case 'a' <= c && c <= 'f':
			output = append(output, string(c))
		default:
			return nil, fmt.Errorf("Unexpected character in path: %x", c)
		}
	}

	// Recover the last value
	if lastValue != "" {
		output = append(output, lastValue)
	}

	return output, nil
}

// checkPathNibbles tests whether the given path can resolve the trie node
// element key, returning the number of nibbles the key has if succeed.
func checkPathNibbles(nibbles []byte, p []string) int {
	if len(p) < len(nibbles) {
		return -1
	}

	for i, n := range nibbles {
		if p[i] != fmt.Sprintf("%x", n) {
			return -1
		}
	}

	return len(nibbles)
}

// getHexIndex returns to you the integer 0 - 15 equivalent to your
// string character if applicable, or -1 otherwise.
func getHexIndex(s string) int {
	if len(s) != 1 {
		return -1
	}

	c := byte(s[0])
	switch {
	case '0' <= c && c <= '9':
		return int(c - '0')
	case 'a' <= c && c <= 'f':
		return int(c - 'a' + 10)
	}

	return -1
}
