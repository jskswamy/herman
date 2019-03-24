package db

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

// Open open the database at the specified location and return error if any
func Open(path string) error {
	instance, err := bolt.Open(path, 0666, nil)

	if err != nil {
		return fmt.Errorf("unable to open db, reason %v", err)
	}
	db = instance
	return nil
}

// Close the db and return error if any
func Close() error {
	if db == nil {
		return fmt.Errorf("db already closed")
	}
	defer func() {
		db = nil
	}()
	return db.Close()
}
