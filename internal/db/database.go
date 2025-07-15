// key logic for the database app
package db

import (
	"sync"
)

// DB is an in-memory key-value store with transaction stack support.
type DB struct {
	mu    sync.RWMutex
	stack []dbState // stack of full database states
}

type dbState struct {
	data     map[string]string
	valCount map[string]int
}

// New returns a new DB instance.
func New() *DB {
	initial := dbState{
		data:     make(map[string]string),
		valCount: make(map[string]int),
	}
	return &DB{
		stack: []dbState{initial},
	}
}

// Set assigns value to key.
func (db *DB) Set(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	cur := &db.stack[len(db.stack)-1]
	if old, ok := cur.data[key]; ok {
		cur.valCount[old]--
		if cur.valCount[old] == 0 {
			delete(cur.valCount, old)
		}
	}
	cur.data[key] = value
	cur.valCount[value]++
}

// Get - value for key, or false if not set.
func (db *DB) Get(key string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	cur := db.stack[len(db.stack)-1]
	v, ok := cur.data[key]
	return v, ok
}

// Unset key.
func (db *DB) Unset(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	cur := &db.stack[len(db.stack)-1]
	if v, ok := cur.data[key]; ok {
		delete(cur.data, key)
		cur.valCount[v]--
		if cur.valCount[v] == 0 {
			delete(cur.valCount, v)
		}
	}
}

// NumEqualTo - count of keys set to value.
func (db *DB) NumEqualTo(value string) int {
	db.mu.RLock()
	defer db.mu.RUnlock()

	cur := db.stack[len(db.stack)-1]
	return cur.valCount[value]
}

// Begin starts a transaction by pushing a deep copy of the current state.
func (db *DB) Begin() {
	db.mu.Lock()
	defer db.mu.Unlock()

	cur := db.stack[len(db.stack)-1]
	copyState := dbState{
		data:     deepCopyMap(cur.data),
		valCount: deepCopyIntMap(cur.valCount),
	}
	db.stack = append(db.stack, copyState)
}

// Rollback undoes the last transaction. Returns false if none.
func (db *DB) Rollback() bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	if len(db.stack) <= 1 {
		return false
	}
	db.stack = db.stack[:len(db.stack)-1]
	return true
}

// Commit makes the top transaction the new base. Returns false if none.
func (db *DB) Commit() bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	if len(db.stack) <= 1 {
		return false
	}
	top := db.stack[len(db.stack)-1]
	db.stack = db.stack[:1]
	db.stack[0] = top
	return true
}

// Helpers for deep copying maps
func deepCopyMap(src map[string]string) map[string]string {
	copy := make(map[string]string, len(src))
	for k, v := range src {
		copy[k] = v
	}
	return copy
}

func deepCopyIntMap(src map[string]int) map[string]int {
	copy := make(map[string]int, len(src))
	for k, v := range src {
		copy[k] = v
	}
	return copy
}
