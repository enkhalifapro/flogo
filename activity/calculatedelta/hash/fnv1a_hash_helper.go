package hash

import (
	"encoding/hex"
	"hash"
	"hash/fnv"
)

type FNV1aHelper struct {
	hasher hash.Hash32
}

func NewFNV1aHelper() *FNV1aHelper {
	return &FNV1aHelper{hasher: fnv.New32a()}
}

func (s *FNV1aHelper) GetHashString(bytes []byte) string {

	hasher := fnv.New32a()
	hasher.Write(bytes)
	resultHash := hex.EncodeToString(hasher.Sum(nil))
	return resultHash
}
