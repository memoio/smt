// Package smt implements a Sparse Merkle tree.
package smt

import(
    "hash"
)

const left = 0
const right = 1
var defaultValue = []byte{0}

// SparseMerkleTree is a Sparse Merkle tree.
type SparseMerkleTree struct {
    hasher hash.Hash
    ms MapStore
    root []byte
}

// Initialise a Sparse Merkle tree on an empty MapStore.
func NewSparseMerkleTree(ms MapStore, hasher hash.Hash) *SparseMerkleTree {
    smt := SparseMerkleTree{
        hasher: hasher,
        ms: ms,
    }

    for i := 0; i < smt.depth() - 1; i++ {
        ms.Put(smt.defaultNode(i), append(smt.defaultNode(i + 1), smt.defaultNode(i + 1)...))
    }

    ms.Put(smt.defaultNode(255), defaultValue)

    rootValue := append(smt.defaultNode(0), smt.defaultNode(0)...)
    rootHash := smt.digest(rootValue)
    ms.Put(rootHash, rootValue)
    smt.root = rootHash

    return &smt
}

func (smt *SparseMerkleTree) depth() int {
    return smt.keySize() * 8
}

func (smt *SparseMerkleTree) keySize() int {
    return smt.hasher.Size()
}

func (smt *SparseMerkleTree) defaultNode(height int) []byte {
    return defaultNodes(smt.hasher)[height]
}

func (smt *SparseMerkleTree) digest(data []byte) []byte {
    smt.hasher.Write(data)
    sum := smt.hasher.Sum(nil)
    smt.hasher.Reset()
    return sum
}

// Get gets a key from the tree.
func (smt *SparseMerkleTree) Get(key []byte) ([]byte, error) {
    path := smt.digest(key)
    currentHash := smt.root
    for i := 0; i < smt.depth(); i++ {
        currentValue, err := smt.ms.Get(currentHash)
        if err != nil {
            return nil, err
        }
        if hasBit(path, i) == right {
            currentHash = currentValue[smt.keySize():]
        } else {
            currentHash = currentValue[:smt.keySize()]
        }
    }

    value, err := smt.ms.Get(currentHash)
    if err != nil {
        return nil, err
    }

    return value, nil
}

// Update sets a new value for a key in the tree.
func (smt *SparseMerkleTree) Update(key []byte, value []byte) error {
    path := smt.digest(key)
    sideNodes, err := smt.sideNodes(path)
    if err != nil {
        return err
    }

    currentHash := smt.digest(value)
    smt.ms.Put(currentHash, value)
    currentValue := currentHash

    for i := smt.depth() - 1; i >= 0; i-- {
        if hasBit(path, i) == right {
            currentValue = append(sideNodes[i], currentValue...)
        } else {
            currentValue = append(currentValue, sideNodes[i]...)
        }
        currentHash = smt.digest(currentValue)
        err := smt.ms.Put(currentHash, currentValue)
        if err != nil {
            return err
        }
        currentValue = currentHash
    }

    smt.root = currentHash
    return nil
}

func (smt *SparseMerkleTree) sideNodes(path []byte) ([][]byte, error) {
    currentValue, err := smt.ms.Get(smt.root)
    if err != nil {
        return nil, err
    }

    sideNodes := make([][]byte, smt.depth())
    for i := 0; i < smt.depth(); i++ {
        if hasBit(path, i) == right {
            sideNodes[i] = currentValue[:smt.keySize()]
            currentValue, err = smt.ms.Get(currentValue[smt.keySize():])
            if err != nil {
                return nil, err
            }
        } else {
            sideNodes[i] = currentValue[smt.keySize():]
            currentValue, err = smt.ms.Get(currentValue[:smt.keySize()])
            if err != nil {
                return nil, err
            }
        }
    }

    return sideNodes, err
}

// Generate a Merkle proof for a key.
func (smt *SparseMerkleTree) Prove(key []byte) ([][]byte, error) {
    sideNodes, err := smt.sideNodes(smt.digest(key))
    return sideNodes, err
}
