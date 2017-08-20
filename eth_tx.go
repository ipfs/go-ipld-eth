package ipldeth

import (
	"bytes"
	"encoding/json"
	"fmt"

	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"

	hexutil "github.com/ethereum/go-ethereum/common/hexutil"
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

	case "gas":
		return t.Gas(), nil, nil
	case "gasPrice":
		return t.GasPrice(), nil, nil
	case "input":
		return fmt.Sprintf("%x", t.Data()), nil, nil
	case "nonce":
		return t.Nonce(), nil, nil
	case "r":
		_, r, _ := t.RawSignatureValues()
		return hexutil.EncodeBig(r), nil, nil
	case "s":
		_, _, s := t.RawSignatureValues()
		return hexutil.EncodeBig(s), nil, nil
	case "toAddress":
		return t.To(), nil, nil
	case "v":
		v, _, _ := t.RawSignatureValues()
		return hexutil.EncodeBig(v), nil, nil
	case "value":
		return hexutil.EncodeBig(t.Value()), nil, nil
	default:
		return nil, nil, fmt.Errorf("no such link")
	}
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (t *EthTx) Tree(p string, depth int) []string {
	if p != "" || depth == 0 {
		return nil
	}
	return []string{"gas", "gasPrice", "input", "nonce", "r", "s", "toAddress", "v", "value"}
}

// ResolveLink is a helper function that calls resolve and asserts the
// output is a link
func (t *EthTx) ResolveLink(p []string) (*node.Link, []string, error) {
	obj, rest, err := t.Resolve(p)
	if err != nil {
		return nil, nil, err
	}

	if lnk, ok := obj.(*node.Link); ok {
		return lnk, rest, nil
	}

	return nil, nil, fmt.Errorf("resolved item was not a link")
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
	v, r, s := t.RawSignatureValues()

	out := map[string]interface{}{
		"gas":       t.Gas(),
		"gasPrice":  t.GasPrice(),
		"input":     fmt.Sprintf("%x", t.Data()),
		"nonce":     t.Nonce(),
		"r":         hexutil.EncodeBig(r),
		"s":         hexutil.EncodeBig(s),
		"toAddress": t.To(),
		"v":         hexutil.EncodeBig(v),
		"value":     hexutil.EncodeBig(t.Value()),
	}
	return json.Marshal(out)
}
