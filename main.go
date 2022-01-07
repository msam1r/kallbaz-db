package main

import (
	"fmt"

	"github.com/msam1r/kallbaz-db/db/record"
)

func main() {
	fmt.Println("Welcome to KallbazDB.")

	val := record.NewValue("name", []byte("mohamed"))

	a, _ := record.FromBytes(val.ToBytes())

	fmt.Println(string(a.Value()))
}
