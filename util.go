package ipldeth

import (
	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"

	common "github.com/ethereum/go-ethereum/common"
)

func castCommonHash(h common.Hash, codec uint64) *cid.Cid {
	mhash, err := mh.Encode(h[:], mh.KECCAK_256)
	if err != nil {
		panic(err)
	}

	return cid.NewCidV1(codec, mhash)
}
