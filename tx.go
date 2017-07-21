package ipldeth

import (
	"bytes"
	"fmt"

	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"
	mh "github.com/multiformats/go-multihash"

	types "github.com/ethereum/go-ethereum/core/types"
	rlp "github.com/ethereum/go-ethereum/rlp"
)

// Tx (eth-tx codec 0x93) represents an ethereum transaction
type Tx struct {
	tx *types.Transaction
}

// Cid returns the cid of the transaction.
func (t *Tx) Cid() *cid.Cid {
	c, err := cid.Prefix{
		Codec:    MEthTx,
		Version:  1,
		MhType:   mh.KECCAK_256,
		MhLength: -1,
	}.Sum(t.RawData())
	if err != nil {
		panic(err)
	}
	return c
}

// ParseTx takes raw binary data and returns a transaction for further processing.
func ParseTx(data []byte) (*Tx, error) {
	var t types.Transaction
	err := rlp.DecodeBytes(data, &t)
	if err != nil {
		return nil, err
	}
	return &Tx{&t}, nil
}

// MarshalJSON processes the transaction into readable JSON format.
func (t *Tx) MarshalJSON() ([]byte, error) {
	return t.tx.MarshalJSON()
}

// Copy is NOT IMPLEMENTED YET
func (t *Tx) Copy() node.Node {
	panic("dont use this yet")
}

// Links is a helper function that returns all links within this object
func (t *Tx) Links() []*node.Link {
	return nil
}

// Loggable returns in a map the type of IPLD Link.
func (t *Tx) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "ethereum_transaction",
	}
}

// RawData returns the binary of the RLP encode of the transaction.
func (t *Tx) RawData() []byte {
	buf := new(bytes.Buffer)
	if err := t.tx.EncodeRLP(buf); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

// Resolve resolves a path through this node, stopping at any link boundary
// and returning the object found as well as the remaining path to traverse
func (t *Tx) Resolve(p []string) (interface{}, []string, error) {
	if len(p) == 0 {
		return t, nil, nil
	}

	switch p[0] {
	case "nonce":
		return t.tx.Nonce(), p[1:], nil
	case "gasPrice":
		return t.tx.GasPrice(), p[1:], nil
	case "gas":
		return t.tx.Gas(), p[1:], nil
	case "toAddress":
		return t.tx.To(), p[1:], nil
	case "value":
		return t.tx.Value(), p[1:], nil
	case "data":
		return t.tx.Data(), p[1:], nil
	default:
		return nil, nil, fmt.Errorf("no such link")
	}
}

// ResolveLink is a helper function that calls resolve and asserts the
// output is a link
func (t *Tx) ResolveLink(p []string) (*node.Link, []string, error) {
	obj, rest, err := t.Resolve(p)
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
func (t *Tx) Size() (uint64, error) {
	return uint64(t.tx.Size().Int64()), nil
}

// Stat helps this struct to comply with the Node interface
func (t *Tx) Stat() (*node.NodeStat, error) {
	return &node.NodeStat{}, nil
}

// String is a helper for output
func (t *Tx) String() string {
	return fmt.Sprintf("<EthereumTx %s>", t.Cid())
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (t *Tx) Tree(p string, depth int) []string {
	return []string{"toAddress", "value", "data", "nonce", "gasPrice", "gas"}
}

// BaseTx returns a go-ethereum/types.Transaction pointer to the object
func (t *Tx) BaseTx() *types.Transaction {
	return t.tx
}
