package internal

import (
	"context"
	"log"
	"strings"

	"github.com/google/go-github/v66/github"
)

// GetLatestSuccessfulDeploymentSha retrieves the Sha of latest successful deployment for a given environment
func GetLatestSuccessfulDeploymentSha(client *github.Client, cfg *InputConfig) string {
	// Create a background context for the GitHub API calls
	ctx := context.Background()
	// Extract the owner and repository names from the full repository path
	owner, repo := extractOwnerRepo(cfg.Repo)

	// Config the options to setup the environment and page size
	opt := &github.DeploymentsListOptions{
		Environment: cfg.Environment,
		ListOptions: github.ListOptions{PerPage: 50},
	}

	// Find the recently successfully deployment by loop on the pagination
	for {
		// List all deployment by page
		deployments, resp, err := client.Repositories.ListDeployments(ctx, owner, repo, opt)
		if err != nil {
			log.Fatalf("Error listing deployments: %v", err)
		}
		// Check if the deployment has a successful state
		for _, deployment := range deployments {

			// Get the statues for the deployment
			statuses, _, err := client.Repositories.ListDeploymentStatuses(ctx, owner, repo, deployment.GetID(), &github.ListOptions{PerPage: 1})
			if err != nil {
				log.Printf("Error getting deployment status: %v", err)
				continue
			}

			// move the next deployment if no statues found
			if len(statuses) == 0 {
				log.Printf("No deployment status found for deployment ID %d", deployment.GetID())
				continue
			}

			// check the status
			status := statuses[0]
			if status.GetState() == "success" {
				log.Printf("Latest successful deployment for %s: ID %d, SHA %s", cfg.Environment, deployment.GetID(), deployment.GetSHA())
				return deployment.GetSHA()
			}
		}

		// loop to next page
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	log.Printf("No successful deployments found for environment: %s", cfg.Environment)
	return ""
}

// GetLatestSHA retrieves the latest commit SHA for a specified branch in a repository.
func GetBranchLatestSHA(client *github.Client, cfg *InputConfig) string {
	// Create a background context for the GitHub API calls
	ctx := context.Background()
	// Extract the owner and repository names from the full repository path
	owner, repo := extractOwnerRepo(cfg.Repo)

	// Get the reference for the specified branch
	ref, _, err := client.Git.GetRef(ctx, owner, repo, "refs/heads/"+cfg.Branch)
	if err != nil {
		log.Printf("Error retrieving SHA for branch '%s' in repository '%s': %v", cfg.Branch, cfg.Repo, err)
		return ""
	}
	log.Printf("Latest successful Sha for Brach %s, SHA %s", cfg.Branch, ref.Object.GetSHA())
	// Return the SHA of the latest commit
	return ref.Object.GetSHA()
}

// extractOwnerRepo takes a repository string in the format "owner/repo"
// and splits it into the owner and repository name.
func extractOwnerRepo(repo string) (string, string) {
	// Split the repo string by the "/" character to separate owner and repo name
	parts := strings.Split(repo, "/")
	// Return the owner (second to last part) and repo name (last part)
	return parts[len(parts)-2], parts[len(parts)-1]
}

// GetClient creates a new GitHub client with the provided configuration.
func GetClient(c *InputConfig) *github.Client {
	// Create a new GitHub client with no initial HTTP client.
	// Authenticate the client using the provided GitHub token.
	return github.NewClient(nil).WithAuthToken(c.GithubToken)
}
