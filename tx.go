package ipldeth

import (
	"bytes"
	"fmt"

	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-node"
	mh "github.com/multiformats/go-multihash"

	types "github.com/ethereum/go-ethereum/core/types"
)

type Tx struct {
	tx *types.Transaction
}

func (b *Tx) Cid() *cid.Cid {
	c, err := cid.Prefix{
		Codec:    cid.EthereumTx,
		Version:  1,
		MhType:   mh.KECCAK_256,
		MhLength: -1,
	}.Sum(b.RawData())
	if err != nil {
		panic(err)
	}
	return c
}

func (b *Tx) Copy() node.Node {
	panic("dont use this yet")
}

func (b *Tx) Links() []*node.Link {
	return nil
}

func (b *Tx) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "ethereum_block",
	}
}

func (b *Tx) RawData() []byte {
	buf := new(bytes.Buffer)
	if err := b.tx.EncodeRLP(buf); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func (b *Tx) Resolve(p []string) (interface{}, []string, error) {
	return nil, nil, nil
}

func (b *Tx) ResolveLink(p []string) (*node.Link, []string, error) {
	return nil, nil, nil
}

func (b *Tx) Size() (uint64, error) {
	panic("don't do size")
}

func (b *Tx) Stat() (*node.NodeStat, error) {
	return &node.NodeStat{}, nil
}

func (b *Tx) String() string {
	return fmt.Sprintf("<EthereumTx %s>", b.Cid())
}

func (b *Tx) Tree(p string, depth int) []string {
	return nil
}
