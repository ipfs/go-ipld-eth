package ipldeth

import (
	cid "gx/ipfs/QmTprEaAA2A9bst5XH7exuyi5KzNMK3SEDNN8rBDnKWcUS/go-cid"
	mh "gx/ipfs/QmU9a9NV9RdPNwZQDYd5uKsm6N6LJLSvLbywDDYFbaaC6P/go-multihash"

	common "github.com/ethereum/go-ethereum/common"
)

func castCommonHash(h common.Hash, codec uint64) *cid.Cid {
	mhash, err := mh.Encode(h[:], mh.KECCAK_256)
	if err != nil {
		panic(err)
	}

	return cid.NewCidV1(codec, mhash)
}
