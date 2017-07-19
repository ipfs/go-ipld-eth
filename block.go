package ipldeth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	cid "gx/ipfs/QmTprEaAA2A9bst5XH7exuyi5KzNMK3SEDNN8rBDnKWcUS/go-cid"
	mh "gx/ipfs/QmU9a9NV9RdPNwZQDYd5uKsm6N6LJLSvLbywDDYFbaaC6P/go-multihash"
	node "gx/ipfs/QmYNyRZJBUYPNrLszFmrBrPJbsBh2vMsefz5gnDpB5M1P6/go-ipld-format"

	common "github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
	rlp "github.com/ethereum/go-ethereum/rlp"
	trie "github.com/ethereum/go-ethereum/trie"
)

// EthBlock will be renamed to EthHeader and refactorized, as per
// https://github.com/paritytech/parity/issues/4172#issue-200744099
// https://github.com/MetaMask/metamask-extension/issues/719#issuecomment-267457567
// https://github.com/ipld/js-ipld-eth-star/blob/master/eth-block/index.js
// eth-block (code 0x90), represents the block header
// TODO
// Activity to be performed after completing the first `golint`
type EthBlock struct {
	header *types.Header
}

// FromRlpBlockMessage takes an RLP message emitted by the BlockByHash ws API
// in go ethereum, and decodes it to return the block header, tx, tx-tries and ommers.
// TODO
// Refactor to block header to comply with eth-block (0x90)
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

	_, err = tr.Commit()
	if err != nil {
		return nil, err
	}

	var out []*TrieNode
	for _, nd := range d.vals {
		out = append(out, nd)
	}

	if len(out) == 0 {
		return []*TrieNode{{val: []byte{0x80}, codec: MEthTxTrie}}, nil
	}

	return out, nil
}

// DecodeBlock takes raw binary data and returns a block header for further processing.
func DecodeBlock(r io.Reader) (*EthBlock, error) {
	var h types.Header
	err := rlp.Decode(r, &h)
	if err != nil {
		return nil, err
	}

	return &EthBlock{&h}, nil
}

var _ node.Node = (*EthBlock)(nil)

// MarshalJSON processes the block header into readable JSON format.
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

// Cid returns the content identifier of the block header.
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

// Parent returns the content identifier of the parent of the block.
func (b *EthBlock) Parent() *cid.Cid {
	return toCid(MEthBlock, b.header.ParentHash.Bytes())
}

// Tx returns the content identifier of the transactionsTrie root of the block.
func (b *EthBlock) Tx() *cid.Cid {
	return castCommonHash(b.header.TxHash, MEthTxTrie)
}

// Copy is NOT IMPLEMENTED YET
// Should return a deep copy of this node.
// TODO
// TBD how deep we want to copy this node.
func (b *EthBlock) Copy() node.Node {
	panic("dont use this yet")
}

// Links is a helper function that returns all links within this object
func (b *EthBlock) Links() []*node.Link {
	return []*node.Link{
		&node.Link{Cid: castCommonHash(b.header.ParentHash, MEthBlock)},
		&node.Link{Cid: castCommonHash(b.header.ReceiptHash, MEthTxReceiptTrie)},
		&node.Link{Cid: castCommonHash(b.header.Root, MEthStateTrie)},
		&node.Link{Cid: castCommonHash(b.header.TxHash, MEthTxTrie)},
		&node.Link{Cid: castCommonHash(b.header.UncleHash, MEthBlockList)},
	}
}

// Loggable returns in a map the type of IPLD Link.
func (b *EthBlock) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "ethereum_block",
	}
}

// RawData returns the binary of the RLP encode of the block header.
func (b *EthBlock) RawData() []byte {
	buf := new(bytes.Buffer)
	if err := rlp.Encode(buf, b.header); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

// Resolve resolves a path through this node, stopping at any link boundary
// and returning the object found as well as the remaining path to traverse
func (b *EthBlock) Resolve(p []string) (interface{}, []string, error) {
	if len(p) == 0 {
		return b, nil, nil
	}

	switch p[0] {
	case "tx":
		return &node.Link{Cid: toCid(MEthTxTrie, b.header.TxHash.Bytes())}, p[1:], nil
	case "parent":
		return &node.Link{Cid: toCid(MEthBlock, b.header.ParentHash.Bytes())}, p[1:], nil
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

// ResolveLink is a helper function that calls resolve and asserts the
// output is a link
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

// Size returns the size in bytes of the serialized object
func (b *EthBlock) Size() (uint64, error) {
	// TODO:
	return 0, nil
}

// Stat helps this struct to comply with the Node interface
// TODO: not sure if stat deserves to stay
func (b *EthBlock) Stat() (*node.NodeStat, error) {
	return &node.NodeStat{}, nil
}

// String is a helper for output
func (b *EthBlock) String() string {
	return fmt.Sprintf("<EthBlock %s>", b.Cid())
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
// TODO
// Implement
func (b *EthBlock) Tree(p string, depth int) []string {
	return nil
}
