package ipldeth

import (
	"encoding/json"
	"fmt"
	"strconv"

	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// EthTxTrie (eth-tx-trie codec 0x95) represents
// an ethereum transaction as a leaf of a merkle tree.
// i.e. its fields can be accessed as a Transaction, but its rawdata is
// the RLP-encoding as a merkle leaf.
type EthTxTrie struct {
	*types.Transaction

	arr     []interface{}
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

// DecodeTxTrie returns an EthTxTrie object from its cid and rawdata.
// It is used by both the trie calculations on block body parsing,
// as well as by the decoding of an object from the IPLD forest.
func DecodeEthTxTrie(c *cid.Cid, b []byte) (*EthTxTrie, error) {
	var i []interface{}
	err := rlp.DecodeBytes(b, &i)
	if err != nil {
		return nil, err
	}

	switch len(i) {
	case 2:
		key := i[0].([]byte)
		valb := i[1].([]byte)

		var val interface{}
		if len(valb) == 32 {
			// This is a reference to another eth-tx-trie object
			val = keccak256ToCid(MEthTxTrie, valb)
		} else {
			// This is a proper transaction
			var t types.Transaction
			err := rlp.DecodeBytes(valb, &t)
			if err != nil {
				return nil, err
			}
			val = &EthTx{
				Transaction: &t,
				cid:         rawdataToCid(MEthTx, valb),
				rawdata:     valb,
			}
		}
		return &EthTxTrie{
			arr:     []interface{}{key, val},
			cid:     c,
			rawdata: b,
		}, nil
	case 17:
		var parsed []interface{}
		for _, v := range i {
			bv := v.([]byte)
			switch len(bv) {
			case 0:
				parsed = append(parsed, nil)
			case 32:
				parsed = append(parsed, keccak256ToCid(MEthTxTrie, bv))
			default:
				return nil, fmt.Errorf("unrecognized object in trie: %v", bv)
			}
		}
		return &EthTxTrie{
			arr:     parsed,
			rawdata: b,
			cid:     c,
		}, nil
	default:
		return nil, fmt.Errorf("unknown trie node type")
	}
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
	if len(p) == 0 {
		return t, nil, nil
	}

	i, err := strconv.Atoi(p[0])
	if err != nil {
		return nil, nil, fmt.Errorf("expected array index to trie: %s", err)
	}

	if i < 0 || i >= len(t.arr) {
		return nil, nil, fmt.Errorf("index in trie out of range")
	}

	switch obj := t.arr[i].(type) {
	case *cid.Cid:
		return &node.Link{Cid: obj}, p[1:], nil
	case *EthTx:
		return obj, p[1:], nil
	default:
		return nil, nil, fmt.Errorf("unexpected object type in trie")
	}
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (t *EthTxTrie) Tree(p string, depth int) []string {
	if p != "" {
		return nil
	}
	if depth > 0 {
		return nil
	}

	if len(t.arr) == 17 {
		var out []string
		for i, v := range t.arr {
			if len(v.([]byte)) == 0 {
				out = append(out, fmt.Sprintf("%x", i))
			}
		}
		return out
	}

	// TODO: not sure what to put here. Most of the 'keys' dont seem to be human readable
	return []string{"VALUE NODE"}

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
	for _, i := range t.arr {
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
	return uint64(t.Transaction.Size().Int64()), nil
}

/*
  EthTxTrie functions
*/

// MarshalJSON processes the transaction trie into readable JSON format.
func (t *EthTxTrie) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.arr)
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
		c := rawdataToCid(MEthTxTrie, rawdata)
		ethTxTrie, err := DecodeEthTxTrie(c, rawdata)
		if err != nil {
			panic(err)
		}
		out = append(out, ethTxTrie)
	}

	return out
}
