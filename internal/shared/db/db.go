package db

import (
	"log"
	"sync"
)

type IDB interface {
	Set(key string, value string) error
	Get(key string) (*string, error)
	Remove(key string) error
}

type DB struct {
	memoryDB map[string]string
	rw       *sync.RWMutex
}

func NewDB() IDB {
	return &DB{memoryDB: make(map[string]string), rw: &sync.RWMutex{}}
}

func (r *DB) Set(key string, value string) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	r.memoryDB[key] = value

	log.Printf("Adding key: %s, value: %s", key, value)

	return nil
}

func (r *DB) Get(key string) (*string, error) {
	log.Printf("Getting key: %s", key)

	r.rw.RLock()
	defer r.rw.RUnlock()

	value, ok := r.memoryDB[key]
	if ok {
		log.Printf("Found key: %s, value: %s", key, value)

		return &value, nil
	}

	log.Printf("Key not found: %s", key)

	return nil, ErrKeyNotFound
}

func (r *DB) Remove(key string) error {
	log.Printf("Removing key: %s", key)

	r.rw.Lock()
	defer r.rw.Unlock()

	if _, ok := r.memoryDB[key]; !ok {
		log.Printf("Key not found: %s", key)

		return ErrKeyNotFound
	}

	delete(r.memoryDB, key)

	return nil
}
