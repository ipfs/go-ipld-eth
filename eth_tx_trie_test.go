package ipldeth

import (
	"encoding/hex"
	"os"
	"testing"

	block "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"
)

func TestTxTriesInBlockBodyJSONParsing(t *testing.T) {
	// HINT: 306 txs
	// cat test_data/eth-block-body-json-4139497 | jsontool | grep transactionIndex | wc -l
	// or, https://etherscan.io/block/4139497
	fi, err := os.Open("test_data/eth-block-body-json-4139497")
	checkError(err, t)

	_, _, output, err := FromBlockJSON(fi)
	checkError(err, t)

	if len(output) != 331 {
		t.Fatal("Wrong number of obtained tx trie nodes")
	}
}

func TestTxTrieDecodeBranch(t *testing.T) {
	ethTxTrie := prepareDecodedEthTxTrieBranch(t)

	if ethTxTrie.nodeKind != "branch" {
		t.Fatal("Wrong nodeKind")
	}

	if len(ethTxTrie.elements) != 17 {
		t.Fatal("Wrong number of elements for a branch node")
	}

	for i, element := range ethTxTrie.elements {
		switch {
		case i < 9:
			if _, ok := element.(*cid.Cid); !ok {
				t.Fatal("Expected element to be a cid")
			}
			continue
		default:
			if element != nil {
				t.Fatal("Expected element to be a nil")
			}
		}
	}
}

func TestTxTrieResolveBranchChildren(t *testing.T) {
	ethTxTrie := prepareDecodedEthTxTrieBranch(t)

	indexes := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}

	for j, index := range indexes {
		obj, rest, err := ethTxTrie.Resolve([]string{index, "nonce"})

		switch {
		case j < 9:
			_, ok := obj.(*node.Link)
			if !ok {
				t.Fatalf("Returned object is not a link (index: %d)", j)
			}

			for k, p := range []string{"nonce"} {
				if rest[k] != p {
					t.Fatal("Wrong rest of the path returned")
				}
			}

			if err != nil {
				t.Fatal("Error should be nil")
			}

		default:
			if obj != nil {
				t.Fatalf("Returned object should have been nil")
			}

			if rest != nil {
				t.Fatalf("Rest of the path returned should be nil")
			}

			if err.Error() != "no such link" {
				t.Fatalf("Wrong error")
			}
		}
	}
}

func TestTxTrieJSONMarshalBranch(t *testing.T) {

}

///// LATER

func TestTxTrieDecodeExtension(t *testing.T) {

}

func TestTxTrieDecodeLeaf(t *testing.T) {

}

/*
  AUXILIARS
*/

// prepareDecodedEthTxTrie simulates an IPLD block available in the datastore,
// checks the source RLP and tests for the absence of errors during the decoding fase.
func prepareDecodedEthTxTrie(branchDataRLP string, t *testing.T) *EthTxTrie {
	b, err := hex.DecodeString(branchDataRLP)
	checkError(err, t)

	c := rawdataToCid(MEthTxTrie, b)

	storedEthTxTrie, err := block.NewBlockWithCid(b, c)
	checkError(err, t)

	ethTxTrie, err := DecodeEthTxTrie(storedEthTxTrie.Cid(), storedEthTxTrie.RawData())
	checkError(err, t)

	return ethTxTrie
}

// prepareDecodedEthTxTrieBranch is a just a helper, to avoid having these lines so many times.
func prepareDecodedEthTxTrieBranch(t *testing.T) *EthTxTrie {
	branchDataRLP :=
		"f90131a051e622bd20e77781a010b9903832e73fd3665e89407ded8c840d8b2db34dd9" +
			"dca0d3f45a40fcad18a6c3d7edbe8e7e92ace9d45e086cbd04a66254b9931375bee1a0" +
			"e15476fc93dc41ef612ac86750dd242d14498c1e48a6ba4fc89fcc501ee7c58ca01363" +
			"826032eeaf1c4540ed2e8e10dc3a34c3fbc4900c7a7c449e69e2ca8a8e1ba094e9d98b" +
			"ebb67807ecd96a6cac608f95a14a07e6a9c06975861e0b86b6c14736a0ec0cfff9d5ab" +
			"a2ac0da8d2c4725bc8253b60f7b6f1c6b4229ea967fcaef319d3a02b652173155b7d9b" +
			"b152ec5d255b82534d3075bcc171a928eba737da9381effaa032a8447e172dc85a1584" +
			"d0f77466ee52a1c00f71caf57e0e1aa01de18a3ca834a0bbc043cc0d03623ba4c7b514" +
			"7d5aca56450b548f797d712d5198f5e8b35f542d8080808080808080"
	return prepareDecodedEthTxTrie(branchDataRLP, t)
}
