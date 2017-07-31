package ipldeth

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	block "github.com/ipfs/go-block-format"
)

func TestBlockBodyRlpParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-block-body-rlp-999999")
	if err != nil {
		t.Fatal(err)
	}

	output, _, _, err := FromBlockRLP(fi)
	if err != nil {
		t.Fatal(err)
	}

	testEthBlockHeaderFields(output, t)
}

func TestBlockHeaderRlpParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-block-header-rlp-999999")
	if err != nil {
		t.Fatal(err)
	}

	output, _, _, err := FromBlockRLP(fi)
	if err != nil {
		t.Fatal(err)
	}

	testEthBlockHeaderFields(output, t)
}

func TestBlockBodyJsonParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-block-body-json-999999")
	if err != nil {
		t.Fatal(err)
	}

	output, _, _, err := FromBlockJSON(fi)
	if err != nil {
		t.Fatal(err)
	}

	testEthBlockHeaderFields(output, t)
}

// TestDecodeBlockHeader should work for both inputs (block header and block body)
// as what we are storing is just the block header
func TestDecodeBlockHeader(t *testing.T) {
	// Prepare the "fetched block". This one is supposed to be in the datastore
	// and given away by github.com/ipfs/go-ipfs/merkledag
	fi, err := os.Open("test_data/eth-block-header-rlp-999999")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(fi)
	if err != nil {
		t.Fatal(err)
	}

	c := rawdataToCid(MEthBlock, b)

	storedBlockHeader, err := block.NewBlockWithCid(b, c)
	if err != nil {
		t.Fatal(err)
	}

	// Now the proper test
	ethBlock, err := DecodeEthBlock(storedBlockHeader.Cid(), storedBlockHeader.RawData())
	if err != nil {
		t.Fatal(err)
	}

	testEthBlockHeaderFields(ethBlock, t)
}

/*
  AUXILIARS
*/

func testEthBlockHeaderFields(ethBlock *EthBlock, t *testing.T) {
	// Was the cid calculated?
	if ethBlock.Cid().String() != "z43AaGF4uHSY4waU68L3DLUKHZP7yfZoo6QbLmid5HomZ4WtbWw" {
		t.Fatal("Wrong cid")
	}

	// Do we have the rawdata available?
	if fmt.Sprintf("%x", ethBlock.RawData()[:10]) != "f90218a0d33c9dde9fff" {
		t.Fatal("Wrong Rawdata")
	}

	// Proper Fields of types.Header
	if fmt.Sprintf("%x", ethBlock.ParentHash) != "d33c9dde9fff0ebaa6e71e8b26d2bda15ccf111c7af1b633698ac847667f0fb4" {
		t.Fatal("Wrong ParentHash")
	}
	if fmt.Sprintf("%x", ethBlock.UncleHash) != "1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347" {
		t.Fatal("Wrong UncleHash field")
	}
	if fmt.Sprintf("%x", ethBlock.Coinbase) != "52bc44d5378309ee2abf1539bf71de1b7d7be3b5" {
		t.Fatal("Wrong Coinbase")
	}
	if fmt.Sprintf("%x", ethBlock.Root) != "ed98aa4b5b19c82fb35364f08508ae0a6dec665fa57663dca94c5d70554cde10" {
		t.Fatal("Wrong Root")
	}
	if fmt.Sprintf("%x", ethBlock.TxHash) != "447cbd8c48f498a6912b10831cdff59c7fbfcbbe735ca92883d4fa06dcd7ae54" {
		t.Fatal("Wrong TxHash")
	}
	if fmt.Sprintf("%x", ethBlock.ReceiptHash) != "7fa0f6ca2a01823208d80801edad37e3e3a003b55c89319b45eb1f97862ad229" {
		t.Fatal("Wrong ReceiptHash")
	}
	if len(ethBlock.Bloom) != 256 {
		t.Fatal("Wrong Bloom Length")
	}
	if fmt.Sprintf("%x", ethBlock.Bloom[71:76]) != "0000000000" { // You wouldn't want me to print out the whole bloom field?
		t.Fatal("Wrong Bloom")
	}
	if ethBlock.Difficulty.String() != "12555463106190" {
		t.Fatal("Wrong Difficulty")
	}
	if ethBlock.Number.String() != "999999" {
		t.Fatal("Wrong Block Number")
	}
	if ethBlock.GasLimit.String() != "3141592" {
		t.Fatal("Wrong Gas Limit")
	}
	if ethBlock.GasUsed.String() != "231000" {
		t.Fatal("Wrong Gas Used")
	}
	if ethBlock.Time.String() != "1455404037" {
		t.Fatal("Wrong Time")
	}
	if fmt.Sprintf("%x", ethBlock.Extra) != "d783010303844765746887676f312e342e32856c696e7578" {
		t.Fatal("Wrong Extra")
	}
	if fmt.Sprintf("%x", ethBlock.Nonce) != "f491f46b60fe04b3" {
		t.Fatal("Wrong Nonce")
	}
	if fmt.Sprintf("%x", ethBlock.MixDigest) != "5b10f4a08a6c209d426f6158bd24b574f4f7b7aa0099c67c14a1f693b4dd04d0" {
		t.Fatal("Wrong MixDigest")
	}
}
