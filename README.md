# Kallbaz DB
Simple Distributed key-value store (in-memory/disk) written with Golang.

# Installation
```
go get github.com/msam1r/kallbaz-db
```

# Usage

### API

```go
// Get - returns the value of the given key or any error occurs. if the
// key was not found it will return a NotFoundError.
Get(key string) ([]byte, error)
// Put - stores the value. it will return a BadRequestError if the provided
// data was invalid or any other error occured.
Put(key string, value []byte) error
// Delete - deletes the value of the given key.
Delete(key string) error
// Close - closes the database and returns when all internal processes
// has stopped. it returns any error occurs.
Close() error
// IsNotFoundError - check if the error is of type NotFoundError.
IsNotFoundError(err error) bool
// IsBadRequestError - check if the error is of type BadRequestError.
IsBadRequestError(err error) bool
```

### In-memory db

```go
package main

import (
    "fmt"
    "github.com/msam1r/kallbaz-db/db/memory"
)

func main() {
    kdb := memory.NewStore(memory.Config{
		MaxRecordSize: 1024,
		Logger:        &log.Logger{},
	})

    // Put value to the database
    err := kdb.Put("name", []byte("mohamed samir"))
    if err != nil {
		fmt.Println(err.Error())
	}

    // Get value from the database.
    var key string = "name-12"
    v, err := kdb.Get(key)
	if err != nil && kdb.IsNotFoundError(err) {
        fmt.Printf("Key: %s does not exists", key)
	}
    fmt.Println(string(v))

    // Delete the value from the database.
    err := kdb.Delete("name")
    if err != nil && (kdb.IsNotFoundError(err) || kdb.IsBadRequestError(err)) {
		fmt.Println(err.Error())
	}
}
```

# TODO
- [x] in memory store
- [ ] disk store {`in progress`}
- [ ] applies aol (append-only log) structure {`in progress`}
- [ ] binary storage format {`in progress`}
- [ ] simple index implementation {`not started`}
- [ ] sharding {`not started`}
- [ ] replication {`not started`}
