package ipldeth

import (
	"encoding/json"
	"fmt"
	"strconv"

	//common "github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
	rlp "github.com/ethereum/go-ethereum/rlp"

	cid "gx/ipfs/QmV5gPoRsjN1Gid3LMdNZTyfCtP2DsvqEbMAmz82RmmiGk/go-cid"
	node "gx/ipfs/QmYDscK7dmdo2GZ9aumS8s5auUUAH5mR1jvj5pYhWusfK7/go-ipld-node"
	mh "gx/ipfs/QmbZ6Cee2uHjG7hf19qLHppgKDRtaG4CVtMzdmK9VCVqLu/go-multihash"
)

type TrieNode struct {
	codec uint64
	Arr   []interface{}
	val   []byte
}

func NewTrieNode(data []byte) (node.Node, error) {
	var i []interface{}
	err := rlp.DecodeBytes(data, &i)
	if err != nil {
		return nil, err
	}

	switch len(i) {
	case 2:
		var out interface{}
		if err := rlp.DecodeBytes(i[0].([]byte), &out); err != nil {
			return nil, err
		}
		var t types.Transaction
		if err := rlp.DecodeBytes(i[1].([]byte), &t); err != nil {
			return nil, err
		}
		return &Tx{&t}, nil
	case 17:
		return &TrieNode{
			Arr:   i,
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

func (b *TrieNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Arr)
}

func (b *TrieNode) Copy() node.Node {
	panic("dont use this yet")
}

func (b *TrieNode) Links() []*node.Link {
	return nil
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

	bytes, ok := b.Arr[i].([]byte)
	if !ok {
		return nil, nil, fmt.Errorf("expected trie array element to be bytes")
	}

	if len(bytes) == 32 {
		// probably a hash
		return &node.Link{Cid: toCid(MEthTxTrie, bytes)}, p[1:], nil
	}

	fmt.Println(len(b.Arr[i].([]byte)))
	return b.Arr[i], p[1:], nil
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
