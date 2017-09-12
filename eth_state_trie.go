package ipldeth

import (
	"fmt"
	"io"
	"io/ioutil"

	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"

	"github.com/ethereum/go-ethereum/rlp"
)

// EthStateTrie (eth-state-trie, codec 0x96), represents
// a node from the satte trie in ethereum.
type EthStateTrie struct {
	*TrieNode
}

// Static (compile time) check that EthStateTrie satisfies the node.Node interface.
var _ node.Node = (*EthStateTrie)(nil)

/*
  INPUT
*/

// FromStateTrieRLP takes the RLP representation of an ethereum
// state trie node to return it as an IPLD node for further processing.
func FromStateTrieRLP(r io.Reader) (*EthStateTrie, error) {
	rawdata, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	c := rawdataToCid(MEthStateTrie, rawdata)

	// Let's run the whole mile and process the nodeKind and
	// its elements, in case somebody would need this function
	// to parse an RLP element from the filesystem
	return DecodeEthStateTrie(c, rawdata)
}

/*
  OUTPUT
*/

// DecodeEthStateTrie returns an EthStateTrie object from its cid and rawdata.
func DecodeEthStateTrie(c *cid.Cid, b []byte) (*EthStateTrie, error) {
	tn, err := decodeTrieNode(c, b, decodeEthStateTrieLeaf)
	if err != nil {
		return nil, err
	}
	return &EthStateTrie{TrieNode: tn}, nil
}

// decodeEthStateTrieLeaf parses a eth-tx-trie leaf
// from decoded RLP elements
func decodeEthStateTrieLeaf(i []interface{}) ([]interface{}, error) {
	var account EthAccount
	err := rlp.DecodeBytes(i[1].([]byte), &account)
	if err != nil {
		return nil, err
	}
	return []interface{}{
		i[0].([]byte),
		&EthAccountSnapshot{
			EthAccount: &account,
			cid:        rawdataToCid(MEthAccountSnapshot, i[1].([]byte)),
			rawdata:    i[1].([]byte),
		},
	}, nil
}

/*
  Block INTERFACE
*/

// RawData returns the binary of the RLP encode of the state trie node.
func (st *EthStateTrie) RawData() []byte {
	return st.rawdata
}

// Cid returns the cid of the state trie node.
func (st *EthStateTrie) Cid() *cid.Cid {
	return st.cid
}

// String is a helper for output
func (st *EthStateTrie) String() string {
	return fmt.Sprintf("<EthereumStateTrie %s>", st.cid)
}

// Loggable returns in a map the type of IPLD Link.
func (st *EthStateTrie) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "eth-state-trie",
	}
}
