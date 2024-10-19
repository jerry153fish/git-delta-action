package main

import (
	"log"

	"github.com/jerry153fish/git-delta/internal"
)

func main() {
	// Get the delta configuration
	// c := internal.GetInputConfig()
	// Print environments from os
	// for _, env := range os.Environ() {
	// 	log.Println(env)
	// }

	// Log the Delta configuration
	// log.Printf("Input Config: %+v\n", c)
	// Create a new GitHub client with authentication

	// client := internal.GetClient(&c)

	// Log the creation of the GitHub client
	log.Println("GitHub client created with authentication")

	// internal.GetLatestSuccessfulDeploymentSha(client, &c)
	// internal.GetDiffBetweenCommits(client, &c)
	sha1 := "c6023e778dac2c67e7ec0c42889e349a76414294"
	sha2 := "7bf5f383a901d5dda65adedfb351fa6f9fffd4f2"

	result, err := internal.GetDiffBetweenCommits(".", sha1, sha2)
	if err != nil {
		log.Fatalf("Error getting diff between commits: %v", err)
	}

	for _, file := range result {
		log.Println(file) // Changed from t.Logf to t.Log
	}
}
