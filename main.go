package main

import (
	"log"

	"github.com/jerry153fish/git-delta/internal"
)

func main() {
	// Get the delta configuration
	c := internal.GetInputConfig()
	// Print environments from os
	// for _, env := range os.Environ() {
	// 	log.Println(env)
	// }

	// Log the Delta configuration
	log.Printf("Input Config: %+v\n", c)
	// Create a new GitHub client with authentication

	client := internal.GetClient(&c)

	// Log the creation of the GitHub client
	log.Println("GitHub client created with authentication")

	internal.GetLatestSuccessfulDeploymentSha(client, &c)
	internal.GetBranchLatestSHA(client, &c)
}
