package ipldeth

import (
	"os"
	"testing"
)

func TestTxTriesInBlockBodyRlpParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-block-body-rlp-999999")
	if err != nil {
		t.Fatal(err)
	}

	_, _, output, err := FromBlockRLP(fi)
	if err != nil {
		t.Fatal(err)
	}

	if len(output) != 13 {
		t.Fatal("Wrong number of obtained tx trie nodes")
	}
}
