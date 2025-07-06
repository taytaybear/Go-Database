// key logic for the database app
package db

import (
	"sync"
)

// DB is an in-memory key-value store with transaction support.
type DB struct {
	mu       sync.RWMutex
	data     map[string]string
	valCount map[string]int
	txns     []*Tx
}

// Tx - a transaction block.
type Tx struct {
	data     map[string]string
	valCount map[string]int
}

// New returns a new DB instance.
func New() *DB {
	return &DB{
		data:     make(map[string]string),
		valCount: make(map[string]int),
		txns:     make([]*Tx, 0),
	}
}

// Set assigns value to key.
func (db *DB) Set(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if len(db.txns) > 0 {
		tx := db.txns[len(db.txns)-1]
		db.setInTx(tx, key, value)
		return
	}

	if old, ok := db.data[key]; ok {
		db.valCount[old]--
		if db.valCount[old] == 0 {
			delete(db.valCount, old)
		}
	}
	db.data[key] = value
	db.valCount[value]++
}

// Get - value for key, or false if not set.
func (db *DB) Get(key string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if len(db.txns) > 0 {
		for i := len(db.txns) - 1; i >= 0; i-- {
			tx := db.txns[i]
			if v, ok := tx.data[key]; ok {
				if v == "" {
					return "", false
				}
				return v, true
			}
		}
	}
	v, ok := db.data[key]
	return v, ok
}

// Unset key.
func (db *DB) Unset(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if len(db.txns) > 0 {
		tx := db.txns[len(db.txns)-1]
		db.unsetInTx(tx, key)
		return
	}

	if v, ok := db.data[key]; ok {
		delete(db.data, key)
		db.valCount[v]--
		if db.valCount[v] == 0 {
			delete(db.valCount, v)
		}
	}
}

// NumEqualTo - count of keys set to value.
func (db *DB) NumEqualTo(value string) int {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if len(db.txns) > 0 {
		return db.numEqualToTx(value)
	}
	return db.valCount[value]
}

// Begin starts a transaction.
func (db *DB) Begin() {
	db.mu.Lock()
	defer db.mu.Unlock()

	tx := &Tx{
		data:     make(map[string]string),
		valCount: make(map[string]int),
	}
	db.txns = append(db.txns, tx)
}

// Rollback undoes the last tx. false if none.
func (db *DB) Rollback() bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	if len(db.txns) == 0 {
		return false
	}
	db.txns = db.txns[:len(db.txns)-1]
	return true
}

// applies all tx. false if none.
func (db *DB) Commit() bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	if len(db.txns) == 0 {
		return false
	}
	for _, tx := range db.txns {
		for k, v := range tx.data {
			if v == "" {
				if old, ok := db.data[k]; ok {
					delete(db.data, k)
					db.valCount[old]--
					if db.valCount[old] == 0 {
						delete(db.valCount, old)
					}
				}
			} else {
				if old, ok := db.data[k]; ok {
					db.valCount[old]--
					if db.valCount[old] == 0 {
						delete(db.valCount, old)
					}
				}
				db.data[k] = v
				db.valCount[v]++
			}
		}
	}
	db.txns = make([]*Tx, 0)
	return true
}

// helpers for tx logic
func (db *DB) setInTx(tx *Tx, key, value string) {
	cur := db.getCurrentTxValue(key)
	if cur != "" {
		tx.valCount[cur]--
		if tx.valCount[cur] == 0 {
			delete(tx.valCount, cur)
		}
	}
	tx.data[key] = value
	tx.valCount[value]++
}

func (db *DB) unsetInTx(tx *Tx, key string) {
	cur := db.getCurrentTxValue(key)
	if cur != "" {
		tx.data[key] = ""
		tx.valCount[cur]--
		if tx.valCount[cur] == 0 {
			delete(tx.valCount, cur)
		}
	}
}

func (db *DB) getCurrentTxValue(key string) string {
	for i := len(db.txns) - 1; i >= 0; i-- {
		tx := db.txns[i]
		if v, ok := tx.data[key]; ok {
			if v == "" {
				return ""
			}
			return v
		}
	}
	return db.data[key]
}

func (db *DB) numEqualToTx(value string) int {
	count := db.valCount[value]
	for _, tx := range db.txns {
		if d, ok := tx.valCount[value]; ok {
			count += d
		}
	}
	return count
}
