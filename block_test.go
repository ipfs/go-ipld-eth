package ipldeth

import (
	"fmt"
	"os"
	"testing"
)

func TestBlockHeaderRlpParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-block-4052384")
	if err != nil {
		t.Fatal(err)
	}

	bh, err := FromBlockHeaderRLP(fi)
	if err != nil {
		t.Fatal(err)
	}

	// See whether it parsed the elements of the header
	if fmt.Sprintf("%x", bh.ParentHash) != "04c59985a1c28774f80e4d27441fbd14bef99bdecf078d80fae0b559d089d670" {
		t.Fatal("Wrong ParentHash")
	}
	if fmt.Sprintf("%x", bh.UncleHash) != "1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347" {
		t.Fatal("Wrong UncleHash")
	}
	if fmt.Sprintf("%x", bh.Coinbase) != "ea674fdde714fd979de3edf0f56aa9716b898ec8" {
		t.Fatal("Wrong Coinbase")
	}
	if fmt.Sprintf("%x", bh.Root) != "221696302b99a3a18740abd59c715d5b543026e69e1a6decddd2a8c0b47520cb" {
		t.Fatal("Wrong Root")
	}
	if fmt.Sprintf("%x", bh.TxHash) != "d8d973df5c50a1ae11cd9e9cdf5ab81612cece33bc6eff70a7beb3fbd94e3458" {
		t.Fatal("Wrong TxHash")
	}
	if fmt.Sprintf("%x", bh.ReceiptHash) != "7e859200171588f01534b1e5b09d1252896e1242f685ecf8186571d471d4a3ae" {
		t.Fatal("Wrong ReceiptHash")
	}
	if fmt.Sprintf("%x", bh.Bloom[71:76]) != "0280000004" { // You wouldn't want me to print out the whole bloom field?
		t.Fatal("Wrong Bloom")
	}
	if bh.Difficulty.String() != "1283310643319628" {
		t.Fatal("Wrong Difficulty")
	}
	if bh.Number.String() != "4052384" {
		t.Fatal("Wrong Number")
	}
	if bh.GasLimit.String() != "6719052" {
		t.Fatal("Wrong Gas Limit")
	}
	if bh.GasUsed.String() != "562490" {
		t.Fatal("Wrong Gas Used")
	}
	if bh.Time.String() != "1500627255" {
		t.Fatal("Wrong Time")
	}
	if fmt.Sprintf("%s", bh.Extra) != "ethermine-eu5" {
		t.Fatal("Wrong Extra")
	}
	if fmt.Sprintf("%x", bh.Nonce) != "1ac84bc00f34b563" {
		t.Fatal("Wrong Nonce")
	}
	if fmt.Sprintf("%x", bh.MixDigest) != "1872a178541e5e263a1a68d797a72baa0bcbac1500b29d51491c586f17a0fab2" {
		t.Fatal("Wrong MixDigest")
	}
}

func TestCidCreationAfterParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-block-4052384")
	if err != nil {
		t.Fatal(err)
	}

	bh, err := FromBlockHeaderRLP(fi)
	if err != nil {
		t.Fatal(err)
	}

	if bh.Cid().String() != "z43AaGEwLuiGeeRYszh2ZVtAe92HK796zn8Qz5REq7ztM1ZBz7d" {
		t.Fatal("Wrong Cid")
	}
}

/*
func TestBlockWithTxParsing(t *testing.T) {
	fi, err := os.Open("test_data/block-with-txs.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer fi.Close()

	blk, txs, tries, uncles, err := FromRlpBlockMessage(fi)
	if err != nil {
		t.Fatal(err)
	}

	_ = txs
	_ = blk
	_ = uncles
	_ = tries
}

func TestBlockWithOddTransactions(t *testing.T) {
	fi, err := os.Open("test_data/odd_block.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer fi.Close()

	blk, err := DecodeBlock(fi)
	if err != nil {
		t.Fatal(err)
	}

	_ = blk
}
*/
