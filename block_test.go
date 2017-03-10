package ipldeth

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
)

func TestBlockParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-block.bin")
	if err != nil {
		t.Fatal(err)
	}

	blk, _, _, _, err := FromRlpBlockMessage(fi)
	if err != nil {
		t.Fatal(err)
	}

	c := blk.Cid()
	fmt.Println(c)
	exp := "8c03e3af302c800b4ef1b96d48c5dd23c9410a9e858df5a1e82cdfa5a71895bf"
	hval, err := hex.DecodeString(exp)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.HasSuffix(c.Bytes(), hval) {
		t.Fatal("expected hashes to match")
	}
}

func TestBlockWithTxParsing(t *testing.T) {
	fi, err := os.Open("test_data/block-with-txs.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer fi.Close()

	blk, txs, _, uncles, err := FromRlpBlockMessage(fi)
	if err != nil {
		t.Fatal(err)
	}

	_ = txs
	_ = blk
	_ = uncles
}
