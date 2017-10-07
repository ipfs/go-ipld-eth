package ipldeth

import (
	"fmt"
	"io"
	"io/ioutil"

	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"
)

// EthStorageTrie (eth-storage-trie, codec 0x98), represents
// a node from the storage trie in ethereum.
type EthStorageTrie struct {
	*TrieNode
}

// Static (compile time) check that EthStorageTrie satisfies the node.Node interface.
var _ node.Node = (*EthStorageTrie)(nil)

/*
  INPUT
*/

// FromStorageTrieRLP takes the RLP representation of an ethereum
// storage trie node to return it as an IPLD node for further processing.
func FromStorageTrieRLP(r io.Reader) (*EthStorageTrie, error) {
	rawdata, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	c := rawdataToCid(MEthStorageTrie, rawdata)

	// Let's run the whole mile and process the nodeKind and
	// its elements, in case somebody would need this function
	// to parse an RLP element from the filesystem
	return DecodeEthStorageTrie(c, rawdata)
}

/*
  OUTPUT
*/

// DecodeEthStorageTrie returns an EthStorageTrie object from its cid and rawdata.
func DecodeEthStorageTrie(c *cid.Cid, b []byte) (*EthStorageTrie, error) {
	tn, err := decodeTrieNode(c, b, decodeEthStorageTrieLeaf)
	if err != nil {
		return nil, err
	}
	return &EthStorageTrie{TrieNode: tn}, nil
}

// decodeEthStorageTrieLeaf parses a eth-tx-trie leaf
// from decoded RLP elements
func decodeEthStorageTrieLeaf(i []interface{}) ([]interface{}, error) {
	return []interface{}{
		i[0].([]byte),
		i[1].([]byte),
	}, nil
}

/*
  Block INTERFACE
*/

// RawData returns the binary of the RLP encode of the storage trie node.
func (st *EthStorageTrie) RawData() []byte {
	return st.rawdata
}

// Cid returns the cid of the storage trie node.
func (st *EthStorageTrie) Cid() *cid.Cid {
	return st.cid
}

// String is a helper for output
func (st *EthStorageTrie) String() string {
	return fmt.Sprintf("<EthereumStorageTrie %s>", st.cid)
}

// Loggable returns in a map the type of IPLD Link.
func (st *EthStorageTrie) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"type": "eth-storage-trie",
	}
}
