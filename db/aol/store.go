package aol

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"

	"github.com/msam1r/kallbaz-db/db"
	"github.com/msam1r/kallbaz-db/db/record"
)

type Store struct {
	db.Store
	storagePath   string
	maxRecordSize int
	logger        *log.Logger
	async         bool
	writeMutex    sync.Mutex
}

type Config struct {
	BasePath      string
	MaxRecordSize *int
	Async         *bool
	Logger        *log.Logger
}

const (
	logFile              = "store.db"
	defaultMaxRecordSize = 1024 * 1024
	defaultAsync         = false
)

var voidLogger = log.New(ioutil.Discard, "", log.LstdFlags)

// NewStore - returns a new disk store.
func NewStore(config Config) (db.Store, error) {
	var (
		maxRecordSize = defaultMaxRecordSize
		storagePath   = path.Join(config.BasePath, logFile)
		async         = defaultAsync
		logger        = voidLogger
	)

	if _, err := os.OpenFile(storagePath, os.O_CREATE, 0600); err != nil {
		return nil, err
	}

	if config.MaxRecordSize != nil {
		maxRecordSize = *config.MaxRecordSize
	}

	if config.Async != nil {
		async = *config.Async
	}

	if config.Logger != nil {
		logger = config.Logger
	}

	return &Store{
		storagePath:   storagePath,
		maxRecordSize: maxRecordSize,
		logger:        logger,
		async:         async,
	}, nil
}

// Get - returns the value of the given key, or a db.NotFoundError if
// the key was not found.
func (s *Store) Get(key string) ([]byte, error) {
	file, err := os.Open(s.storagePath)
	defer file.Close()

	if err != nil {
		return nil, fmt.Errorf("could not open file: %s, %w", s.storagePath, err)
	}

	scanner, err := record.NewScanner(file, s.maxRecordSize)
	if err != nil {
		return nil, fmt.Errorf("could not create scanner for file: %s, %w", s.storagePath, err)
	}

	var found *record.Record
	for scanner.Scan() {
		record := scanner.Record()
		if record.Key() == key {
			found = record
		}
	}

	if scanner.Err() != nil {
		s.logger.Printf("error encountered: %s", scanner.Err())
		return nil, scanner.Err()
	}

	if found == nil || found.IsTombstone() {
		return nil, db.NewNotFoundError(key)
	}

	return found.Value(), nil
}

// Put - save the value to the database, and returns an error occured.
func (s *Store) Put(key string, value []byte) error {
	record := record.NewValue(key, value)
	return s.append(record)
}

// Delete - deletes the value from the database and returns occured error.
func (s *Store) Delete(key string) error {
	record := record.NewTombstone(key)
	return s.append(record)
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

func (s *Store) append(record *record.Record) error {
	file, err := os.OpenFile(s.storagePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	defer file.Close()

	if err != nil {
		return fmt.Errorf("could not open file: %s for write, %w", s.storagePath, err)
	}

	_, err = record.Write(file)

	if err != nil {
		return fmt.Errorf("could not write record to file: %s, %w", s.storagePath, err)
	}

	return nil
}
