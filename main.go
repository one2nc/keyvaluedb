package main

import (
	"fmt"
	"keyvaluedb/domain"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Invalid arguments")
		return
	}

	argsWithoutProg := os.Args[1:]

	// Convert the []string slice to []interface{} slice
	args := make([]interface{}, len(argsWithoutProg[1:]))
	for i, arg := range argsWithoutProg[1:] {
		args[i] = arg
	}

	kvdb := domain.NewKeyValueDB()
	cmd := domain.NewCommand(argsWithoutProg[0], args...)

	fmt.Println(kvdb.Execute(cmd))
}
