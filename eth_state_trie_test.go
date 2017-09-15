package ipldeth

import (
	"fmt"
	"os"
	"testing"

	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"
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

/*
  TRIE NODE (Through EthStateTrie)
  Node INTERFACE
*/

func TestTraverseStateTrieWithResolve(t *testing.T) {
	var err error

	stMap := prepareStateTrieMap(t)

	// This is the cid of the root of the block 0
	// z45oqTS97WG4WsMjquajJ8PB9Ubt3ks7rGmo14P5XWjnPL7LHDM
	currentNode := stMap["z45oqTS97WG4WsMjquajJ8PB9Ubt3ks7rGmo14P5XWjnPL7LHDM"]

	// This is the path we want to traverse
	// The eth address is 0x5abfec25f74cd88437631a7731906932776356f9
	// Its keccak-256 is cdd3e25edec0a536a05f5e5ab90a5603624c0ed77453b2e8f955cf8b43d4d0fb
	// We use the keccak-256(addr) to traverse the state trie in ethereum.
	var traversePath []string
	for _, s := range "cdd3e25edec0a536a05f5e5ab90a5603624c0ed77453b2e8f955cf8b43d4d0fb" {
		traversePath = append(traversePath, string(s))
	}
	traversePath = append(traversePath, "balance")

	var obj interface{}
	for {
		obj, traversePath, err = currentNode.Resolve(traversePath)
		link, ok := obj.(*node.Link)
		if !ok {
			break
		}
		if err != nil {
			t.Fatal("Error should be nil")
		}

		currentNode = stMap[link.Cid.String()]
		if currentNode == nil {
			t.Fatal("state trie node not found in memory map")
		}
	}

	if fmt.Sprintf("%v", obj) != "11901484239480000000000000" {
		t.Fatal("Wrong value, expected a balance")
	}
}

func prepareStateTrieMap(t *testing.T) map[string]*EthStateTrie {
	filepaths := []string{
		"test_data/eth-state-trie-rlp-0e8b34",
		"test_data/eth-state-trie-rlp-56864f",
		"test_data/eth-state-trie-rlp-6fc2d7",
		"test_data/eth-state-trie-rlp-727994",
		"test_data/eth-state-trie-rlp-c9070d",
		"test_data/eth-state-trie-rlp-d5be90",
		"test_data/eth-state-trie-rlp-d7f897",
		"test_data/eth-state-trie-rlp-eb2f5f",
	}

	out := make(map[string]*EthStateTrie)

	for _, fp := range filepaths {
		fi, err := os.Open(fp)
		checkError(err, t)

		stateTrieNode, err := FromStateTrieRLP(fi)
		checkError(err, t)

		out[stateTrieNode.Cid().String()] = stateTrieNode

		// DEBUG
		// fmt.Printf("%x")
		// DEBUG
	}

	return out
}
