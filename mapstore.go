package smt

import (
	"fmt"
)

// MapStore is a key-value store.
type MapStore interface {
	Put(key []byte, value []byte) error // Put updates the value for a key.
	Get(key []byte) ([]byte, error)     // Get gets the value for a key.
	Has(key []byte) (bool, error)
	Delete(key []byte) error // Delete deletes a key.
	Close() error
}

// InvalidKeyError is thrown when a key that does not exist is being accessed.
type InvalidKeyError struct {
	Key []byte
}

func (e *InvalidKeyError) Error() string {
	return fmt.Sprintf("invalid key: %x", e.Key)
}

type SimpleValue struct {
	data  []byte
	count uint32
}

// SimpleMap is a simple in-memory map.
type SimpleMap struct {
	m map[string]SimpleValue
}

// NewSimpleMap creates a new empty SimpleMap.
func NewSimpleMap() *SimpleMap {
	return &SimpleMap{
		m: make(map[string]SimpleValue),
	}
}

// Get gets the value for a key.
func (sm *SimpleMap) Get(key []byte) ([]byte, error) {
	if value, ok := sm.m[string(key)]; ok {
		return value.data, nil
	}
	return nil, &InvalidKeyError{Key: key}
}

// Put updates the value for a key.
func (sm *SimpleMap) Put(key []byte, value []byte) error {
	if data, ok := sm.m[string(key)]; ok {
		sm.m[string(key)] = SimpleValue{
			data:  value,
			count: data.count + 1,
		}
	} else {
		sm.m[string(key)] = SimpleValue{
			data:  value,
			count: 1,
		}
	}
	return nil
}

func (sm *SimpleMap) Has(key []byte) (bool, error) {
	if _, ok := sm.m[string(key)]; ok {
		return true, nil
	}
	return false, nil
}

// Delete deletes a key.
func (sm *SimpleMap) Delete(key []byte) error {
	data, ok := sm.m[string(key)]
	if ok {
		data.count -= 1
		if data.count == 0 {
			delete(sm.m, string(key))
		} else {
			sm.m[string(key)] = data
		}
		return nil
	}
	return &InvalidKeyError{Key: key}
}

func (sm *SimpleMap) Size() int64 {
	return int64(len(sm.m))
}

func (sm *SimpleMap) Close() error {
	sm.m = nil
	return nil
}
