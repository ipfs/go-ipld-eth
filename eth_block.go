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
// an ethereum block header or body (header, ommers and txs)
// to return it as a set of IPLD nodes for further processing.
func FromBlockRLP(r io.Reader) (*EthBlock, []*EthTx, []*EthTxTrie, error) {
	// We may want to use this stream several times
	rawdata, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, nil, nil, err
	}

	// Let's try to decode the received element as a block body
	var decodedBlock types.Block
	err = rlp.Decode(bytes.NewBuffer(rawdata), &decodedBlock)
	if err != nil {
		if err.Error()[:41] != "rlp: expected input list for types.Header" {
			return nil, nil, nil, err
		}

		// Maybe it is just a header... (body sans ommers and txs)
		var decodedHeader types.Header
		err := rlp.Decode(bytes.NewBuffer(rawdata), &decodedHeader)
		if err != nil {
			return nil, nil, nil, err
		}

		// It was a header
		return &EthBlock{
			Header:  &decodedHeader,
			cid:     rawdataToCid(MEthBlock, rawdata),
			rawdata: rawdata,
		}, nil, nil, nil
	}

	// This is a block body (header + ommers + txs)
	// We'll extract the header bits here
	headerRawData := getRLP(decodedBlock.Header())
	ethBlock := &EthBlock{
		Header:  decodedBlock.Header(),
		cid:     rawdataToCid(MEthBlock, headerRawData),
		rawdata: headerRawData,
	}

	// Process the found eth-tx objects
	ethTxNodes, ethTxTrieNodes, err := processTransactions(decodedBlock.Transactions(),
		decodedBlock.Header().TxHash[:])
	if err != nil {
		return nil, nil, nil, err
	}

	return ethBlock, ethTxNodes, ethTxTrieNodes, nil
}

// FromBlockJSON takes the output of an ethereum client JSON API
// (i.e. parity or geth) and returns a set of IPLD nodes.
func FromBlockJSON(r io.Reader) (*EthBlock, []*EthTx, []*EthTxTrie, error) {
	var obj objJSONBlock
	dec := json.NewDecoder(r)
	err := dec.Decode(&obj)
	if err != nil {
		return nil, nil, nil, err
	}

	headerRawData := getRLP(obj.Result.Header)
	ethBlock := &EthBlock{
		Header:  &obj.Result.Header,
		cid:     rawdataToCid(MEthBlock, headerRawData),
		rawdata: headerRawData,
	}

	// Process the found eth-tx objects
	ethTxNodes, ethTxTrieNodes, err := processTransactions(obj.Result.Transactions,
		obj.Result.Header.TxHash[:])
	if err != nil {
		return nil, nil, nil, err
	}

	return ethBlock, ethTxNodes, ethTxTrieNodes, nil
}

// processTransactions will take the found transactions in a parsed block body
// to return IPLD node slices for eth-tx and eth-tx-trie
func processTransactions(txs []*types.Transaction, expectedTxRoot []byte) ([]*EthTx, []*EthTxTrie, error) {
	var ethTxNodes []*EthTx
	transactionTrie := newTxTrie()

	for idx, tx := range txs {
		ethTx := NewTx(tx) // Will panic if it finds an error while parsing a tx
		ethTxNodes = append(ethTxNodes, ethTx)
		transactionTrie.add(idx, ethTx.RawData())
	}

	if !bytes.Equal(transactionTrie.rootHash(), expectedTxRoot) {
		return nil, nil, fmt.Errorf("wrong transaction hash computed")
	}

	ethTxTrieNodes := transactionTrie.getNodes()

	return ethTxNodes, ethTxTrieNodes, nil
}

/*
  OUTPUT
*/

// DecodeEthBlock takes a cid and its raw binary data
// from IPFS and returns an EthBlock object for further processing.
func DecodeEthBlock(c *cid.Cid, b []byte) (*EthBlock, error) {
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
	return fmt.Sprintf("<EthBlock %s>", b.cid)
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

	first, rest := p[0], p[1:]

	switch first {
	case "parent":
		return &node.Link{Cid: commonHashToCid(MEthBlock, b.ParentHash)}, rest, nil
	case "receipts":
		return &node.Link{Cid: commonHashToCid(MEthTxReceiptTrie, b.ReceiptHash)}, rest, nil
	case "root":
		return &node.Link{Cid: commonHashToCid(MEthStateTrie, b.Root)}, rest, nil
	case "tx":
		return &node.Link{Cid: commonHashToCid(MEthTxTrie, b.TxHash)}, rest, nil
	case "uncles":
		return &node.Link{Cid: commonHashToCid(MEthBlockList, b.UncleHash)}, rest, nil
	}

	if len(p) != 1 {
		return nil, nil, fmt.Errorf("unexpected path elements past %s", first)
	}

	switch first {
	case "bloom":
		return b.Bloom, nil, nil
	case "coinbase":
		return b.Coinbase, nil, nil
	case "difficulty":
		return b.Difficulty, nil, nil
	case "extra":
		// This is a []byte. By default they are marshalled into Base64.
		return fmt.Sprintf("0x%x", b.Extra), nil, nil
	case "gaslimit":
		return b.GasLimit, nil, nil
	case "gasused":
		return b.GasUsed, nil, nil
	case "mixdigest":
		return b.MixDigest, nil, nil
	case "nonce":
		return b.Nonce, nil, nil
	case "number":
		return b.Number, nil, nil
	case "time":
		return b.Time, nil, nil
	default:
		return nil, nil, fmt.Errorf("no such link")
	}
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (b *EthBlock) Tree(p string, depth int) []string {
	if p != "" || depth == 0 {
		return nil
	}

	return []string{
		"time",
		"bloom",
		"coinbase",
		"difficulty",
		"extra",
		"gaslimit",
		"gasused",
		"mixdigest",
		"nonce",
		"number",
		"parent",
		"receipts",
		"root",
		"tx",
		"uncles",
	}
}

// ResolveLink is a helper function that allows easier traversal of links through blocks
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
		&node.Link{Cid: commonHashToCid(MEthBlock, b.ParentHash)},
		&node.Link{Cid: commonHashToCid(MEthTxReceiptTrie, b.ReceiptHash)},
		&node.Link{Cid: commonHashToCid(MEthStateTrie, b.Root)},
		&node.Link{Cid: commonHashToCid(MEthTxTrie, b.TxHash)},
		&node.Link{Cid: commonHashToCid(MEthBlockList, b.UncleHash)},
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

// MarshalJSON processes the block header into readable JSON format,
// converting the right links into their cids, and keeping the original
// hex hash, allowing the user to simplify external queries.
func (b *EthBlock) MarshalJSON() ([]byte, error) {
	out := map[string]interface{}{
		"time":       b.Time,
		"bloom":      b.Bloom,
		"coinbase":   b.Coinbase,
		"difficulty": b.Difficulty,
		"extra":      fmt.Sprintf("0x%x", b.Extra),
		"gaslimit":   b.GasLimit,
		"gasused":    b.GasUsed,
		"mixdigest":  b.MixDigest,
		"nonce":      b.Nonce,
		"number":     b.Number,
		"parent":     commonHashToCid(MEthBlock, b.ParentHash),
		"receipts":   commonHashToCid(MEthTxReceiptTrie, b.ReceiptHash),
		"root":       commonHashToCid(MEthStateTrie, b.Root),
		"tx":         commonHashToCid(MEthTxTrie, b.TxHash),
		"uncles":     commonHashToCid(MEthBlockList, b.UncleHash),
	}
	return json.Marshal(out)
}

// objJSONBlock defines the output of the JSON RPC API for either
// "eth_BlockByHash" or "eth_BlockByHeader".
type objJSONBlock struct {
	Result objJSONBlockResult `json:"result"`
}

// objJSONBLockResult is the  nested struct that takes
// the contents of the JSON field "result".
type objJSONBlockResult struct {
	types.Header           // Use its fields and unmarshaler
	*objJSONBlockResultExt // Add these fields to the parsing
}

// objJSONBLockResultExt facilitates the composition
// of the field "result", adding to the
// `types.Header` fields, both ommers (their hashes) and transactions.
type objJSONBlockResultExt struct {
	OmmerHashes  []common.Hash        `json:"uncles"`
	Transactions []*types.Transaction `json:"transactions"`
}

// UnmarshalJSON overrides the function types.Header.UnmarshalJSON, allowing us
// to parse the fields of Header, plus ommer hashes and transactions.
// (yes, ommer hashes. You will need to "eth_getUncleCountByBlockHash" per each ommer)
func (o *objJSONBlockResult) UnmarshalJSON(input []byte) error {
	err := o.Header.UnmarshalJSON(input)
	if err != nil {
		return err
	}

	o.objJSONBlockResultExt = &objJSONBlockResultExt{}
	err = json.Unmarshal(input, o.objJSONBlockResultExt)
	return err
}
