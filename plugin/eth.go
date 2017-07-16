package main

import (
	"bytes"
	"io"

	block "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-ipfs/core/coredag"
	plugin "github.com/ipfs/go-ipfs/plugin"
	eth "github.com/ipfs/go-ipld-eth"

	node "github.com/ipfs/go-ipld-format"
)

type EthereumPlugin struct{}

var _ plugin.PluginIPLD = (*EthereumPlugin)(nil)

func (ep *EthereumPlugin) Init() error {
	return nil
}

func (ep *EthereumPlugin) Name() string {
	return "ethereum"
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
