package ipldeth

import (
	"fmt"

	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"
)

// EthTxTrie (eth-tx-trie codec 0x95) represents
// a node from the transaction trie in ethereum.
type EthTxTrie struct {
	// leaf, extension or branch
	nodeKind string

	// If leaf or extension: [0] is key, [1] is val.
	// If branch: [0] - [16] are children.
	elements []interface{}

	// IPLD block information
	cid     *cid.Cid
	rawdata []byte
}

// Static (compile time) check that EthTxTrie satisfies the node.Node interface.
var _ node.Node = (*EthTxTrie)(nil)

/*
 INPUT
*/

// To create a proper trie of the eth-tx-trie objects, it is required
// to input all transactions belonging to a forest in a single step.
// We are adding the transactions, and creating its trie on
// block body parsing time.

/*
  OUTPUT
*/

// DecodeEthTxTrie returns an EthTxTrie object from its cid and rawdata.
func DecodeEthTxTrie(c *cid.Cid, b []byte) (*EthTxTrie, error) {
	nodeKind, elements, err := decodeTrieNode(b)
	if err != nil {
		return nil, err
	}

	return &EthTxTrie{
		nodeKind: nodeKind,
		elements: elements,
		rawdata:  b,
		cid:      c,
	}, nil
}

/*
  Block INTERFACE
*/

// RawData returns the binary of the RLP encode of the transaction.
func (t *EthTxTrie) RawData() []byte {
	return t.rawdata
}

// Cid returns the cid of the transaction.
func (t *EthTxTrie) Cid() *cid.Cid {
	return t.cid
}

// String is a helper for output
func (t *EthTxTrie) String() string {
	return fmt.Sprintf("<EthereumTxTrie %s>", t.cid)
}

// Loggable returns in a map the type of IPLD Link.
func (t *EthTxTrie) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "eth-tx-trie",
	}
}

// Resolve resolves a path through this node, stopping at any link boundary
// and returning the object found as well as the remaining path to traverse
func (t *EthTxTrie) Resolve(p []string) (interface{}, []string, error) {
	p, err := validateTriePath(p, getTxFields())
	if err != nil {
		return nil, nil, err
	}

	switch t.nodeKind {
	case "branch":
		child := t.elements[getHexIndex(p[0])]
		if child != nil {
			return &node.Link{Cid: child.(*cid.Cid)}, p[1:], nil
		}
		return nil, nil, fmt.Errorf("no such link")
	default:
		return nil, nil, fmt.Errorf("nodeKind case not implemented")
	}
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (t *EthTxTrie) Tree(p string, depth int) []string {
	// PLACEHOLDER
	return nil
	// PLACEHOLDER
}

// ResolveLink is a helper function that calls resolve and asserts the
// output is a link
func (t *EthTxTrie) ResolveLink(p []string) (*node.Link, []string, error) {
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
func (t *EthTxTrie) Copy() node.Node {
	panic("dont use this yet")
}

// Links is a helper function that returns all links within this object
func (t *EthTxTrie) Links() []*node.Link {
	// PLACEHOLDER
	return nil
	// PLACEHOLDER
}

// Stat will go away. It is here to comply with the interface.
func (t *EthTxTrie) Stat() (*node.NodeStat, error) {
	return &node.NodeStat{}, nil
}

// Size will go away. It is here to comply with the interface.
func (t *EthTxTrie) Size() (uint64, error) {
	return 0, nil
}

/*
  EthTxTrie functions
*/

// MarshalJSON processes the transaction trie into readable JSON format.
func (t *EthTxTrie) MarshalJSON() ([]byte, error) {
	// PLACEHOLDER
	return nil, fmt.Errorf("Not implemented")
	// PLACEHOLDER
}

// txTrie wraps a localTrie for use on the transaction trie.
type txTrie struct {
	*localTrie
}

// newTxTrie initializes and returns a txTrie.
func newTxTrie() *txTrie {
	return &txTrie{
		localTrie: newLocalTrie(),
	}
}

// addTx takes the rawdata of an EthTx to incorporate it to the
// transaction trie.
func (tt *txTrie) addTx(idx int, rawdata []byte) {
	tt.add(idx, rawdata)
}

// rootHash returns the computed trie root.
// Useful for sanity checks on parsed data.
func (tt *txTrie) rootHash() []byte {
	return tt.trie.Hash().Bytes()
}

// getNodes invokes the localTrie, which computes the root hash of the
// transaction trie and returns its database keys, to return a slice
// of EthTxTrie nodes.
func (tt *txTrie) getNodes() []*EthTxTrie {
	keys := tt.getKeys()
	var out []*EthTxTrie

	for _, k := range keys {
		rawdata, err := tt.db.Get(k)
		if err != nil {
			panic(err)
		}

		out = append(out, &EthTxTrie{
			cid:     rawdataToCid(MEthTxTrie, rawdata),
			rawdata: rawdata,
		})
	}

	return out
}

// getTxFields returns the fields defined in an ethereum transaction
func getTxFields() map[string]interface{} {
	return map[string]interface{}{
		"nonce":     nil,
		"gasPrice":  nil,
		"gas":       nil,
		"toAddress": nil,
		"value":     nil,
		"data":      nil,
	}
}
