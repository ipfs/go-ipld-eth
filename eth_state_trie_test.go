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

func TestStateTrieNodeEvenExtensionParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-state-trie-rlp-eb2f5f")
	checkError(err, t)

	output, err := FromStateTrieRLP(fi)
	checkError(err, t)

	if output.nodeKind != "extension" {
		t.Fatal("Wrong nodeKind")
	}

	if len(output.elements) != 2 {
		t.Fatal("Wrong number of elements for an extension node")
	}

	if fmt.Sprintf("%x", output.elements[0]) != "0d08" {
		t.Fatal("Wrong key")
	}

	if output.elements[1].(*cid.Cid).String() !=
		"z45oqTRzjRMBR6EUKt42Kta9wkJgWxCK22WY2rcKgomw6gAqqaL" {
		t.Fatal("Wrong Value")
	}
}

func TestStateTrieNodeOddExtensionParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-state-trie-rlp-56864f")
	checkError(err, t)

	output, err := FromStateTrieRLP(fi)
	checkError(err, t)

	if output.nodeKind != "extension" {
		t.Fatal("Wrong nodeKind")
	}

	if len(output.elements) != 2 {
		t.Fatal("Wrong number of elements for an extension node")
	}

	if fmt.Sprintf("%x", output.elements[0]) != "02" {
		t.Fatal("Wrong key")
	}

	if output.elements[1].(*cid.Cid).String() !=
		"z45oqTRyJrBDs3CgRfadhvGKM5vvs2hbDE1hSTsQFQDG3XsrWfJ" {
		t.Fatal("Wrong Value")
	}
}

func TestStateTrieNodeEvenLeafParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-state-trie-rlp-0e8b34")
	checkError(err, t)

	output, err := FromStateTrieRLP(fi)
	checkError(err, t)

	if output.nodeKind != "leaf" {
		t.Fatal("Wrong nodeKind")
	}

	if len(output.elements) != 2 {
		t.Fatal("Wrong number of elements for an extension node")
	}

	// bd66f60e5b954e1af93ded1b02cb575ff0ed6d9241797eff7576b0bf0637
	if fmt.Sprintf("%x", output.elements[0].([]byte)[0:10]) != "0b0d06060f06000e050b" {
		t.Fatal("Wrong key")
	}

	if output.elements[1].(*EthAccountSnapshot).String() !=
		"<EthereumAccountSnapshot z46FNzJ7k5MtoWgsk7scSFK224n7wohAv3xuWibghx6qzXCLFfo>" {
		t.Fatal("Wrong Value")
	}
}

func TestStateTrieNodeOddLeafParsing(t *testing.T) {
	fi, err := os.Open("test_data/eth-state-trie-rlp-c9070d")
	checkError(err, t)

	output, err := FromStateTrieRLP(fi)
	checkError(err, t)

	if output.nodeKind != "leaf" {
		t.Fatal("Wrong nodeKind")
	}

	if len(output.elements) != 2 {
		t.Fatal("Wrong number of elements for an extension node")
	}

	// 6c9db9bb545a03425e300f3ee72bae098110336dd3eaf48c20a2e5b6865fc
	if fmt.Sprintf("%x", output.elements[0].([]byte)[0:10]) != "060c090d0b090b0b0504" {
		t.Fatal("Wrong key")
	}

	if output.elements[1].(*EthAccountSnapshot).String() !=
		"<EthereumAccountSnapshot z46FNzJEHSt2Pf9kUR88VNRbjf64BiPEK36iDexE28r81VZ9VMm>" {
		t.Fatal("Wrong Value")
	}
}

/*
  Block INTERFACE
*/
func TestStateTrieBlockElements(t *testing.T) {
	fi, err := os.Open("test_data/eth-state-trie-rlp-d7f897")
	checkError(err, t)

	output, err := FromStateTrieRLP(fi)
	checkError(err, t)

	if fmt.Sprintf("%x", output.RawData())[:10] != "f90211a090" {
		t.Fatal("Wrong Data")
	}

	if output.Cid().String() !=
		"z45oqTS97WG4WsMjquajJ8PB9Ubt3ks7rGmo14P5XWjnPL7LHDM" {
		t.Fatal("Wrong Cid")
	}
}

func TestStateTrieString(t *testing.T) {
	fi, err := os.Open("test_data/eth-state-trie-rlp-d7f897")
	checkError(err, t)

	output, err := FromStateTrieRLP(fi)
	checkError(err, t)

	if output.String() !=
		"<EthereumStateTrie z45oqTS97WG4WsMjquajJ8PB9Ubt3ks7rGmo14P5XWjnPL7LHDM>" {
		t.Fatalf("Wrong String()")
	}
}

func TestStateTrieLoggable(t *testing.T) {
	fi, err := os.Open("test_data/eth-state-trie-rlp-d7f897")
	checkError(err, t)

	output, err := FromStateTrieRLP(fi)
	checkError(err, t)

	l := output.Loggable()
	if _, ok := l["type"]; !ok {
		t.Fatal("Loggable map expected the field 'type'")
	}

	if l["type"] != "eth-state-trie" {
		t.Fatal("Wrong Loggable 'type' value")
	}
}
