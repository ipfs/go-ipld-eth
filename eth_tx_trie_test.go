package ipldeth

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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

func TestTxTrieDecodeExtension(t *testing.T) {
	ethTxTrie := prepareDecodedEthTxTrieExtension(t)

	if ethTxTrie.nodeKind != "extension" {
		t.Fatal("Wrong nodeKind")
	}

	if len(ethTxTrie.elements) != 2 {
		t.Fatal("Wrong number of elements for an extension node")
	}

	if fmt.Sprintf("%x", ethTxTrie.elements[0].([]byte)) != "0001" {
		t.Fatal("Wrong key")
	}

	if ethTxTrie.elements[1].(*cid.Cid).String() !=
		"z443fKyJaFfaE7Hsozvv7HGEHqNWPEhkNgzgnXjVKdxqCE74PgF" {
		t.Fatal("Wrong Value")
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

func TestTxTrieResolveExtensionReference(t *testing.T) {
	ethTxTrie := prepareDecodedEthTxTrieExtension(t)

	badCases := []string{"0", "00"}

	for _, bc := range badCases {
		obj, rest, err := ethTxTrie.Resolve([]string{bc})
		if obj != nil {
			t.Fatal("obj should be nil")
		}

		if rest != nil {
			t.Fatal("rest should be nil")
		}

		if err.Error() != "no such link in this extension" {
			t.Fatalf("Wrong error")
		}
	}

	goodCases := []string{"01", "01a", "01ab"}
	for _, gc := range goodCases {
		obj, rest, err := ethTxTrie.Resolve([]string{gc})
		_, ok := obj.(*node.Link)
		if !ok {
			t.Fatalf("Returned object is not a link")
		}

		if strings.Join(rest, "") != gc[2:] {
			t.Fatal("Wrong rest of the path returned")
		}

		if err != nil {
			t.Fatal("Error should be nil")
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

			if err.Error() != "no such link in this branch" {
				t.Fatalf("Wrong error")
			}
		}
	}
}

func TestTxTrieJSONMarshalExtension(t *testing.T) {
	ethTxTrie := prepareDecodedEthTxTrieExtension(t)

	jsonOutput, err := ethTxTrie.MarshalJSON()
	checkError(err, t)

	var data map[string]interface{}
	err = json.Unmarshal(jsonOutput, &data)
	checkError(err, t)

	if parseMapElement(data["01"]) !=
		"z443fKyJaFfaE7Hsozvv7HGEHqNWPEhkNgzgnXjVKdxqCE74PgF" {
		t.Fatal("Wrong Marshaled Value")
	}

	if data["type"] != "extension" {
		t.Fatal("Expected type to be extension")
	}
}

func TestTxTrieJSONMarshalBranch(t *testing.T) {
	ethTxTrie := prepareDecodedEthTxTrieBranch(t)

	jsonOutput, err := ethTxTrie.MarshalJSON()
	checkError(err, t)

	var data map[string]interface{}
	err = json.Unmarshal(jsonOutput, &data)
	checkError(err, t)

	desiredValues := map[string]string{
		"0": "z443fKyJBiTCxynCqKP1r3BSvJ4nvR4bSEpWFMc7ZJ57L6NJUdH",
		"1": "z443fKySwQ2WU6av9YcvidCRCYcBrcY1FbsJfdxtTTeKbpiZD8k",
		"2": "z443fKyTqcL3923Cwqeun2Lo1qs9MPXNV16KFJBHRs6ghNHaFpf",
		"3": "z443fKyDyheaZ5qTSjSS6XLj6trLWasneACqrkBfwNpnjN2Fuia",
		"4": "z443fKyNhK436C7wMxoiM9NfjcnHpmdWgbW6CKvtA4f9kUnoD9P",
		"5": "z443fKyUZTcKeGxvmCfecLxAF8rHEAzCFNVaTwonX2Atd6BB4CS",
		"6": "z443fKyFbQsGGz5fuym6Gv8hyHErR962okHt1zTNKwXebjwUo3w",
		"7": "z443fKyG5m6cHmnhBfi4qNvRXNmL18w71XZGxJifbtUPyUNfk5Z",
		"8": "z443fKyRJvB8PQEdWTL44qqoo2DeZr8QwkasSAfEcWJ6uDUWyh6",
	}

	for k, v := range desiredValues {
		if parseMapElement(data[k]) != v {
			t.Fatal("Wrong Marshaled Value")
		}
	}

	for _, v := range []string{"a", "b", "c", "d", "e", "f"} {
		if data[v] != nil {
			t.Fatal("Expected value to be nil")
		}
	}

	if data["type"] != "branch" {
		t.Fatal("Expected type to be branch")
	}
}

func TestTxTrieLinksBranch(t *testing.T) {
	ethTxTrie := prepareDecodedEthTxTrieBranch(t)

	desiredValues := []string{
		"z443fKyJBiTCxynCqKP1r3BSvJ4nvR4bSEpWFMc7ZJ57L6NJUdH",
		"z443fKySwQ2WU6av9YcvidCRCYcBrcY1FbsJfdxtTTeKbpiZD8k",
		"z443fKyTqcL3923Cwqeun2Lo1qs9MPXNV16KFJBHRs6ghNHaFpf",
		"z443fKyDyheaZ5qTSjSS6XLj6trLWasneACqrkBfwNpnjN2Fuia",
		"z443fKyNhK436C7wMxoiM9NfjcnHpmdWgbW6CKvtA4f9kUnoD9P",
		"z443fKyUZTcKeGxvmCfecLxAF8rHEAzCFNVaTwonX2Atd6BB4CS",
		"z443fKyFbQsGGz5fuym6Gv8hyHErR962okHt1zTNKwXebjwUo3w",
		"z443fKyG5m6cHmnhBfi4qNvRXNmL18w71XZGxJifbtUPyUNfk5Z",
		"z443fKyRJvB8PQEdWTL44qqoo2DeZr8QwkasSAfEcWJ6uDUWyh6",
	}

	links := ethTxTrie.Links()

	for i, v := range desiredValues {
		if links[i].Cid.String() != v {
			t.Fatalf("Wrong cid for link %d", i)
		}
	}
}

func TessTxTrieTreeBadParams(t *testing.T) {
	ethTxTrie := prepareDecodedEthTxTrieBranch(t)

	tree := ethTxTrie.Tree("non-empty-string", 0)
	if tree != nil {
		t.Fatal("Expected nil to be returned")
	}

	tree = ethTxTrie.Tree("non-empty-string", 1)
	if tree != nil {
		t.Fatal("Expected nil to be returned")
	}

	tree = ethTxTrie.Tree("", 0)
	if tree != nil {
		t.Fatal("Expected nil to be returned")
	}
}

func TestTxTrieTreeBranch(t *testing.T) {
	ethTxTrie := prepareDecodedEthTxTrieBranch(t)

	tree := ethTxTrie.Tree("", -1)

	lookupElements := map[string]interface{}{
		"0": nil,
		"1": nil,
		"2": nil,
		"3": nil,
		"4": nil,
		"5": nil,
		"6": nil,
		"7": nil,
		"8": nil,
	}

	if len(tree) != len(lookupElements) {
		t.Fatalf("Wrong number of elements. Got %d. Expecting %d", len(tree), len(lookupElements))
	}

	for _, te := range tree {
		if _, ok := lookupElements[te]; !ok {
			t.Fatalf("Unexpected Element: %v", te)
		}
	}
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

func prepareDecodedEthTxTrieExtension(t *testing.T) *EthTxTrie {
	extensionDataRLP :=
		"e4820001a057ac34d6471cc3f5c6ab992c4c0fe5ec131d8d9961fe6d5de8e5e367513243b4"
	return prepareDecodedEthTxTrie(extensionDataRLP, t)
}

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
