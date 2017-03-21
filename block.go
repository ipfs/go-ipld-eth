package ipldeth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-node"
	mh "github.com/multiformats/go-multihash"

	common "github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
	rlp "github.com/ethereum/go-ethereum/rlp"
	trie "github.com/ethereum/go-ethereum/trie"
)

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

type EthBlock struct {
	header *types.Header
}

func FromRlpBlockMessage(r io.Reader) (*EthBlock, []*Tx, []*TrieNode, []*EthBlock, error) {
	var b types.Block
	s := rlp.NewStream(r, 0)
	err := b.DecodeRLP(s)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	var txs []*Tx
	for _, tx := range b.Transactions() {
		txs = append(txs, &Tx{tx})
	}

	triends, err := buildTreeFromTxs(txs)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var uncles []*EthBlock
	for _, u := range b.Uncles() {
		uncles = append(uncles, &EthBlock{u})
	}

	return &EthBlock{b.Header()}, txs, triends, uncles, nil
}

type db struct {
	vals map[string]*TrieNode
}

func (db *db) Get(k []byte) ([]byte, error) {
	mhval, err := mh.Encode(k, mh.KECCAK_256)
	if err != nil {
		return nil, err
	}

	h, err := mh.Cast(mhval)
	if err != nil {
		return nil, err
	}

	c := cid.NewCidV1(MEthTxTrie, h)

	out, ok := db.vals[c.KeyString()]
	if !ok {
		return nil, nil
	}
	return out.RawData(), nil
}

func (db *db) Put(k []byte, val []byte) error {
	mval := make([]byte, len(val))
	copy(mval, val)
	tn := &TrieNode{
		codec: MEthTxTrie,
		val:   mval,
	}
	db.vals[tn.Cid().KeyString()] = tn
	return nil
}

func newdb() *db {
	return &db{make(map[string]*TrieNode)}
}

func buildTreeFromTxs(txs []*Tx) ([]*TrieNode, error) {
	d := newdb()
	tr, err := trie.New(common.Hash{}, d)
	if err != nil {
		return nil, err
	}

	for i, tx := range txs {
		key, err := rlp.EncodeToBytes(uint(i))
		if err != nil {
			return nil, err
		}

		tr.Update(key, tx.RawData())
	}

	tr.Commit()
	var out []*TrieNode
	for _, nd := range d.vals {
		out = append(out, nd)
	}

	return out, nil
}

func DecodeBlock(r io.Reader) (*EthBlock, error) {
	var h types.Header
	err := rlp.Decode(r, &h)
	if err != nil {
		return nil, err
	}

	return &EthBlock{&h}, nil
}

var _ node.Node = (*EthBlock)(nil)

func (b *EthBlock) MarshalJSON() ([]byte, error) {
	out := map[string]interface{}{
		"time":       b.header.Time,
		"bloom":      b.header.Bloom,
		"coinbase":   b.header.Coinbase,
		"difficulty": b.header.Difficulty,
		"extra":      b.header.Extra,
		"gaslimit":   b.header.GasLimit,
		"gasused":    b.header.GasUsed,
		"mixdigest":  b.header.MixDigest,
		"nonce":      b.header.Nonce,
		"number":     b.header.Number,
		"parent":     castCommonHash(b.header.ParentHash, MEthBlock),
		"receipts":   castCommonHash(b.header.ReceiptHash, MEthTxReceiptTrie),
		"root":       castCommonHash(b.header.Root, MEthStateTrie),
		"tx":         castCommonHash(b.header.TxHash, MEthTxTrie),
		"uncles":     castCommonHash(b.header.UncleHash, MEthBlockList),
	}
	return json.Marshal(out)
}

func (b *EthBlock) Cid() *cid.Cid {
	c, err := cid.Prefix{
		Codec:    MEthBlock,
		Version:  1,
		MhType:   mh.KECCAK_256,
		MhLength: -1,
	}.Sum(b.RawData())
	if err != nil {
		panic(err)
	}
	return c
}

func (b *EthBlock) Copy() node.Node {
	panic("dont use this yet")
}

func (b *EthBlock) Links() []*node.Link {
	return nil
}

func (b *EthBlock) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "ethereum_block",
	}
}

func (b *EthBlock) RawData() []byte {
	buf := new(bytes.Buffer)
	if err := rlp.Encode(buf, b.header); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func (b *EthBlock) Resolve(p []string) (interface{}, []string, error) {
	if len(p) == 0 {
		return b, nil, nil
	}

	switch p[0] {
	case "tx":
		return &node.Link{Cid: toCid(MEthTxTrie, b.header.TxHash.Bytes())}, p[1:], nil
	default:
		return nil, nil, fmt.Errorf("no such link")
	}
}

func toCid(ctype uint64, h []byte) *cid.Cid {
	buf, err := mh.Encode(h, mh.KECCAK_256)
	if err != nil {
		panic(err)
	}

	return cid.NewCidV1(ctype, mh.Multihash(buf))
}

func (b *EthBlock) ResolveLink(p []string) (*node.Link, []string, error) {
	obj, rest, err := b.Resolve(p)
	if err != nil {
		return nil, nil, err
	}

	if lnk, ok := obj.(*node.Link); ok {
		return lnk, rest, nil
	}

	return nil, nil, fmt.Errorf("resolved item was not a link")
}

func (b *EthBlock) Size() (uint64, error) {
	panic("don't do size")
}

func (b *EthBlock) Stat() (*node.NodeStat, error) {
	return &node.NodeStat{}, nil
}

func (b *EthBlock) String() string {
	return fmt.Sprintf("<EthBlock %s>", b.Cid())
}

func (b *EthBlock) Tree(p string, depth int) []string {
	return nil
}
