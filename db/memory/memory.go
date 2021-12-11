package memory

import (
	"fmt"
	"log"
	"sync"

	"github.com/msam1r/kallbaz-db/db"
)

// Store - a simple key-value store which keeps the state in the memory.
// It has the performance of go's map structure.
type Store struct {
	db.Store
	maxRecordSize int
	logger        *log.Logger

	sync.RWMutex
	table map[string][]byte
}

// Config - contains the configuration properties for the in memory store.
type Config struct {
	MaxRecordSize int
	Logger        *log.Logger
}

// NewStore - returns a new memory store.
func NewStore(config Config) db.Store {
	return &Store{
		maxRecordSize: config.MaxRecordSize,
		logger:        config.Logger,
		table:         map[string][]byte{},
	}
}

// Get - returns the value of the given key, or a db.NotFoundError if
// the key was not found.
func (s *Store) Get(key string) ([]byte, error) {
	s.RLock()
	val, ok := s.table[key]
	s.RUnlock()

	if !ok {
		return nil, db.NewNotFoundError(key)
	}

	return val, nil
}

// Put - save the value to the database, and returns an error occured.
func (s *Store) Put(key string, value []byte) error {
	size := len([]byte(key)) + len(value)
	if size > s.maxRecordSize {
		return db.NewBadRequestError(
			fmt.Sprintf("key-value too big, max size allowd: %d", s.maxRecordSize),
		)
	}

	s.Lock()
	s.table[key] = value
	s.Unlock()

	return nil
}

// Delete - deletes the value from the database and returns occured error.
func (s *Store) Delete(key string) error {
	s.Lock()
	delete(s.table, key)
	s.Unlock()

	return nil
}

// Close - closes the database.
func (s *Store) Close() error {
	s.logger.Print("Closing the database")
	return nil
}

// IsNotFoundError - check if the error is of type db.NotFoundError
func (s *Store) IsNotFoundError(err error) bool {
	return db.IsNotFoundError(err)
}

// IsBadRequestError - check if the error is of type db.BadRequestError
func (s *Store) IsBadRequestError(err error) bool {
	return db.IsBadRequestError(err)
}
