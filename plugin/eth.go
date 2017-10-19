package plugin

import (
	"io"

	block "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-ipfs/core/coredag"
	plugin "github.com/ipfs/go-ipfs/plugin"
	eth "github.com/ipfs/go-ipld-eth"
	node "github.com/ipfs/go-ipld-format"
)

// Plugins declare what and how many of these will be defined.
var Plugins = []plugin.Plugin{
	&EthereumPlugin{},
}

// EthereumPlugin is the main structure.
type EthereumPlugin struct{}

// Static (compile time) check that EthereumPlugin satisfies the plugin.PluginIPLD interface.
var _ plugin.PluginIPLD = (*EthereumPlugin)(nil)

// Init complies with the plugin.Plugin interface.
// Use node.BlockDecoder.Register() instead.
// see https://github.com/ipfs/go-ipld-eth/issues/1#issuecomment-316777885
func (ep *EthereumPlugin) Init() error {
	return nil
}

// Name returns the name of this plugin.
func (ep *EthereumPlugin) Name() string {
	return "ipld-ethereum"
}

// Version returns the version of this plugin.
func (ep *EthereumPlugin) Version() string {
	return "0.0.3"
}

/*
  INPUT PARSERS
*/

// RegisterInputEncParsers enters the encode parsers needed to put the blocks into the DAG.
func (ep *EthereumPlugin) RegisterInputEncParsers(iec coredag.InputEncParsers) error {
	iec.AddParser("raw", "eth-block", EthBlockRawInputParser)
	iec.AddParser("json", "eth-block", EthBlockJSONInputParser)
	iec.AddParser("raw", "eth-state-trie", EthStateTrieRawInputParser)
	iec.AddParser("raw", "eth-storage-trie", EthStorageTrieRawInputParser)
	return nil
}

// EthBlockRawInputParser will take the piped input, which could an RLP binary
// of either an RLP block header, or an RLP body (header + uncles + txs)
// to return an IPLD Node slice.
func EthBlockRawInputParser(r io.Reader, mhtype uint64, mhLen int) ([]node.Node, error) {
	blockHeader, txs, txTrieNodes, err := eth.FromBlockRLP(r)
	if err != nil {
		return nil, err
	}

	var out []node.Node
	out = append(out, blockHeader)
	for _, tx := range txs {
		out = append(out, tx)
	}
	for _, ttn := range txTrieNodes {
		out = append(out, ttn)
	}
	return out, nil
}

// EthBlockJSONInputParser will take the piped input, a JSON representation of
// a block header or body (header + uncles + txs), to return an IPLD Node slice.
func EthBlockJSONInputParser(r io.Reader, mhtype uint64, mhLen int) ([]node.Node, error) {
	blockHeader, txs, txTrieNodes, err := eth.FromBlockJSON(r)
	if err != nil {

		return nil, err
	}

	var out []node.Node
	out = append(out, blockHeader)
	for _, tx := range txs {
		out = append(out, tx)
	}
	for _, ttn := range txTrieNodes {
		out = append(out, ttn)
	}
	return out, nil
}

// EthStateTrieRawInputParser will take the piped input, which is an RLP binary
// representation of a state trie node, to return an IPLD Node.
func EthStateTrieRawInputParser(r io.Reader, mhtype uint64, mhLen int) ([]node.Node, error) {
	stateTrieNode, err := eth.FromStateTrieRLP(r)
	if err != nil {
		return nil, err
	}

	return []node.Node{stateTrieNode}, nil
}

// EthStorageTrieRawInputParser will take the piped input, which is an RLP binary
// representation of a storage trie node, to return an IPLD Node.
func EthStorageTrieRawInputParser(r io.Reader, mhtype uint64, mhLen int) ([]node.Node, error) {
	storageTrieNode, err := eth.FromStorageTrieRLP(r)
	if err != nil {
		return nil, err
	}

	return []node.Node{storageTrieNode}, nil
}

/*
  OUTPUT BLOCK DECODERS
*/

// RegisterBlockDecoders enters which functions will help us to decode the requested IPLD blocks.
func (ep *EthereumPlugin) RegisterBlockDecoders(dec node.BlockDecoder) error {
	dec.Register(eth.MEthBlock, EthBlockParser)             // eth-block
	dec.Register(eth.MEthTx, EthTxParser)                   // eth-tx
	dec.Register(eth.MEthTxTrie, EthTxTrieParser)           // eth-tx-trie
	dec.Register(eth.MEthStateTrie, EthStateTrieParser)     // eth-state-trie
	dec.Register(eth.MEthStorageTrie, EthStorageTrieParser) // eth-storage-trie
	return nil
}

// EthBlockParser takes care of the eth-block IPLD objects (ethereum block headers)
func EthBlockParser(b block.Block) (node.Node, error) {
	return eth.DecodeEthBlock(b.Cid(), b.RawData())
}

// EthTxParser takes care of the eth-tx IPLD objects (ethereum transactions)
func EthTxParser(b block.Block) (node.Node, error) {
	return eth.DecodeEthTx(b.Cid(), b.RawData())
}

// EthTxTrieParser takes care of the eth-tx-trie IPLD objects
// (ethereum transactions as patricia merkle tree leaves)
func EthTxTrieParser(b block.Block) (node.Node, error) {
	return eth.DecodeEthTxTrie(b.Cid(), b.RawData())
}

// EthStateTrieParser takes care of the eth-state-trie IPLD objects
// (ethereum patricia merkle tree state nodes)
func EthStateTrieParser(b block.Block) (node.Node, error) {
	return eth.DecodeEthStateTrie(b.Cid(), b.RawData())
}

// EthStorageTrieParser takes care of the eth-storage-trie IPLD objects
// (ethereum patricia merkle tree state nodes)
func EthStorageTrieParser(b block.Block) (node.Node, error) {
	return eth.DecodeEthStorageTrie(b.Cid(), b.RawData())
}
