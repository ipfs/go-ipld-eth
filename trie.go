package ipldeth

import (
	"fmt"

	//common "github.com/ethereum/go-ethereum/common"
	rlp "github.com/ethereum/go-ethereum/rlp"
	//trie "github.com/ethereum/go-ethereum/trie"

	cid "gx/ipfs/QmV5gPoRsjN1Gid3LMdNZTyfCtP2DsvqEbMAmz82RmmiGk/go-cid"
	node "gx/ipfs/QmYDscK7dmdo2GZ9aumS8s5auUUAH5mR1jvj5pYhWusfK7/go-ipld-node"
	mh "gx/ipfs/QmbZ6Cee2uHjG7hf19qLHppgKDRtaG4CVtMzdmK9VCVqLu/go-multihash"
)

type TrieNode struct {
	codec uint64
	Arr   []interface{}
	Val   []byte
}

func NewTrieNode(data []byte) *TrieNode {
	var i []interface{}
	err := rlp.DecodeBytes(data, &i)
	if err != nil {
		panic(err)
	}
	fmt.Println(i)
	switch len(i) {
	case 2:
		fmt.Println("Its a value")
	case 17:
		fmt.Println("its a shard")
	default:
		fmt.Println("who knows what this is")
	}

	return &TrieNode{
		Arr:   i,
		Val:   data,
		codec: MEthTxTrie,
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
	return b.Val
}

func (b *TrieNode) Resolve(p []string) (interface{}, []string, error) {
	return nil, nil, nil
}

func (b *TrieNode) ResolveLink(p []string) (*node.Link, []string, error) {
	return nil, nil, nil
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
	return []string{"VALUE NODE"}
}
