package ipldeth

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	block "github.com/ipfs/go-block-format"
)

func TestTxInBlockBodyRlpParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-block-body-rlp-999999")
	checkError(err, t)

	_, output, _, err := FromBlockRLP(fi)
	checkError(err, t)

	if len(output) != 11 {
		t.Fatal("Wrong number of parsed txs")
	}

	// Oh, let's just grab the last element and one from the middle
	testTx05Fields(output[5], t)
	testTx10Fields(output[10], t)
}

func TestTxInBlockHeaderRlpParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-block-header-rlp-999999")
	checkError(err, t)

	_, output, _, err := FromBlockRLP(fi)
	checkError(err, t)

	if len(output) != 0 {
		t.Fatal("No transactions should have been gotten from here")
	}
}

func TestTxInBlockBodyJsonParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-block-body-json-999999")
	checkError(err, t)

	_, output, _, err := FromBlockJSON(fi)
	checkError(err, t)

	if len(output) != 11 {
		t.Fatal("Wrong number of parsed txs")
	}

	testTx05Fields(output[5], t)
	testTx10Fields(output[10], t)
}

func TestDecodeTransaction(t *testing.T) {
	// Prepare the "fetched transaction".
	// This one is supposed to be in the datastore already,
	// and given away by github.com/ipfs/go-ipfs/merkledag
	rawTransactionString :=
		"f86c34850df84758008252089432be343b94f860124dc4fee278fdcbd38c102d88880f25" +
			"8512af0d4000801ba0e9a25c929c26d1a95232ba75aef419a91b470651eb77614695e16c" +
			"5ba023e383a0679fb2fc0d0b0f3549967c0894ee7d947f07d238a83ef745bc3ced5143a4af36"
	rawTransaction, err := hex.DecodeString(rawTransactionString)
	checkError(err, t)
	c := rawdataToCid(MEthTx, rawTransaction)

	// Just to clarify: This `block` is an IPFS block
	storedTransaction, err := block.NewBlockWithCid(rawTransaction, c)
	checkError(err, t)

	// Now the proper test
	ethTransaction, err := DecodeEthTx(storedTransaction.Cid(), storedTransaction.RawData())
	checkError(err, t)

	testTx05Fields(ethTransaction, t)
}

/*
  AUXILIARS
*/

func testTx05Fields(ethTx *EthTx, t *testing.T) {
	// Was the cid calculated?
	if ethTx.Cid().String() != "z44VCrqacegDLXw385vC4tZi84ifPengFdSqbLveMRmsFBeDdNs" {
		t.Fatal("Wrong cid")
	}

	// Do we have the rawdata available?
	if fmt.Sprintf("%x", ethTx.RawData()[:10]) != "f86c34850df847580082" {
		t.Fatal("Wrong Rawdata")
	}

	// Proper Fields of types.Transaction
	if fmt.Sprintf("%x", ethTx.To()) != "32be343b94f860124dc4fee278fdcbd38c102d88" {
		t.Fatal("Wrong Recipient")
	}
	if len(ethTx.Data()) != 0 {
		t.Fatal("Wrong Data")
	}
	if fmt.Sprintf("%v", ethTx.Gas()) != "21000" {
		t.Fatal("Wrong Gas")
	}
	if fmt.Sprintf("%v", ethTx.Value()) != "1091424800000000000" {
		t.Fatal("Wrong Value")
	}
	if fmt.Sprintf("%v", ethTx.Nonce()) != "52" {
		t.Fatal("Wrong Nonce")
	}
	if fmt.Sprintf("%v", ethTx.GasPrice()) != "60000000000" {
		t.Fatal("Wrong Gas Price")
	}
}

func testTx10Fields(ethTx *EthTx, t *testing.T) {
	// Was the cid calculated?
	if ethTx.Cid().String() != "z44VCrqbjszozB5K5Xqm3tm9YDqrWPE5H9QRpKAZRjCLQFbrctT" {
		t.Fatal("Wrong cid")
	}

	// Do we have the rawdata available?
	if fmt.Sprintf("%x", ethTx.RawData()[:10]) != "f8708302a120850ba43b" {
		t.Fatal("Wrong Rawdata")
	}

	// Proper Fields of types.Transaction
	if fmt.Sprintf("%x", ethTx.To()) != "1c51bf013add0857c5d9cf2f71a7f15ca93d4816" {
		t.Fatal("Wrong Recipient")
	}
	if len(ethTx.Data()) != 0 {
		t.Fatal("Wrong Data")
	}
	if fmt.Sprintf("%v", ethTx.Gas()) != "90000" {
		t.Fatal("Wrong Gas")
	}
	if fmt.Sprintf("%v", ethTx.Value()) != "1049756850000000000" {
		t.Fatal("Wrong Value")
	}
	if fmt.Sprintf("%v", ethTx.Nonce()) != "172320" {
		t.Fatal("Wrong Nonce")
	}
	if fmt.Sprintf("%v", ethTx.GasPrice()) != "50000000000" {
		t.Fatal("Wrong Gas Price")
	}
}
