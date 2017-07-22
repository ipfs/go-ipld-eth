package main

import (
	"bytes"
	"io"

	block "github.com/ipfs/go-block-format"
	node "github.com/ipfs/go-ipld-format"

	"github.com/ipfs/go-ipfs/core/coredag"
	plugin "github.com/ipfs/go-ipfs/plugin"
	eth "github.com/ipfs/go-ipld-eth"
)

// Plugins declare what and how many of these will be defined.
var Plugins = []plugin.Plugin{
	&EthereumPlugin{},
}

// EthereumPlugin is the main structure.
type EthereumPlugin struct{}

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
  INPUT
*/

// RegisterInputEncParsers enters the encode parsers needed to put the blocks into IPLD.
func (ep *EthereumPlugin) RegisterInputEncParsers(iec coredag.InputEncParsers) error {
	iec.AddParser("raw", "eth-block", EthBlockInputParser)
	return nil
}

// EthBlockInputParser will take the piped input,
// an RLP binary of a block header, to return an IPLD node slice.
func EthBlockInputParser(r io.Reader) ([]node.Node, error) {
	blockHeader, err := eth.FromBlockHeaderRLP(r)
	if err != nil {
		return nil, err
	}
	var out []node.Node
	out = append(out, blockHeader)
	return out, nil
}

/*
  OUTPUT
*/

// RegisterBlockDecoders enters which functions will help us to decode the requested IPLD blocks.
func (ep *EthereumPlugin) RegisterBlockDecoders(dec node.BlockDecoder) error {
	dec.Register(eth.MEthBlock, EthBlockParser) // eth-block
	dec.Register(eth.MEthTx, EthTxParser)
	dec.Register(eth.MEthTxTrie, EthTxTrieParser)
	return nil
}

// EthBlockParser takes care of the eth-block IPLD objects (ethereum block headers)
func EthBlockParser(b block.Block) (node.Node, error) {
	return eth.DecodeBlock(bytes.NewReader(b.RawData()))
}

// EthTxParser takes care of the eth-tx IPLD objects (ethereum transactions)
func EthTxParser(b block.Block) (node.Node, error) {
	return eth.ParseTx(b.RawData())
}

// EthTxTrieParser takes care of the eth-tx-trie objects (ethereum transaction trie nodes)
func EthTxTrieParser(b block.Block) (node.Node, error) {
	return eth.NewTrieNode(b.RawData())
}
