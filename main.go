package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Hello, World!")
	fmt.Println(os.Environ())
	// Define the command and its arguments
	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	// Print the names of the files
	fmt.Println("Files in current directory:")
	for _, file := range files {
		fmt.Println(file.Name())
	}
}
