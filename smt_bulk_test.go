package smt

import (
	"bytes"
	"time"

	// "fmt"
	"crypto/sha256"
	"math/rand"
	"reflect"
	"testing"
)

// Test all block-ops in bulk.
func TestSparseMerkleTree1(t *testing.T) {
	for i := 0; i < 5; i++ {
		// Test more inserts/updates than deletions.
		bulkOperations1(t, 100, 200, 200, 50)
	}
	for i := 0; i < 5; i++ {
		// Test extreme deletions.
		bulkOperations1(t, 200, 100, 100, 500)
	}
}

// Test all block-ops in bulk, with specified ratio probabilities of insert, update and delete.
func bulkOperations1(t *testing.T, blocks int, insert int, update int, del int) {
	smn, smv := NewSimpleMap(), NewSimpleMap()
	smt := NewSparseMerkleTree(smn, smv, sha256.New())

	max := insert + update + del
	// record every version's kv
	kv := make([]map[string]string, blocks)
	for i := 0; i < blocks; i++ {
		kv[i] = make(map[string]string)
	}
	roots := make([][]byte, blocks)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < blocks; i++ {
		if i != 0 {
			for k, v := range kv[i-1] {
				kv[i][k] = v
			}
		}
		for j := 0; j < 10; j++ {
			n := rand.Intn(max)
			oldRoot := smt.Root()
			if n < insert { // Insert
				keyLen := 16 + rand.Intn(32)
				key := make([]byte, keyLen)
				rand.Read(key)

				valLen := 1 + rand.Intn(64)
				val := make([]byte, valLen)
				rand.Read(val)
				// val := []byte(fmt.Sprintf("val%d", i*10+j))

				kv[i][string(key)] = string(val)
				_, err := smt.Update(key, val)

				if j != 0 {
					if i > 0 {
						err = smt.RemovePath(key, oldRoot, roots[i-1])
					} else {
						err = smt.RemovePathForRoot(key, oldRoot)
					}
				}
				if err != nil {
					t.Errorf("error: %v", err)
				}
			} else if n > insert && n < insert+update { // Update
				keys := reflect.ValueOf(kv[i]).MapKeys()
				if len(keys) == 0 {
					continue
				}
				key := []byte(keys[rand.Intn(len(keys))].Interface().(string))

				valLen := 1 + rand.Intn(64)
				val := make([]byte, valLen)
				rand.Read(val)
				// val := []byte(fmt.Sprintf("val%d", i*10+j))

				kv[i][string(key)] = string(val)
				_, err := smt.Update(key, val)

				if j != 0 {
					if i > 0 {
						err = smt.RemovePath(key, oldRoot, roots[i-1])
					} else {
						err = smt.RemovePathForRoot(key, oldRoot)
					}
				}
				if err != nil {
					t.Errorf("error: %v", err)
				}
			} else { // Delete
				keys := reflect.ValueOf(kv[i]).MapKeys()
				if len(keys) == 0 {
					continue
				}
				key := []byte(keys[rand.Intn(len(keys))].Interface().(string))

				delete(kv[i], string(key))
				_, err := smt.Delete(key)

				if j != 0 {
					if i > 0 {
						err = smt.RemovePath(key, oldRoot, roots[i-1])
					} else {
						err = smt.RemovePathForRoot(key, oldRoot)
					}
					// fmt.Printf("Remove from root[%x]\n", oldRoot)
				}
				if err != nil {
					t.Errorf("error: %v", err)
				}
			}

		}
		roots[i] = smt.Root()
		checkOne(t, smt, &kv[i])
	}

	checkAll(t, smt, roots, kv)
}

func checkOne(t *testing.T, smt *SparseMerkleTree, kv *map[string]string) {
	for k, v := range *kv {
		actualVal, err := smt.Get([]byte(k))
		if err != nil {
			t.Errorf("error: %v", err)
			continue
			// return false;
		}

		if !bytes.Equal([]byte(v), actualVal) {
			t.Error("got incorrect value when bulk testing operations")
		}
	}
}

func checkAll(t *testing.T, smt *SparseMerkleTree, roots [][]byte, kv []map[string]string) {
	for i, root := range roots {
		for k, v := range kv[i] {
			actualVal, err := smt.GetFromRoot([]byte(k), root)
			if err != nil {
				t.Errorf("error: %v", err)
			}

			if !bytes.Equal([]byte(v), actualVal) {
				t.Error("got incorrect value when bulk testing operations")
			}
		}
	}
}
