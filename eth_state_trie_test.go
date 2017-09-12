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
		t.Fatalf("Wrong Value")
	}
}

func TestStateTrieNodeOddExtensionParsing(t *testing.T) {

}

func TestStateTrieNodeEvenLeafParsing(t *testing.T) {

}

func TestStateTrieNodeOddLeafParsing(t *testing.T) {

}

/*
  AUXILIARS
*/
