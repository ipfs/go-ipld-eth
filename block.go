package ipldeth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"
	mh "github.com/multiformats/go-multihash"

	types "github.com/ethereum/go-ethereum/core/types"
	rlp "github.com/ethereum/go-ethereum/rlp"
)

// EthBlock (eth-block, codec 0x90), represents an ethereum block header
type EthBlock struct {
	*types.Header

	cid *cid.Cid
}

// Static (compile time) check that EthBlock satisfies the node.Node interface.
var _ node.Node = (*EthBlock)(nil)

/*
  INPUT
*/

// FromBlockHeaderRLP takes an RLP message representing an ethereum block header,
// parses it, calculate its cid, and wraps it into the EthBlock struct for further processing.
func FromBlockHeaderRLP(r io.Reader) (*EthBlock, error) {
	ethBlock := &EthBlock{}

	// We will read this buffer twice
	var r1 bytes.Buffer
	r0 := io.TeeReader(r, &r1)

	// Parse the RLP into a geth types.Header object
	var h types.Header
	s := rlp.NewStream(r0, 0)

	err := s.Decode(&h)
	if err != nil {
		return nil, err
	}

	ethBlock.Header = &h

	// Now, let's create the cid
	c, err := cid.Prefix{
		Codec:    MEthBlock,
		Version:  1,
		MhType:   mh.KECCAK_256,
		MhLength: -1,
	}.Sum(r1.Bytes())
	if err != nil {
		panic(err)
	}

	ethBlock.cid = c

	return ethBlock, nil
}

/*
  OUTPUT
*/

// DecodeBlock takes raw binary data from IPFS and returns
// a block header for further processing.
func DecodeBlock(r io.Reader) (*EthBlock, error) {
	var h types.Header
	err := rlp.Decode(r, &h)
	if err != nil {
		return nil, err
	}

	return &EthBlock{Header: &h}, nil
}

// MarshalJSON processes the block header into readable JSON format.
func (b *EthBlock) MarshalJSON() ([]byte, error) {
	out := map[string]interface{}{
		"time":       b.Time,
		"bloom":      b.Bloom,
		"coinbase":   b.Coinbase,
		"difficulty": b.Difficulty,
		"extra":      b.Extra,
		"gaslimit":   b.GasLimit,
		"gasused":    b.GasUsed,
		"mixdigest":  b.MixDigest,
		"nonce":      b.Nonce,
		"number":     b.Number,
		"parent":     castCommonHash(b.ParentHash, MEthBlock),
		"receipts":   castCommonHash(b.ReceiptHash, MEthTxReceiptTrie),
		"root":       castCommonHash(b.Root, MEthStateTrie),
		"tx":         castCommonHash(b.TxHash, MEthTxTrie),
		"uncles":     castCommonHash(b.UncleHash, MEthBlockList),
	}
	return json.Marshal(out)
}

// Cid returns the cid of the block header.
func (b *EthBlock) Cid() *cid.Cid {
	return b.cid
}

// Parent returns the cid of the parent of the block.
func (b *EthBlock) Parent() *cid.Cid {
	return toCid(MEthBlock, b.ParentHash.Bytes())
}

// Tx returns the cid of the transactionsTrie root of the block.
func (b *EthBlock) Tx() *cid.Cid {
	return castCommonHash(b.TxHash, MEthTxTrie)
}

// Links is a helper function that returns all links within this object
func (b *EthBlock) Links() []*node.Link {
	return []*node.Link{
		&node.Link{Cid: castCommonHash(b.ParentHash, MEthBlock)},
		&node.Link{Cid: castCommonHash(b.ReceiptHash, MEthTxReceiptTrie)},
		&node.Link{Cid: castCommonHash(b.Root, MEthStateTrie)},
		&node.Link{Cid: castCommonHash(b.TxHash, MEthTxTrie)},
		&node.Link{Cid: castCommonHash(b.UncleHash, MEthBlockList)},
	}
}

// Loggable returns in a map the type of IPLD Link.
func (b *EthBlock) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "eth_block",
	}
}

// RawData returns the binary of the RLP encode of the block header.
func (b *EthBlock) RawData() []byte {
	buf := new(bytes.Buffer)
	if err := rlp.Encode(buf, b); err != nil {
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
		return &node.Link{Cid: toCid(MEthTxTrie, b.TxHash.Bytes())}, p[1:], nil
	case "parent":
		return &node.Link{Cid: toCid(MEthBlock, b.ParentHash.Bytes())}, p[1:], nil
	default:
		return nil, nil, fmt.Errorf("no such link")
	}
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

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (b *EthBlock) Tree(p string, depth int) []string {
	return nil
}

func toCid(ctype uint64, h []byte) *cid.Cid {
	buf, err := mh.Encode(h, mh.KECCAK_256)
	if err != nil {
		panic(err)
	}

	return cid.NewCidV1(ctype, mh.Multihash(buf))
}

// String is a helper for output
func (b *EthBlock) String() string {
	return fmt.Sprintf("<EthBlock %s>", b.Cid())
}

/*
  GRAVEYARD
*/

// Copy will go away. It is here to comply with the interface.
func (b *EthBlock) Copy() node.Node {
	panic("dont use this yet")
}

// Size will go away. It is here to comply with the interface.
func (b *EthBlock) Size() (uint64, error) {
	return 0, nil
}

// Stat will go away. It is here to comply with the interface.
func (b *EthBlock) Stat() (*node.NodeStat, error) {
	return &node.NodeStat{}, nil
}
