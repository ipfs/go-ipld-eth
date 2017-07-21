package ipldeth

import (
	"bytes"
	"encoding/hex"
	"os"
	"testing"
)

func TestBlockParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-block.bin")
	if err != nil {
		t.Fatal(err)
	}

	blk, txs, _, _, err := FromRlpBlockMessage(fi)
	if err != nil {
		t.Fatal(err)
	}

	c := blk.Cid()

	_ = c
	// f_ = cmt.Println(c)
	exp := "8c03e3af302c800b4ef1b96d48c5dd23c9410a9e858df5a1e82cdfa5a71895bf"
	hval, err := hex.DecodeString(exp)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.HasSuffix(c.Bytes(), hval) {
		t.Fatal("expected hashes to match")
	}
	_ = txs
}

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
	// fmt.Printf("%x\n", blk.Tx().Bytes())
}
