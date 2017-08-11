package ipldeth

import (
	"bytes"
	"fmt"

	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"

	types "github.com/ethereum/go-ethereum/core/types"
	rlp "github.com/ethereum/go-ethereum/rlp"
)

// EthTx (eth-tx codec 0x93) represents an ethereum transaction
type EthTx struct {
	*types.Transaction

	cid     *cid.Cid
	rawdata []byte
}

// Static (compile time) check that EthTx satisfies the node.Node interface.
var _ node.Node = (*EthTx)(nil)

/*
  INPUT
*/

// NewTx computes the cid and rlp-encodes a types.Transaction object
// returning a proper EthTx node
func NewTx(t *types.Transaction) *EthTx {
	buf := new(bytes.Buffer)
	if err := t.EncodeRLP(buf); err != nil {
		panic(err)
	}
	rawdata := buf.Bytes()

	return &EthTx{
		Transaction: t,
		cid:         rawdataToCid(MEthTx, rawdata),
		rawdata:     rawdata,
	}
}

/*
 OUTPUT
*/

// DecodeEthTx takes a cid and its raw binary data
// from IPFS and returns an EthTx object for further processing.
func DecodeEthTx(c *cid.Cid, b []byte) (*EthTx, error) {
	var t types.Transaction
	err := rlp.DecodeBytes(b, &t)
	if err != nil {
		return nil, err
	}

	return &EthTx{
		Transaction: &t,
		cid:         c,
		rawdata:     b,
	}, nil
}

/*
  Block INTERFACE
*/

// RawData returns the binary of the RLP encode of the transaction.
func (t *EthTx) RawData() []byte {
	return t.rawdata
}

// Cid returns the cid of the transaction.
func (t *EthTx) Cid() *cid.Cid {
	return t.cid
}

// String is a helper for output
func (t *EthTx) String() string {
	return fmt.Sprintf("<EthereumTx %s>", t.cid)
}

// Loggable returns in a map the type of IPLD Link.
func (t *EthTx) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "eth-tx",
	}
}

/*
  Node INTERFACE
*/

// Resolve resolves a path through this node, stopping at any link boundary
// and returning the object found as well as the remaining path to traverse
func (t *EthTx) Resolve(p []string) (interface{}, []string, error) {
	if len(p) == 0 {
		return t, nil, nil
	}

	if len(p) > 1 {
		return nil, nil, fmt.Errorf("unexpected path elements past %s", p[0])
	}

	switch p[0] {
	case "nonce":
		return t.Nonce(), p[1:], nil
	case "gasPrice":
		return t.GasPrice(), p[1:], nil
	case "gas":
		return t.Gas(), p[1:], nil
	case "toAddress":
		return t.To(), p[1:], nil
	case "value":
		return t.Value(), p[1:], nil
	case "data":
		return t.Data(), p[1:], nil
	default:
		return nil, nil, fmt.Errorf("no such link")
	}
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (t *EthTx) Tree(p string, depth int) []string {
	return []string{"toAddress", "value", "data", "nonce", "gasPrice", "gas"}
}

// ResolveLink is a helper function that calls resolve and asserts the
// output is a link
func (t *EthTx) ResolveLink(p []string) (*node.Link, []string, error) {
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

// Copy will go away. It is here to comply with the interface.
func (t *EthTx) Copy() node.Node {
	panic("dont use this yet")
}

// Links is a helper function that returns all links within this object
func (t *EthTx) Links() []*node.Link {
	return nil
}

// Stat will go away. It is here to comply with the interface.
func (t *EthTx) Stat() (*node.NodeStat, error) {
	return &node.NodeStat{}, nil
}

// Size will go away. It is here to comply with the interface.
func (t *EthTx) Size() (uint64, error) {
	return uint64(t.Transaction.Size().Int64()), nil
}

/*
  EthTx functions
*/

// MarshalJSON processes the transaction into readable JSON format.
func (t *EthTx) MarshalJSON() ([]byte, error) {
	return t.Transaction.MarshalJSON()
}
