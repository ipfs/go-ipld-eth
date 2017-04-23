package ipldeth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	//common "github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
	rlp "github.com/ethereum/go-ethereum/rlp"

	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"
	mh "github.com/multiformats/go-multihash"
)

type TrieNode struct {
	codec uint64
	Arr   []interface{}
	val   []byte
}

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

func (b *TrieNode) Cid() *cid.Cid {
	c, err := cid.Prefix{
		Codec:    b.codec,
		Version:  1,
		MhType:   mh.KECCAK_256,
		MhLength: -1,
	}.Sum(b.RawData())
	if err != nil {
		panic(err)
	}
	return c
}

func (tn *TrieNode) HexHash() string {
	return fmt.Sprintf("%x", tn.Cid().Bytes()[4:])
}

func (b *TrieNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Arr)
}

func (b *TrieNode) Copy() node.Node {
	panic("dont use this yet")
}

func (b *TrieNode) Links() []*node.Link {
	var out []*node.Link
	for _, i := range b.Arr {
		c, ok := i.(*cid.Cid)
		if ok {
			out = append(out, &node.Link{Cid: c})
		}
	}
	return out
}

func (b *TrieNode) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "ethereum_block",
	}
}

func (b *TrieNode) RawData() []byte {
	return b.val
}

func (b *TrieNode) Resolve(p []string) (interface{}, []string, error) {
	if len(p) == 0 {
		return b, nil, nil
	}

	i, err := strconv.Atoi(p[0])
	if err != nil {
		return nil, nil, fmt.Errorf("expected array index to trie: %s", err)
	}

	if i < 0 || i >= len(b.Arr) {
		return nil, nil, fmt.Errorf("index in trie out of range")
	}

	switch obj := b.Arr[i].(type) {
	case *cid.Cid:
		return &node.Link{Cid: obj}, p[1:], nil
	case *Tx:
		return obj, p[1:], nil
	default:
		return nil, nil, fmt.Errorf("unexpected object type in trie")
	}

}

func (b *TrieNode) ResolveLink(p []string) (*node.Link, []string, error) {
	obj, rest, err := b.Resolve(p)
	if err != nil {
		return nil, nil, err
	}

	lnk, ok := obj.(*node.Link)
	if !ok {
		return nil, nil, fmt.Errorf("was not a link")
	}

	return lnk, rest, nil
}

func (b *TrieNode) Size() (uint64, error) {
	panic("don't do size")
}

func (b *TrieNode) Stat() (*node.NodeStat, error) {
	return &node.NodeStat{}, nil
}

func (b *TrieNode) String() string {
	return fmt.Sprintf("<EthereumTrieNode %s>", b.Cid())
}

func (b *TrieNode) Tree(p string, depth int) []string {
	if p != "" {
		return nil
	}
	if depth > 0 {
		return nil
	}

	if len(b.Arr) == 17 {
		var out []string
		for i, v := range b.Arr {
			if len(v.([]byte)) == 0 {
				out = append(out, fmt.Sprintf("%x", i))
			}
		}
		return out
	}

	// TODO: not sure what to put here. Most of the 'keys' dont seem to be human readable
	return []string{"VALUE NODE"}
}
