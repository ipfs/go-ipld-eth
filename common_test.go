package ipldeth

import (
	"testing"
)

func TestValidateTriePathGoodCases(t *testing.T) {
	var (
		path []string
		err  error
	)

	goodCases := [][]string{
		[]string{"b", "0d010", "1"},
		[]string{"b", "0", "d", "0", "1", "0", "1"},
		[]string{"0", "1", "1", "B"},
		[]string{"0", "1", "1", "b"},
		[]string{"cC001d4"},
		[]string{"c", "c", "0", "0", "1", "d", "4"},
		[]string{"0123456789abcdef"},
		[]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"},
		[]string{"cC001d4", "nonce"},
		[]string{"c", "c", "0", "0", "1", "d", "4", "nonce"},
		[]string{"cC001d4", "gasPrice"},
		[]string{"c", "c", "0", "0", "1", "d", "4", "gasPrice"},
		[]string{"cC001d4", "gas"},
		[]string{"c", "c", "0", "0", "1", "d", "4", "gas"},
		[]string{"cC001d4", "toAddress"},
		[]string{"c", "c", "0", "0", "1", "d", "4", "toAddress"},
		[]string{"cC001d4", "value"},
		[]string{"c", "c", "0", "0", "1", "d", "4", "value"},
		[]string{"cC001d4", "data"},
		[]string{"c", "c", "0", "0", "1", "d", "4", "data"},
	}

	for i := 0; i < len(goodCases); i = i + 2 {
		path, err = validateTriePath(goodCases[i], getTxFields())
		if err != nil {
			t.Fatal("Unexpected Error")
		}

		if !compareStringSlices(path, goodCases[i+1]) {
			t.Fatal("Wrong returned path")
		}
	}
}

func TestValidateTriePathBadCases(t *testing.T) {
	var (
		path []string
		err  error
	)

	badCases := [][]string{
		[]string{"b", " ", "1"},
		[]string{"", "1", "1", "B"},
		[]string{"cC00--d4"},
		[]string{"0123n56789abcdef"},
		[]string{"012356789abcdef", "banana"},
		[]string{"0", "0", "0", "m0m0ney"},
	}

	for _, bc := range badCases {
		path, err = validateTriePath(bc, getTxFields())
		if err.Error()[:10] != "Unexpected" {
			t.Fatal("Should have an error")
		}

		if path != nil {
			t.Fatal("Should have returned nil path")
		}
	}
}

func TestGetHexIndexGoodCases(t *testing.T) {
	goodCases := map[string]int{
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
		"7": 7,
		"8": 8,
		"9": 9,
		"a": 10,
		"b": 11,
		"c": 12,
		"d": 13,
		"e": 14,
		"f": 15,
	}

	for k, v := range goodCases {
		if getHexIndex(k) != v {
			t.Fatal("Wrong hex index")
		}
	}
}

/*
  AUXILIARS
*/
func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
