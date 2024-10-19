package main

import (
	"log"

	"github.com/jerry153fish/git-delta/internal"
)

func main() {
	// Get the delta configuration
	c := internal.GetInputConfig()

	// Log the Delta configuration
	log.Printf("Delta Config: %+v\n", c)
}
