package database

import (
	"sync"
)

type InMemoryDatabase struct {
	store map[string]map[string]interface{}
	mu    sync.RWMutex
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		store: make(map[string]map[string]interface{}),
	}
}

func (db *InMemoryDatabase) Create(location string, data map[string]interface{}) (bool, string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if the location already exists.
	if _, exists := db.store[location]; exists {
		return false, "Location already exists"
	}

	// Store the data.
	db.store[location] = data
	return true, ""
}

func (db *InMemoryDatabase) Read(location string) (bool, string, map[string]interface{}) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Check if the location exists.
	data, exists := db.store[location]
	if !exists {
		return false, "Location does not exist", nil
	}

	// Return a copy of the data to avoid modification.
	dataCopy := make(map[string]interface{})
	for k, v := range data {
		dataCopy[k] = v
	}

	return true, "", dataCopy
}

func (db *InMemoryDatabase) Update(location string, data map[string]interface{}) (bool, string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if the location exists.
	if _, exists := db.store[location]; !exists {
		return false, "Location does not exist"
	}

	// Update the data.
	db.store[location] = data
	return true, ""
}

func (db *InMemoryDatabase) Delete(location string) (bool, string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if the location exists.
	if _, exists := db.store[location]; !exists {
		return false, "Location does not exist"
	}

	// Delete the data.
	delete(db.store, location)
	return true, ""
}
