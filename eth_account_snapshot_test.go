package ipldeth

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"testing"
)

/*
  Block INTERFACE
*/

func TestAccountSnapshotBlockElements(t *testing.T) {
	eas := prepareEthAccountSnapshot(t)

	if fmt.Sprintf("%x", eas.RawData())[:10] != "f84e808a03" {
		t.Fatal("Wrong Data")
	}

	if eas.Cid().String() !=
		"z46FNzJEHSt2Pf9kUR88VNRbjf64BiPEK36iDexE28r81VZ9VMm" {
		t.Fatal("Wrong Cid")
	}
}

func TestAccountSnapshotString(t *testing.T) {
	eas := prepareEthAccountSnapshot(t)

	if eas.String() !=
		"<EthereumAccountSnapshot z46FNzJEHSt2Pf9kUR88VNRbjf64BiPEK36iDexE28r81VZ9VMm>" {
		t.Fatalf("Wrong String()")
	}
}

func TestAccountSnapshotLoggable(t *testing.T) {
	eas := prepareEthAccountSnapshot(t)

	l := eas.Loggable()
	if _, ok := l["type"]; !ok {
		t.Fatal("Loggable map expected the field 'type'")
	}

	if l["type"] != "eth-account-snapshot" {
		t.Fatal("Wrong Loggable 'type' value")
	}
}

/*
  Node INTERFACE
*/
func TestAccountSnapshotResolve(t *testing.T) {
	eas := prepareEthAccountSnapshot(t)

	// Empty path
	obj, rest, err := eas.Resolve([]string{})
	reas, ok := obj.(*EthAccountSnapshot)
	if !ok {
		t.Fatal("Wrong type of returned object")
	}
	if reas.Cid() != eas.Cid() {
		t.Fatal("wrong returned object")
	}
	if rest != nil {
		t.Fatal("rest should be nil")
	}
	if err != nil {
		t.Fatal("err should be nil")
	}

	// len(p) > 1
	badCases := [][]string{
		[]string{"two", "elements"},
		[]string{"here", "three", "elements"},
		[]string{"and", "here", "four", "elements"},
	}

	for _, bc := range badCases {
		obj, rest, err = eas.Resolve(bc)
		if obj != nil {
			t.Fatal("obj should be nil")
		}
		if rest != nil {
			t.Fatal("rest should be nil")
		}
		if err.Error() != fmt.Sprintf("unexpected path elements past %s", bc[0]) {
			t.Fatal("wrong error")
		}
	}

	moreBadCases := []string{
		"i",
		"am",
		"not",
		"an",
		"account",
		"field",
	}
	for _, mbc := range moreBadCases {
		obj, rest, err = eas.Resolve([]string{mbc})
		if obj != nil {
			t.Fatal("obj should be nil")
		}
		if rest != nil {
			t.Fatal("rest should be nil")
		}
		if err.Error() != fmt.Sprintf("no such link") {
			t.Fatal("wrong error")
		}
	}

	goodCases := []string{
		"balance",
		"codeHash",
		"nonce",
		"root",
	}
	for _, gc := range goodCases {
		_, _, err = eas.Resolve([]string{gc})
		if err != nil {
			t.Fatalf("error should be nil %v", gc)
		}
	}

}

func TestAccountSnapshotTree(t *testing.T) {
	eas := prepareEthAccountSnapshot(t)

	// Bad cases
	tree := eas.Tree("non-empty-string", 0)
	if tree != nil {
		t.Fatal("Expected nil to be returned")
	}

	tree = eas.Tree("non-empty-string", 1)
	if tree != nil {
		t.Fatal("Expected nil to be returned")
	}

	tree = eas.Tree("", 0)
	if tree != nil {
		t.Fatal("Expected nil to be returned")
	}

	// Good cases
	tree = eas.Tree("", 1)
	lookupElements := map[string]interface{}{
		"balance":  nil,
		"codeHash": nil,
		"nonce":    nil,
		"root":     nil,
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

func TestAccountSnapshotResolveLink(t *testing.T) {
	eas := prepareEthAccountSnapshot(t)

	// bad case
	obj, rest, err := eas.ResolveLink([]string{"supercalifragilist"})
	if obj != nil {
		t.Fatalf("Expected obj to be nil")
	}
	if rest != nil {
		t.Fatal("Expected rest to be nil")
	}
	if err.Error() != "no such link" {
		t.Fatal("Wrong error")
	}

	// good case
	obj, rest, err = eas.ResolveLink([]string{"nonce"})
	if obj != nil {
		t.Fatalf("Expected obj to be nil")
	}
	if rest != nil {
		t.Fatal("Expected rest to be nil")
	}
	if err.Error() != "resolved item was not a link" {
		t.Fatal("Wrong error")
	}
}

func TestAccountSnapshotCopy(t *testing.T) {
	eas := prepareEthAccountSnapshot(t)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic")
		}
		if r != "dont use this yet" {
			t.Fatal("Expected panic message 'dont use this yet'")
		}
	}()

	_ = eas.Copy()
}

func TestAccountSnapshotLinks(t *testing.T) {
	eas := prepareEthAccountSnapshot(t)

	if eas.Links() != nil {
		t.Fatal("Links() expected to return nil")
	}
}

func TestAccountSnapshotStat(t *testing.T) {
	eas := prepareEthAccountSnapshot(t)

	obj, err := eas.Stat()
	if obj == nil {
		t.Fatal("Expected a not null object node.NodeStat")
	}

	if err != nil {
		t.Fatal("Expected a nil error")
	}
}

func TestAccountSnapshotSize(t *testing.T) {
	eas := prepareEthAccountSnapshot(t)

	size, err := eas.Size()
	if size != uint64(0) {
		t.Fatal("Expected a size equal to 0")
	}

	if err != nil {
		t.Fatal("Expected a nil error")
	}
}

/*
  EthAccountSnapshot functions
*/

func TestAccountSnapshotMarshalJSON(t *testing.T) {
	eas := prepareEthAccountSnapshot(t)

	jsonOutput, err := eas.MarshalJSON()
	checkError(err, t)

	var data map[string]interface{}
	err = json.Unmarshal(jsonOutput, &data)
	checkError(err, t)

	balanceExpression := regexp.MustCompile(`{"balance":16011846000000000000000,`)
	if !balanceExpression.MatchString(string(jsonOutput)) {
		t.Fatal("Balance expression not found")
	}

	code, _ := data["codeHash"].(map[string]interface{})
	if fmt.Sprintf("%s", code["/"]) !=
		"z46gvXAFfCZuCCVi5sXRVYrdHFHQwFw6VdSGumL22MKk4qazAEa" {
		t.Fatal("Wrong Marshaled Value")
	}

	if fmt.Sprintf("%v", data["nonce"]) != "0" {
		t.Fatal("Wrong Marshaled Value")
	}

	root, _ := data["root"].(map[string]interface{})
	if fmt.Sprintf("%s", root["/"]) !=
		"z46gvXAN45x7Y6NtefinhLtJ2UNDwHnbsggxNcuEN89hWvy6ktt" {
		t.Fatal("Wrong Marshaled Value")
	}
}

/*
  AUXILIARS
*/
func prepareEthAccountSnapshot(t *testing.T) *EthAccountSnapshot {
	fi, err := os.Open("test_data/eth-state-trie-rlp-c9070d")
	checkError(err, t)

	output, err := FromStateTrieRLP(fi)
	checkError(err, t)

	return output.elements[1].(*EthAccountSnapshot)
}
