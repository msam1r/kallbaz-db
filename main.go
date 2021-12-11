package main

import (
	"fmt"
	"log"

	"github.com/msam1r/kallbaz-db/db/memory"
)

func main() {
	fmt.Println("Welcome to KallbazDB.")

	kdb := memory.NewStore(memory.Config{
		MaxRecordSize: 1024,
		Logger:        &log.Logger{},
	})

	val, err := kdb.Get("name")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(val))
}
