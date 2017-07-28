package ipldeth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"

	"github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
	rlp "github.com/ethereum/go-ethereum/rlp"
)

// EthBlock (eth-block, codec 0x90), represents an ethereum block header
type EthBlock struct {
	*types.Header

	cid     *cid.Cid
	rawdata []byte
}

// Static (compile time) check that EthBlock satisfies the node.Node interface.
var _ node.Node = (*EthBlock)(nil)

/*
  INPUT
*/

// FromBlockRLP takes an RLP message representing
// an ethereum block header or body (header, uncles and txs)
// to return it as an slice of IPLD nodes for further processing.
func FromBlockRLP(r io.Reader) (*EthBlock, error) {
	// We may want to use this stream several times
	rawdata, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Let's try to decode the received element as a block body
	var decodedBlock types.Block
	err = rlp.Decode(bytes.NewBuffer(rawdata), &decodedBlock)
	if err != nil {
		if err.Error()[:41] != "rlp: expected input list for types.Header" {
			return nil, err
		}

		// Maybe it is just a header... (body sans uncles and txs)
		var decodedHeader types.Header
		err := rlp.Decode(bytes.NewBuffer(rawdata), &decodedHeader)
		if err != nil {
			return nil, err
		}

		// It was a header
		return &EthBlock{
			Header:  &decodedHeader,
			cid:     rawdataToCid(MEthBlock, rawdata),
			rawdata: rawdata,
		}, nil
	}

	// This is a block body (header + uncles + txs)
	// We'll extract the header bits here
	headerRawData := getRLP(decodedBlock.Header())
	ethBlock := &EthBlock{
		Header:  decodedBlock.Header(),
		cid:     rawdataToCid(MEthBlock, headerRawData),
		rawdata: headerRawData,
	}

	// TODO
	// eth-block-list, eth-tx, eth-tx-trie

	return ethBlock, nil
}

// FromBlockJSON takes the output of an ethereum client JSON API
// (i.e. parity or geth) and returns a slice of IPLD nodes.
func FromBlockJSON(r io.Reader) (*EthBlock, error) {
	var obj objJSONBlock
	dec := json.NewDecoder(r)
	err := dec.Decode(&obj)
	if err != nil {
		return nil, err
	}

	headerRawData := getRLP(obj.Result.Header)
	ethBlock := &EthBlock{
		Header:  &obj.Result.Header,
		cid:     rawdataToCid(MEthBlock, headerRawData),
		rawdata: headerRawData,
	}

	// TODO
	// eth-block-list, eth-tx, eth-tx-trie

	return ethBlock, nil
}

/*
  OUTPUT
*/

// DecodeBlockHeader takes raw binary data from IPFS and returns
// a block header for further processing.
func DecodeBlockHeader(c *cid.Cid, b []byte) (*EthBlock, error) {
	var h types.Header
	err := rlp.Decode(bytes.NewReader(b), &h)
	if err != nil {
		return nil, err
	}

	return &EthBlock{
		Header:  &h,
		cid:     c,
		rawdata: b,
	}, nil
}

/*
  Block INTERFACE
*/

// RawData returns the binary of the RLP encode of the block header.
func (b *EthBlock) RawData() []byte {
	return b.rawdata
}

// Cid returns the cid of the block header.
func (b *EthBlock) Cid() *cid.Cid {
	return b.cid
}

// String is a helper for output
func (b *EthBlock) String() string {
	return fmt.Sprintf("<EthBlock %s>", b.Cid())
}

// Loggable returns a map the type of IPLD Link.
func (b *EthBlock) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "eth-block",
	}
}

/*
  Node INTERFACE
*/

// Resolve resolves a path through this node, stopping at any link boundary
// and returning the object found as well as the remaining path to traverse
func (b *EthBlock) Resolve(p []string) (interface{}, []string, error) {
	if len(p) == 0 {
		return b, nil, nil
	}

	switch p[0] {
	case "tx":
		return &node.Link{Cid: castCommonHash(MEthTxTrie, b.TxHash)}, p[1:], nil
	case "parent":
		return &node.Link{Cid: castCommonHash(MEthBlock, b.ParentHash)}, p[1:], nil
	default:
		return nil, nil, fmt.Errorf("no such link")
	}
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (b *EthBlock) Tree(p string, depth int) []string {
	return nil
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

// Copy will go away. It is here to comply with the Node interface.
func (b *EthBlock) Copy() node.Node {
	panic("dont use this yet")
}

// Links is a helper function that returns all links within this object
// HINT: Use `ipfs refs <cid>`
func (b *EthBlock) Links() []*node.Link {
	return []*node.Link{
		&node.Link{Cid: castCommonHash(MEthBlock, b.ParentHash)},
		&node.Link{Cid: castCommonHash(MEthTxReceiptTrie, b.ReceiptHash)},
		&node.Link{Cid: castCommonHash(MEthStateTrie, b.Root)},
		&node.Link{Cid: castCommonHash(MEthTxTrie, b.TxHash)},
		&node.Link{Cid: castCommonHash(MEthBlockList, b.UncleHash)},
	}
}

// Stat will go away. It is here to comply with the Node interface.
func (b *EthBlock) Stat() (*node.NodeStat, error) {
	return &node.NodeStat{}, nil
}

// Size will go away. It is here to comply with the Node interface.
func (b *EthBlock) Size() (uint64, error) {
	return 0, nil
}

/*
  EthBlock functions
*/

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
		"parent":     castCommonHash(MEthBlock, b.ParentHash),
		"receipts":   castCommonHash(MEthTxReceiptTrie, b.ReceiptHash),
		"root":       castCommonHash(MEthStateTrie, b.Root),
		"tx":         castCommonHash(MEthTxTrie, b.TxHash),
		"uncles":     castCommonHash(MEthBlockList, b.UncleHash),
	}
	return json.Marshal(out)
}

// Defines the output of the JSON RPC API for either
// "eth_BlockByHash" or "eth_BlockByHeader".
type objJSONBlock struct {
	Result objJSONBlockResult `json:"result"`
}

// Nested struct that takes the contents of the JSON field "result".
type objJSONBlockResult struct {
	types.Header           // Use its fields and unmarshaler
	*objJSONBlockResultExt // Add these fields to the parsing
}

// Facilitates the composition of the field "result", adding to the
// Header fields, both uncles and transactions.
type objJSONBlockResultExt struct {
	UncleHashes  []common.Hash        `json:"uncles"`
	Transactions []*types.Transaction `json:"transactions"`
}

// Overrides the function types.Header.UnmarshalJSON, allowing us
// to parse the fields of Header, plus uncle hashes and transactions.
// (yes, uncle hashes. You will need to "eth_getUncleCountByBlockHash"
// per uncle... Don't kill the messenger)
func (o *objJSONBlockResult) UnmarshalJSON(input []byte) error {
	err := o.Header.UnmarshalJSON(input)
	if err != nil {
		return err
	}

	o.objJSONBlockResultExt = &objJSONBlockResultExt{}
	err = json.Unmarshal(input, o.objJSONBlockResultExt)
	return err
}
