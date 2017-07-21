package ipldeth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	types "github.com/ethereum/go-ethereum/core/types"
	rlp "github.com/ethereum/go-ethereum/rlp"

	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"
	mh "github.com/multiformats/go-multihash"
)

// TrieNode represents a generic ethereum trie object
type TrieNode struct {
	codec uint64
	Arr   []interface{}
	val   []byte
}

// NewTrieNode returns a formatted trie node from a raw dump
func NewTrieNode(data []byte) (node.Node, error) {
	if bytes.Equal(data, []byte{0x80}) {
		return &TrieNode{val: []byte{0x80}, codec: MEthTxTrie}, nil
	}

	var i []interface{}
	err := rlp.DecodeBytes(data, &i)
	if err != nil {
		return nil, err
	}

	switch len(i) {
	case 2:
		key := i[0].([]byte)

		valb := i[1].([]byte)

		var val interface{}
		if len(valb) == 32 {
			val = toCid(MEthTxTrie, valb)
		} else {
			var t types.Transaction
			if err := rlp.DecodeBytes(i[1].([]byte), &t); err != nil {
				return nil, err
			}
			val = &Tx{&t}
		}
		return &TrieNode{
			Arr:   []interface{}{key, val},
			val:   data,
			codec: MEthTxTrie,
		}, nil
	case 17:
		var parsed []interface{}
		for _, v := range i {
			bv := v.([]byte)
			switch len(bv) {
			case 0:
				parsed = append(parsed, nil)
			case 32:
				parsed = append(parsed, toCid(MEthTxTrie, bv))
			default:
				return nil, fmt.Errorf("unrecognized object in trie: %v", bv)
			}
		}
		return &TrieNode{
			Arr:   parsed,
			val:   data,
			codec: MEthTxTrie,
		}, nil
	default:
		return nil, fmt.Errorf("unknown trie node type")
	}
}

// Cid returns the cid of the trie
func (tn *TrieNode) Cid() *cid.Cid {
	c, err := cid.Prefix{
		Codec:    tn.codec,
		Version:  1,
		MhType:   mh.KECCAK_256,
		MhLength: -1,
	}.Sum(tn.RawData())
	if err != nil {
		panic(err)
	}
	return c
}

// HexHash returns the hex hash of the trie
func (tn *TrieNode) HexHash() string {
	return fmt.Sprintf("%x", tn.Cid().Bytes()[4:])
}

// MarshalJSON processes the trie node into readable JSON format.
func (tn *TrieNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(tn.Arr)
}

// Copy is NOT IMPLEMENTED YET
func (tn *TrieNode) Copy() node.Node {
	panic("dont use this yet")
}

// Links is a helper function that returns all links within this object
func (tn *TrieNode) Links() []*node.Link {
	var out []*node.Link
	for _, i := range tn.Arr {
		c, ok := i.(*cid.Cid)
		if ok {
			out = append(out, &node.Link{Cid: c})
		}
	}
	return out
}

// Loggable returns in a map the type of IPLD Link.
func (tn *TrieNode) Loggable() map[string]interface{} {
	// TODO
	// Should change the value based on the codec?
	return map[string]interface{}{
		"type": "ethereum_trie",
	}
}

// RawData returns the raw data of the trie node
func (tn *TrieNode) RawData() []byte {
	return tn.val
}

// Resolve resolves a path through this node, stopping at any link boundary
// and returning the object found as well as the remaining path to traverse
func (tn *TrieNode) Resolve(p []string) (interface{}, []string, error) {
	if len(p) == 0 {
		return tn, nil, nil
	}

	i, err := strconv.Atoi(p[0])
	if err != nil {
		return nil, nil, fmt.Errorf("expected array index to trie: %s", err)
	}

	if i < 0 || i >= len(tn.Arr) {
		return nil, nil, fmt.Errorf("index in trie out of range")
	}

	switch obj := tn.Arr[i].(type) {
	case *cid.Cid:
		return &node.Link{Cid: obj}, p[1:], nil
	case *Tx:
		return obj, p[1:], nil
	default:
		return nil, nil, fmt.Errorf("unexpected object type in trie")
	}

}

// ResolveLink is a helper function that calls resolve and asserts the
// output is a link
func (tn *TrieNode) ResolveLink(p []string) (*node.Link, []string, error) {
	obj, rest, err := tn.Resolve(p)
	if err != nil {
		return nil, nil, err
	}

	lnk, ok := obj.(*node.Link)
	if !ok {
		return nil, nil, fmt.Errorf("was not a link")
	}

	return lnk, rest, nil
}

// Size returns the size in bytes of the serialized object
func (tn *TrieNode) Size() (uint64, error) {
	panic("don't do size")
}

// Stat helps this struct to comply with the Node interface
func (tn *TrieNode) Stat() (*node.NodeStat, error) {
	return &node.NodeStat{}, nil
}

// String is a helper for output
func (tn *TrieNode) String() string {
	return fmt.Sprintf("<EthereumTrieNode %s>", tn.Cid())
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (tn *TrieNode) Tree(p string, depth int) []string {
	if p != "" {
		return nil
	}
	if depth > 0 {
		return nil
	}

	if len(tn.Arr) == 17 {
		var out []string
		for i, v := range tn.Arr {
			if len(v.([]byte)) == 0 {
				out = append(out, fmt.Sprintf("%x", i))
			}
		}
		return out
	}

	// TODO: not sure what to put here. Most of the 'keys' dont seem to be human readable
	return []string{"VALUE NODE"}
}
