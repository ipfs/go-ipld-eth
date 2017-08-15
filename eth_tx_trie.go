package ipldeth

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
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
	var elements []interface{}

	nodeKind, i, err := decodeTrieNode(b)
	if err != nil {
		return nil, err
	}

	switch nodeKind {
	case "extension":
		elements, err = parseEthTxTrieExtension(i)
	case "leaf":
		elements, err = parseEthTxTrieLeaf(i)
	case "branch":
		elements, err = parseEthTxTrieBranch(i)
	default:
		return nil, fmt.Errorf("nodeKind not supported")
	}

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

// parseEthTxTrieExtension helper improves readability
func parseEthTxTrieExtension(i []interface{}) ([]interface{}, error) {
	return []interface{}{
		i[0].([]byte),
		keccak256ToCid(MEthTxTrie, i[1].([]byte)),
	}, nil
}

// parseEthTxTrieLeaf helper improves readability
func parseEthTxTrieLeaf(i []interface{}) ([]interface{}, error) {
	var t types.Transaction
	err := rlp.DecodeBytes(i[1].([]byte), &t)
	if err != nil {
		return nil, err
	}
	return []interface{}{
		i[0].([]byte),
		&EthTx{
			Transaction: &t,
			cid:         rawdataToCid(MEthTx, i[1].([]byte)),
			rawdata:     i[1].([]byte),
		},
	}, nil
}

// parseEthTxTrieBranch helper improves readability
func parseEthTxTrieBranch(i []interface{}) ([]interface{}, error) {
	var out []interface{}

	for _, vi := range i {
		v := vi.([]byte)

		switch len(v) {
		case 0:
			out = append(out, nil)
		case 32:
			out = append(out, keccak256ToCid(MEthTxTrie, v))
		default:
			return nil, fmt.Errorf("unrecognized object: %v", v)
		}
	}

	return out, nil
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
	obj, rest, err := resolveTriePath(p, t.nodeKind, t.elements)
	if err != nil {
		return nil, nil, err
	}

	return obj, rest, nil
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (t *EthTxTrie) Tree(p string, depth int) []string {
	if p != "" || depth == 0 {
		return nil
	}

	var out []string

	switch t.nodeKind {
	case "branch":
		for i, elem := range t.elements {
			if _, ok := elem.(*cid.Cid); ok {
				out = append(out, fmt.Sprintf("%x", i))
			}
		}
		return out

	default:
		return nil
	}
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
	var out []*node.Link

	for _, i := range t.elements {
		c, ok := i.(*cid.Cid)
		if ok {
			out = append(out, &node.Link{Cid: c})
		}
	}

	return out
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
	var out map[string]interface{}

	switch t.nodeKind {
	case "extension":
		var val string
		for _, e := range t.elements[0].([]byte) {
			val += fmt.Sprintf("%x", e)
		}
		out = map[string]interface{}{
			"type": "extension",
			val:    t.elements[1],
		}
	case "branch":
		out = map[string]interface{}{
			"type": "branch",
			"0":    t.elements[0],
			"1":    t.elements[1],
			"2":    t.elements[2],
			"3":    t.elements[3],
			"4":    t.elements[4],
			"5":    t.elements[5],
			"6":    t.elements[6],
			"7":    t.elements[7],
			"8":    t.elements[8],
			"9":    t.elements[9],
			"a":    t.elements[10],
			"b":    t.elements[11],
			"c":    t.elements[12],
			"d":    t.elements[13],
			"e":    t.elements[14],
			"f":    t.elements[15],
		}
	default:
		return nil, fmt.Errorf("nodeKind %s not supported", t.nodeKind)
	}

	return json.Marshal(out)
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
