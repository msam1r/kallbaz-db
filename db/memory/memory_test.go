package memory

import (
	"log"
	"testing"

	"github.com/msam1r/kallbaz-db/db"
)

var (
	key   = "name"
	value = []byte("mohamed samir")
)

func TestStore(t *testing.T) {
	kdb := createDatabase()
	err := kdb.Put(key, []byte(value))
	errCheck(err, t)

	val, err := kdb.Get(key)
	errCheck(err, t)

	if string(val) != string(value) {
		t.Fatalf("Expected: %s, got: %s", string(value), string(val))
	}
}

func createDatabase() db.Store {
	return NewStore(Config{
		MaxRecordSize: 1024,
		Logger:        &log.Logger{},
	})
}

func errCheck(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("Expected no errors got: %v", err)
	}
}
