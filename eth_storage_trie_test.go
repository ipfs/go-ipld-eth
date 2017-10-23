package ipldeth

import (
	"fmt"
	"os"
	"testing"

	cid "github.com/ipfs/go-cid"
)

/*
  INPUT
  OUTPUT
*/

func TestStorageTrieNodeExtensionParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-storage-trie-rlp-113049")
	checkError(err, t)

	output, err := FromStateTrieRLP(fi)
	checkError(err, t)

	if output.nodeKind != "extension" {
		t.Fatal("Wrong nodeKind")
	}

	if len(output.elements) != 2 {
		t.Fatal("Wrong number of elements for an extension node")
	}

	if fmt.Sprintf("%x", output.elements[0]) != "0a" {
		t.Fatal("Wrong key")
	}

	if output.elements[1].(*cid.Cid).String() !=
		"z45oqTS5gGMhx63LYCz1cgTFEEPPd9o6wX5py1tBY6Mpzbez12t" {
		t.Fatal("Wrong Value")
	}
}

func TestStateTrieNodeLeafParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-storage-trie-rlp-ffbcad")
	checkError(err, t)

	output, err := FromStorageTrieRLP(fi)
	checkError(err, t)

	if output.nodeKind != "leaf" {
		t.Fatal("Wrong nodeKind")
	}

	if len(output.elements) != 2 {
		t.Fatal("Wrong number of elements for an leaf node")
	}

	// 2ee1ae9c502e48e0ed528b7b39ac569cef69d7844b5606841a7f3fe898a2
	if fmt.Sprintf("%x", output.elements[0].([]byte)[:10]) != "020e0e010a0e090c0500" {
		t.Fatal("Wrong key")
	}

	if fmt.Sprintf("%x", output.elements[1]) != "89056c31f304b2530000" {
		t.Fatal("Wrong Value")
	}
}

func TestStateTrieNodeBranchParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-storage-trie-rlp-ffc25c")
	checkError(err, t)

	output, err := FromStateTrieRLP(fi)
	checkError(err, t)

	if output.nodeKind != "branch" {
		t.Fatal("Wrong nodeKind")
	}

	if len(output.elements) != 17 {
		t.Fatal("Wrong number of elements for an branch node")
	}

	if fmt.Sprintf("%s", output.elements[4]) !=
		"z45oqTRvTxd9n2n4jvEpCrAjzoUvDoficbHN3j7ha59BhsNDSXj" {
		t.Fatal("Wrong Cid")
	}

	if fmt.Sprintf("%s", output.elements[10]) !=
		"z45oqTSBnuNXCDvZzfESmWyJLUvyMSWj71xctFG15RArd6QannU" {
		t.Fatal("Wrong Cid")
	}
}

/*
  Block INTERFACE
*/
func TestStorageTrieBlockElements(t *testing.T) {
	fi, err := os.Open("test_data/eth-storage-trie-rlp-ffbcad")
	checkError(err, t)

	output, err := FromStateTrieRLP(fi)
	checkError(err, t)

	if fmt.Sprintf("%x", output.RawData())[:10] != "eb9f202ee1" {
		t.Fatal("Wrong Data")
	}

	if output.Cid().String() !=
		"z45oqTSBnjagWJLxnbpFLK7hdHQzneN5SagVwJB1cNvFvw68QmX" {
		t.Fatal("Wrong Cid")
	}
}

func TestStorageTrieString(t *testing.T) {
	fi, err := os.Open("test_data/eth-storage-trie-rlp-ffbcad")
	checkError(err, t)

	output, err := FromStateTrieRLP(fi)
	checkError(err, t)

	if output.String() !=
		"<EthereumStateTrie z45oqTSBnjagWJLxnbpFLK7hdHQzneN5SagVwJB1cNvFvw68QmX>" {
		t.Fatalf("Wrong String()")
	}
}

func TestStorageTrieLoggable(t *testing.T) {
	fi, err := os.Open("test_data/eth-storage-trie-rlp-ffbcad")
	checkError(err, t)

	output, err := FromStorageTrieRLP(fi)
	checkError(err, t)

	l := output.Loggable()
	if _, ok := l["type"]; !ok {
		t.Fatal("Loggable map expected the field 'type'")
	}

	if l["type"] != "eth-storage-trie" {
		t.Fatal("Wrong Loggable 'type' value")
	}
}
