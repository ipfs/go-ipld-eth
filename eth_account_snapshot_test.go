package ipldeth

import (
	"fmt"
	"os"
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
  AUXILIARS
*/
func prepareEthAccountSnapshot(t *testing.T) *EthAccountSnapshot {
	fi, err := os.Open("test_data/eth-state-trie-rlp-c9070d")
	checkError(err, t)

	output, err := FromStateTrieRLP(fi)
	checkError(err, t)

	return output.elements[1].(*EthAccountSnapshot)
}
