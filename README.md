# smt

A Go library that implements a Sparse Merkle tree for a key-value map. The tree implements the same optimisations specified in the [Libra whitepaper][libra whitepaper], to reduce the number of hash operations required per tree operation to O(k) where k is the number of non-empty elements in the tree.

[![Tests](https://github.com/celestiaorg/smt/actions/workflows/test.yml/badge.svg)](https://github.com/celestiaorg/smt/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/celestiaorg/smt/branch/master/graph/badge.svg?token=U3GGEDSA94)](https://codecov.io/gh/celestiaorg/smt)
[![GoDoc](https://godoc.org/github.com/celestiaorg/smt?status.svg)](https://godoc.org/github.com/celestiaorg/smt)

## Example

```go
package main

import (
	"crypto/sha256"
	"fmt"

	"github.com/celestiaorg/smt"
)

func main() {
	// Initialise two new key-value store to store the nodes and values of the tree
	nodeStore := smt.NewSimpleMap()
	valueStore := smt.NewSimpleMap()
	// Initialise the tree
	tree := smt.NewSparseMerkleTree(nodeStore, valueStore, sha256.New())

	// Update the key "foo" with the value "bar"
	_, _ = tree.Update([]byte("foo"), []byte("bar"))

	// Generate a Merkle proof for foo=bar
	proof, _ := tree.Prove([]byte("foo"))
	root := tree.Root() // We also need the current tree root for the proof

	// Verify the Merkle proof for foo=bar
	if smt.VerifyProof(proof, root, []byte("foo"), []byte("bar"), sha256.New()) {
		fmt.Println("Proof verification succeeded.")
	} else {
		fmt.Println("Proof verification failed.")
	}
}
```

[libra whitepaper]: https://diem-developers-components.netlify.app/papers/the-diem-blockchain/2020-05-26.pdf

## Patch log

> v0.2.1

**Fix**

* Multiversion: modifies the value's valuehash to support multiversion.
* Same Value: same value with same hash (recognized as key in `valuestore`) will cause over deletion - a deletion of a kv may cause another key with same value cannot find the value. Using `[keyhash,valuehash]` as a value's valuehash can mitigate this problem.
* Get Exception: adds judgment when a key's path is not equal to the keyhash in the leaf (i.e., requests a non-existent key, instead of returning a false value, it returns a empty value)

**Modify**

* Mapstore Interface: fit to badger store.
* Related call.

**Supplement**

* Multithread Support: replaces the public hasher with the instance (to avoid race condition).
* Multiversion Same Leaf Support: simple value store adds `count` FYI (avoid the deletion of same kv, old version leaf to affect the newer one).
* Removing Intermidiate Version: `RemovePath` provides a parameter `keepRoot` to retain old version root, e.g., a state transition `a -> b -> c`, `RemovePath` can remove state `b` without affecting state `a`.
