package ipldeth

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/rlp"
	cid "github.com/ipfs/go-cid"
	node "github.com/ipfs/go-ipld-format"
)

// TrieNode is the general abstraction for
//ethereum IPLD trie nodes.
type TrieNode struct {
	// leaf, extension or branch
	nodeKind string

	// If leaf or extension: [0] is key, [1] is val.
	// If branch: [0] - [16] are children.
	elements []interface{}

	// IPLD block information
	cid     *cid.Cid
	rawdata []byte
}

/*
  OUTPUT
*/

type trieNodeLeafDecoder func([]interface{}) ([]interface{}, error)

// decodeTrieNode returns a TrieNode object from an IPLD block's
// cid and rawdata.
func decodeTrieNode(c *cid.Cid, b []byte,
	leafDecoder trieNodeLeafDecoder) (*TrieNode, error) {
	var (
		i, decoded, elements []interface{}
		nodeKind             string
		err                  error
	)

	err = rlp.DecodeBytes(b, &i)
	if err != nil {
		return nil, err
	}

	codec := c.Type()
	switch len(i) {
	case 2:
		nodeKind, decoded, err = decodeCompactKey(i)
		if err != nil {
			return nil, err
		}

		if nodeKind == "extension" {
			elements, err = parseTrieNodeExtension(decoded, codec)
		}
		if nodeKind == "leaf" {
			elements, err = leafDecoder(decoded)
		}
		if err != nil {
			return nil, err
		}
	case 17:
		nodeKind = "branch"
		elements, err = parseTrieNodeBranch(i, codec)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown trie node type")
	}

	return &TrieNode{
		nodeKind: nodeKind,
		elements: elements,
		rawdata:  b,
		cid:      c,
	}, nil
}

// decodeCompactKey takes a compact key, and returns its nodeKind and value.
func decodeCompactKey(i []interface{}) (string, []interface{}, error) {
	first := i[0].([]byte)
	last := i[1].([]byte)

	switch first[0] / 16 {
	case '\x00':
		return "extension", []interface{}{
			nibbleToByte(first)[2:],
			last,
		}, nil
	case '\x01':
		return "extension", []interface{}{
			nibbleToByte(first)[1:],
			last,
		}, nil
	case '\x02':
		return "leaf", []interface{}{
			nibbleToByte(first)[2:],
			last,
		}, nil
	case '\x03':
		return "leaf", []interface{}{
			nibbleToByte(first)[1:],
			last,
		}, nil
	default:
		return "", nil, fmt.Errorf("unknown hex prefix")
	}
}

// parseTrieNodeExtension helper improves readability
func parseTrieNodeExtension(i []interface{}, codec uint64) ([]interface{}, error) {
	return []interface{}{
		i[0].([]byte),
		keccak256ToCid(codec, i[1].([]byte)),
	}, nil
}

// parseTrieNodeBranch helper improves readability
func parseTrieNodeBranch(i []interface{}, codec uint64) ([]interface{}, error) {
	var out []interface{}

	for _, vi := range i {
		v := vi.([]byte)

		switch len(v) {
		case 0:
			out = append(out, nil)
		case 32:
			out = append(out, keccak256ToCid(codec, v))
		default:
			return nil, fmt.Errorf("unrecognized object: %v", v)
		}
	}

	return out, nil
}

/*
  Node INTERFACE
*/

// Resolve resolves a path through this node, stopping at any link boundary
// and returning the object found as well as the remaining path to traverse
func (t *TrieNode) Resolve(p []string) (interface{}, []string, error) {
	p, err := validateTriePath(p, getTxFields())
	if err != nil {
		return nil, nil, err
	}

	switch t.nodeKind {
	case "extension":
		nibblesCount := checkPathNibbles(t.elements[0].([]byte), p)
		if nibblesCount == -1 {
			return nil, nil, fmt.Errorf("no such link in this extension")
		}
		return &node.Link{Cid: t.elements[1].(*cid.Cid)}, p[nibblesCount:], nil
	case "branch":
		child := t.elements[getHexIndex(p[0])]
		if child != nil {
			return &node.Link{Cid: child.(*cid.Cid)}, p[1:], nil
		}
		return nil, nil, fmt.Errorf("no such link in this branch")
	default:
		return nil, nil, fmt.Errorf("nodeKind case not implemented")
	}
}

// Tree lists all paths within the object under 'path', and up to the given depth.
// To list the entire object (similar to `find .`) pass "" and -1
func (t *TrieNode) Tree(p string, depth int) []string {
	if p != "" || depth == 0 {
		return nil
	}

	var out []string

	switch t.nodeKind {
	case "extension":
		var val string
		for _, e := range t.elements[0].([]byte) {
			val += fmt.Sprintf("%x", e)
		}
		return []string{val}
	case "branch":
		for i, elem := range t.elements {
			if _, ok := elem.(*cid.Cid); ok {
				out = append(out, fmt.Sprintf("%x", i))
			}
		}
		return out

	default:
		return nil
	}
}

// ResolveLink is a helper function that calls resolve and asserts the
// output is a link
func (t *TrieNode) ResolveLink(p []string) (*node.Link, []string, error) {
	obj, rest, err := t.Resolve(p)
	if err != nil {
		return nil, nil, err
	}

	lnk, ok := obj.(*node.Link)
	if !ok {
		return nil, nil, fmt.Errorf("was not a link")
	}

	return lnk, rest, nil
}

// Copy will go away. It is here to comply with the interface.
func (t *TrieNode) Copy() node.Node {
	panic("dont use this yet")
}

// Links is a helper function that returns all links within this object
func (t *TrieNode) Links() []*node.Link {
	var out []*node.Link

	for _, i := range t.elements {
		c, ok := i.(*cid.Cid)
		if ok {
			out = append(out, &node.Link{Cid: c})
		}
	}

	return out
}

// Stat will go away. It is here to comply with the interface.
func (t *TrieNode) Stat() (*node.NodeStat, error) {
	return &node.NodeStat{}, nil
}

// Size will go away. It is here to comply with the interface.
func (t *TrieNode) Size() (uint64, error) {
	return 0, nil
}

/*
  TrieNode functions
*/

// MarshalJSON processes the transaction trie into readable JSON format.
func (t *TrieNode) MarshalJSON() ([]byte, error) {
	var out map[string]interface{}

	switch t.nodeKind {
	case "extension":
		fallthrough
	case "leaf":
		var val string
		for _, e := range t.elements[0].([]byte) {
			val += fmt.Sprintf("%x", e)
		}
		out = map[string]interface{}{
			"type": t.nodeKind,
			val:    t.elements[1],
		}
	case "branch":
		out = map[string]interface{}{
			"type": "branch",
			"0":    t.elements[0],
			"1":    t.elements[1],
			"2":    t.elements[2],
			"3":    t.elements[3],
			"4":    t.elements[4],
			"5":    t.elements[5],
			"6":    t.elements[6],
			"7":    t.elements[7],
			"8":    t.elements[8],
			"9":    t.elements[9],
			"a":    t.elements[10],
			"b":    t.elements[11],
			"c":    t.elements[12],
			"d":    t.elements[13],
			"e":    t.elements[14],
			"f":    t.elements[15],
		}
	default:
		return nil, fmt.Errorf("nodeKind %s not supported", t.nodeKind)
	}

	return json.Marshal(out)
}

// nibbleToByte expands the nibbles of a byte slice into their own bytes.
func nibbleToByte(k []byte) []byte {
	var out []byte

	for _, b := range k {
		out = append(out, b/16)
		out = append(out, b%16)
	}

	return out
}

// validateTriePath takes a trie path, checking whether each element represents
// an hexadecimal character, and returns a slice of one hex character elements,
// allowing the input of paths such as /b/0d010/1 /0/1/1/b /cc001d4 possible.
func validateTriePath(p []string, specialFields map[string]interface{}) ([]string, error) {
	var (
		testString string
		output     []string
	)

	//
	lastValue := p[len(p)-1]
	if _, ok := specialFields[lastValue]; ok {
		// Remove this lastValue and add it after the validation.
		// Examples of lastValue: nonce, gasPrice for txs. balance for states.
		p = p[:len(p)-1]
	} else {
		lastValue = ""
	}

	for _, v := range p {
		if v == "" {
			return nil, fmt.Errorf("Unexpected blank element in path")
		}
		testString += v
	}

	testString = strings.ToLower(testString)

	for _, v := range testString {
		c := byte(v)

		switch {
		case '0' <= c && c <= '9':
			fallthrough
		case 'a' <= c && c <= 'f':
			output = append(output, string(c))
		default:
			return nil, fmt.Errorf("Unexpected character in path: %x", c)
		}
	}

	// Recover the last value
	if lastValue != "" {
		output = append(output, lastValue)
	}

	return output, nil
}

// checkPathNibbles tests whether the given path can resolve the trie node
// element key, returning the number of nibbles the key has if succeed.
func checkPathNibbles(nibbles []byte, p []string) int {
	if len(p) < len(nibbles) {
		return -1
	}

	for i, n := range nibbles {
		if p[i] != fmt.Sprintf("%x", n) {
			return -1
		}
	}

	return len(nibbles)
}

// getHexIndex returns to you the integer 0 - 15 equivalent to your
// string character if applicable, or -1 otherwise.
func getHexIndex(s string) int {
	if len(s) != 1 {
		return -1
	}

	c := byte(s[0])
	switch {
	case '0' <= c && c <= '9':
		return int(c - '0')
	case 'a' <= c && c <= 'f':
		return int(c - 'a' + 10)
	}

	return -1
}
