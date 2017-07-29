package main

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
	return nil
}

// EthBlockRawInputParser will take the piped input, which could an RLP binary
// of either an RLP block header, or an RLP body (header + uncles + txs)
// to return an IPLD Node slice.
func EthBlockRawInputParser(r io.Reader) ([]node.Node, error) {
	blockHeader, err := eth.FromBlockRLP(r)
	if err != nil {
		return nil, err
	}
	var out []node.Node
	out = append(out, blockHeader)
	return out, nil
}

// EthBlockJSONInputParser will take the piped input, a JSON representation of
// a block header or body (header + uncles + txs), to return an IPLD Node slice.
func EthBlockJSONInputParser(r io.Reader) ([]node.Node, error) {
	blockHeader, err := eth.FromBlockJSON(r)
	if err != nil {
		return nil, err
	}
	var out []node.Node
	out = append(out, blockHeader)
	return out, nil
}

/*
  OUTPUT BLOCK DECODERS
*/

// RegisterBlockDecoders enters which functions will help us to decode the requested IPLD blocks.
func (ep *EthereumPlugin) RegisterBlockDecoders(dec node.BlockDecoder) error {
	dec.Register(eth.MEthBlock, EthBlockParser) // eth-block
	// TODO
	// Let's deal with these two elements later
	// dec.Register(eth.MEthTx, EthTxParser)
	// dec.Register(eth.MEthTxTrie, EthTxTrieParser)
	return nil
}

// EthBlockParser takes care of the eth-block IPLD objects (ethereum block headers)
func EthBlockParser(b block.Block) (node.Node, error) {
	return eth.DecodeBlockHeader(b.Cid(), b.RawData())
}

// EthTxParser takes care of the eth-tx IPLD objects (ethereum transactions)
func EthTxParser(b block.Block) (node.Node, error) {
	return eth.ParseTx(b.RawData())
}

// EthTxTrieParser takes care of the eth-tx-trie objects (ethereum transaction trie nodes)
func EthTxTrieParser(b block.Block) (node.Node, error) {
	return eth.NewTrieNode(b.RawData())
}
