package main

import (
	"fmt"
	"github.com/google/uuid"
)

func main() {
	fmt.Println("Hello, World!")

	// Generate a random UUID (version 4)
	id := uuid.New()
	fmt.Println("Generated UUID:", id.String())
}