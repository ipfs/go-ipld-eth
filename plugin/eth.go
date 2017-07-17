package main

import (
	"bytes"
	"io"

	block "gx/ipfs/QmVA4mafxbfH5aEvNz8fyoxC6J1xhAtw88B4GerPznSZBg/go-block-format"
	node "gx/ipfs/QmYNyRZJBUYPNrLszFmrBrPJbsBh2vMsefz5gnDpB5M1P6/go-ipld-format"

	"github.com/ipfs/go-ipfs/core/coredag"
	plugin "github.com/ipfs/go-ipfs/plugin"
	eth "github.com/ipfs/go-ipld-eth"
)

var Plugins = []plugin.Plugin{
	&EthereumPlugin{},
}

type EthereumPlugin struct{}

var _ plugin.PluginIPLD = (*EthereumPlugin)(nil)

func (ep *EthereumPlugin) Init() error {
	return nil
}

func (ep *EthereumPlugin) Name() string {
	return "ipld-ethereum"
}

func (ep *EthereumPlugin) Version() string {
	return "0.0.1"
}

func (ep *EthereumPlugin) RegisterBlockDecoders(dec node.BlockDecoder) error {
	dec.Register(eth.MEthBlock, EthBlockParser)
	dec.Register(eth.MEthTx, EthTxParser)
	dec.Register(eth.MEthTxTrie, EthTxTrieParser)
	return nil
}

func (ep *EthereumPlugin) RegisterInputEncParsers(iec coredag.InputEncParsers) error {
	iec.AddParser("raw", "eth", BlockInputParser)
	return nil
}

func EthBlockParser(b block.Block) (node.Node, error) {
	return eth.DecodeBlock(bytes.NewReader(b.RawData()))
}

func EthTxParser(b block.Block) (node.Node, error) {
	return eth.ParseTx(b.RawData())
}

func EthTxTrieParser(b block.Block) (node.Node, error) {
	return eth.NewTrieNode(b.RawData())
}

func BlockInputParser(r io.Reader) ([]node.Node, error) {
	blk, txs, tries, uncles, err := eth.FromRlpBlockMessage(r)
	if err != nil {
		return nil, err
	}

	var out []node.Node
	out = append(out, blk)
	for _, tx := range txs {
		out = append(out, tx)
	}
	for _, t := range tries {
		out = append(out, t)
	}
	for _, unc := range uncles {
		out = append(out, unc)
	}
	return out, nil
}
