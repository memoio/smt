package smt

import (
	"bytes"
	"crypto/sha256"
	"hash"
)

var leafPrefix = []byte{0}
var nodePrefix = []byte{1}

type treeHasher struct {
	hasher    hash.Hash
	hashL     int
	zeroValue []byte
}

func newTreeHasher(hasher hash.Hash) *treeHasher {
	var th treeHasher
	th.hashL = sha256.New().Size()
	th.zeroValue = make([]byte, th.pathSize())

	return &th
}

func (th *treeHasher) digest(data []byte) []byte {
	// should instantiates a hasher to avoid race condition
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func (th *treeHasher) path(key []byte) []byte {
	return th.digest(key)
}

func (th *treeHasher) digestLeaf(path []byte, leafData []byte) ([]byte, []byte) {
	value := make([]byte, 0, len(leafPrefix)+len(path)+len(leafData))
	value = append(value, leafPrefix...)
	value = append(value, path...)
	value = append(value, leafData...)

	sum := th.digest(value)

	return sum, value
}

func (th *treeHasher) parseLeaf(data []byte) ([]byte, []byte, []byte) {
	return data[len(leafPrefix) : th.pathSize()+len(leafPrefix)], data[len(leafPrefix)+th.pathSize():], data[len(leafPrefix):]
}

func (th *treeHasher) isLeaf(data []byte) bool {
	return bytes.Equal(data[:len(leafPrefix)], leafPrefix)
}

func (th *treeHasher) digestNode(leftData []byte, rightData []byte) ([]byte, []byte) {
	value := make([]byte, 0, len(nodePrefix)+len(leftData)+len(rightData))
	value = append(value, nodePrefix...)
	value = append(value, leftData...)
	value = append(value, rightData...)

	sum := th.digest(value)

	return sum, value
}

func (th *treeHasher) parseNode(data []byte) ([]byte, []byte) {
	return data[len(nodePrefix) : th.pathSize()+len(nodePrefix)], data[len(nodePrefix)+th.pathSize():]
}

func (th *treeHasher) pathSize() int {
	return th.hashL
}

func (th *treeHasher) placeholder() []byte {
	return th.zeroValue
}
