package ipldeth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	block "github.com/ipfs/go-block-format"
	node "github.com/ipfs/go-ipld-format"
)

func TestBlockBodyRlpParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-block-body-rlp-999999")
	checkError(err, t)

	output, _, _, err := FromBlockRLP(fi)
	checkError(err, t)

	testEthBlockFields(output, t)
}

func TestBlockHeaderRlpParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-block-header-rlp-999999")
	checkError(err, t)

	output, _, _, err := FromBlockRLP(fi)
	checkError(err, t)

	testEthBlockFields(output, t)
}

func TestBlockBodyJsonParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-block-body-json-999999")
	checkError(err, t)

	output, _, _, err := FromBlockJSON(fi)
	checkError(err, t)

	testEthBlockFields(output, t)
}

// TestDecodeBlockHeader should work for both inputs (block header and block body)
// as what we are storing is just the block header
func TestDecodeBlockHeader(t *testing.T) {
	storedEthBlock := prepareStoredEthBlock("test_data/eth-block-header-rlp-999999", t)

	ethBlock, err := DecodeEthBlock(storedEthBlock.Cid(), storedEthBlock.RawData())
	checkError(err, t)

	testEthBlockFields(ethBlock, t)
}

func TestEthBlockJSONMarshal(t *testing.T) {
	// Get the block from the datastore and decode it.
	storedEthBlock := prepareStoredEthBlock("test_data/eth-block-header-rlp-999999", t)
	ethBlock, err := DecodeEthBlock(storedEthBlock.Cid(), storedEthBlock.RawData())
	checkError(err, t)

	jsonOutput, err := ethBlock.MarshalJSON()
	checkError(err, t)

	var data map[string]interface{}
	err = json.Unmarshal(jsonOutput, &data)
	checkError(err, t)

	// Testing all fields is boring, but can help us to avoid
	// that dreaded regression
	if data["bloom"].(string)[:10] != "0x00000000" {
		t.Fatal("Wrong Bloom")
	}
	if data["coinbase"] != "0x52bc44d5378309ee2abf1539bf71de1b7d7be3b5" {
		t.Fatal("Wrong Coinbase")
	}
	if parseFloat(data["difficulty"]) != "12555463106190" {
		t.Fatal("Wrong Difficulty")
	}
	if data["extra"] != "14MBAwOER2V0aIdnbzEuNC4yhWxpbnV4" {
		t.Fatal("Wrong Extra")
	}
	if parseFloat(data["gaslimit"]) != "3141592" {
		t.Fatal("Wrong Gas limit")
	}
	if parseFloat(data["gasused"]) != "231000" {
		t.Fatal("Wrong Gas used")
	}
	if data["mixdigest"] != "0x5b10f4a08a6c209d426f6158bd24b574f4f7b7aa0099c67c14a1f693b4dd04d0" {
		t.Fatal("Wrong Mix digest")
	}
	if data["nonce"] != "0xf491f46b60fe04b3" {
		t.Fatal("Wrong nonce")
	}
	if parseFloat(data["number"]) != "999999" {
		t.Fatal("Wrong block number")
	}
	if parseMapElement(data["parent"]) != "z43AaGF6wP6uoLFEauru5oLK5JS5MGfNuGDK1xWEpQK4BqkJkL3" {
		t.Fatal("Wrong Parent cid")
	}
	if data["parentHash"] != "0xd33c9dde9fff0ebaa6e71e8b26d2bda15ccf111c7af1b633698ac847667f0fb4" {
		t.Fatal("Wrong Parent hash")
	}
	if data["receiptHash"] != "0x7fa0f6ca2a01823208d80801edad37e3e3a003b55c89319b45eb1f97862ad229" {
		t.Fatal("Wrong Receipt hash")
	}
	if parseMapElement(data["receipts"]) != "z44vkPhhDSTXPAswvC1rdDunzkgZ7FgAAnhGQtNDNDk9m9N2BZA" {
		t.Fatal("Wrong Receipt root cid")
	}
	if parseMapElement(data["root"]) != "z45oqTSAZvPiiPV8hMZDH5fi4NkaAkMYTJC6PmaeWBmYUpbMpoh" {
		t.Fatal("Wrong root hash cid")
	}
	if data["rootHash"] != "0xed98aa4b5b19c82fb35364f08508ae0a6dec665fa57663dca94c5d70554cde10" {
		t.Fatal("Wrong Root Hash")
	}
	if parseFloat(data["time"]) != "1455404037" {
		t.Fatal("Wrong Time")
	}
	if parseMapElement(data["tx"]) != "z443fKyHHMwVy13VXtD4fdRcUXSqkr79Q5E8hcmEravVBq3Dc51" {
		t.Fatal("Wrong Tx root cid")
	}
	if data["txHash"] != "0x447cbd8c48f498a6912b10831cdff59c7fbfcbbe735ca92883d4fa06dcd7ae54" {
		t.Fatal("Wrong Tx root hash")
	}
	if data["uncleHash"] != "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347" {
		t.Fatal("Wrong Uncle hash")
	}
	if parseMapElement(data["uncles"]) != "z43c7o74hjCAqnyneWetkyXU2i5KuGQLbYfVWZMvJMG4VTYABtz" {
		t.Fatal("Wrong Uncle hash cid")
	}
}

func TestEthBlockLinks(t *testing.T) {
	// Get the block from the datastore and decode it.
	storedEthBlock := prepareStoredEthBlock("test_data/eth-block-header-rlp-999999", t)
	ethBlock, err := DecodeEthBlock(storedEthBlock.Cid(), storedEthBlock.RawData())
	checkError(err, t)

	links := ethBlock.Links()
	if links[0].Cid.String() != "z43AaGF6wP6uoLFEauru5oLK5JS5MGfNuGDK1xWEpQK4BqkJkL3" {
		t.Fatal("Wrong cid for parent link")
	}
	if links[1].Cid.String() != "z44vkPhhDSTXPAswvC1rdDunzkgZ7FgAAnhGQtNDNDk9m9N2BZA" {
		t.Fatal("Wrong cid for receipt root link")
	}
	if links[2].Cid.String() != "z45oqTSAZvPiiPV8hMZDH5fi4NkaAkMYTJC6PmaeWBmYUpbMpoh" {
		t.Fatal("Wrong cid for state root link")
	}
	if links[3].Cid.String() != "z443fKyHHMwVy13VXtD4fdRcUXSqkr79Q5E8hcmEravVBq3Dc51" {
		t.Fatal("Wrong cid for transaction root link")
	}
	if links[4].Cid.String() != "z43c7o74hjCAqnyneWetkyXU2i5KuGQLbYfVWZMvJMG4VTYABtz" {
		t.Fatal("Wrong cid for uncle root link")
	}
}

func TestEthBlockResolveBloom(t *testing.T) {
	// Get the block from the datastore and decode it.
	storedEthBlock := prepareStoredEthBlock("test_data/eth-block-header-rlp-999999", t)
	ethBlock, err := DecodeEthBlock(storedEthBlock.Cid(), storedEthBlock.RawData())
	checkError(err, t)

	obj, rest, err := ethBlock.Resolve([]string{"bloom", "anything", "goes", "here"})
	checkError(err, t)

	// The marshaler of types.Bloom should output it as 0x
	bloomInText := fmt.Sprintf("%x", obj.(types.Bloom))
	if bloomInText[:10] != "0000000000" {
		t.Fatal("Wrong Bloom")
	}

	if rest[2] != "here" {
		t.Fatal("Wrong rest of the path returned")
	}
}

func TestEthBlockResolveCoinbase(t *testing.T) {
	// Get the block from the datastore and decode it.
	storedEthBlock := prepareStoredEthBlock("test_data/eth-block-header-rlp-999999", t)
	ethBlock, err := DecodeEthBlock(storedEthBlock.Cid(), storedEthBlock.RawData())
	checkError(err, t)

	obj, rest, err := ethBlock.Resolve([]string{"coinbase", "anything", "goes", "here"})
	checkError(err, t)

	// The marshaler of common.Address should output it as 0x
	coinbaseInText := fmt.Sprintf("%x", obj.(common.Address))
	if coinbaseInText != "52bc44d5378309ee2abf1539bf71de1b7d7be3b5" {
		t.Fatal("Wrong Coinbase")
	}

	if rest[2] != "here" {
		t.Fatal("Wrong rest of the path returned")
	}
}

// TODO
// Test for the following elements resolving
/*
   difficulty
   extra
   gaslimit
   gasused
   mixdigest
   nonce
   number
   parentHash
   receiptHash
   receipts
   root
   rootHash
   time
   tx
   txHash
   uncleHash
   uncles
*/

func TestEthBlockResolveParent(t *testing.T) {
	// Get the block from the datastore and decode it.
	storedEthBlock := prepareStoredEthBlock("test_data/eth-block-header-rlp-999999", t)
	ethBlock, err := DecodeEthBlock(storedEthBlock.Cid(), storedEthBlock.RawData())
	checkError(err, t)

	obj, rest, err := ethBlock.Resolve([]string{"parent", "rest", "of", "the", "path"})
	checkError(err, t)

	lnk, ok := obj.(*node.Link)
	if !ok {
		t.Fatal("Returned object is not a link")
	}

	if lnk.Cid.String() != "z43AaGF6wP6uoLFEauru5oLK5JS5MGfNuGDK1xWEpQK4BqkJkL3" {
		t.Fatal("Wrong parent")
	}

	if rest[3] != "path" {
		t.Fatal("Wrong rest of the path returned")
	}
}

/*
  AUXILIARS
*/

// checkError makes 3 lines into 1.
func checkError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

// parseFloat is a convenience function to test json output
func parseFloat(v interface{}) string {
	return strconv.FormatFloat(v.(float64), 'f', 0, 64)
}

// parseMapElement is a convenience function to tets json output
func parseMapElement(v interface{}) string {
	return v.(map[string]interface{})["/"].(string)
}

// prepareStoredEthBlock reads the block from a file source to get its rawdata
// and computes its cid, for then, feeding it into a new IPLD block function.
// So we can pretend that we got this block from the datastore
func prepareStoredEthBlock(filepath string, t *testing.T) *block.BasicBlock {
	// Prepare the "fetched block". This one is supposed to be in the datastore
	// and given away by github.com/ipfs/go-ipfs/merkledag
	fi, err := os.Open(filepath)
	checkError(err, t)

	b, err := ioutil.ReadAll(fi)
	checkError(err, t)

	c := rawdataToCid(MEthBlock, b)
	storedEthBlock, err := block.NewBlockWithCid(b, c)
	checkError(err, t)

	return storedEthBlock
}

// testEthBlockFields checks the fields of EthBlock one by one.
func testEthBlockFields(ethBlock *EthBlock, t *testing.T) {
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
