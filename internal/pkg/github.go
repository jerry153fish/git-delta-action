package internal

import (
	"context"
	"fmt"
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

// extractOwnerRepo takes a repository string in the format "owner/repo"
// and splits it into the owner and repository name.
func extractOwnerRepo(repo string) (owner, repoName string) {
	parts := strings.Split(repo, "/")
	if len(parts) == 0 {
		return "", ""
	}

	if len(parts) == 1 {
		return "", parts[0]
	}

	repoName = parts[len(parts)-1]
	owner = strings.Join(parts[:len(parts)-1], "/")
	return owner, repoName
}

// GetClient creates a new GitHub client with the provided configuration.
func GetClient(c *InputConfig) *github.Client {
	// Create a new GitHub client with no initial HTTP client.
	// Authenticate the client using the provided GitHub token.
	return github.NewClient(nil).WithAuthToken(c.GithubToken)
}

// CompareSHAs compares the commits between the base SHA and the current SHA for the
// specified repository, and returns a list of the filenames that have changed.
//
// client is the GitHub API client to use for the comparison.
// cfg is the input configuration containing the repository information and the current SHA.
// baseSHA is the base SHA to compare against.
//
// Returns a slice of filenames that have changed between the base SHA and the current SHA.
// If an error occurs during the comparison, an error is returned.
func CompareSHAs(client *github.Client, cfg *InputConfig, baseSHA string) ([]string, error) {
	// Create a background context for the GitHub API calls
	ctx := context.Background()
	// Extract the owner and repository names from the full repository path
	owner, repo := extractOwnerRepo(cfg.Repo)

	// Compare the commits between the base SHA and the current SHA
	comparison, _, err := client.Repositories.CompareCommits(ctx, owner, repo, baseSHA, cfg.Sha, nil)
	if err != nil {
		return nil, fmt.Errorf("error comparing commits: %v", err)
	}

	// Extract the filenames from the comparison
	var fileNames []string
	for _, file := range comparison.Files {
		fileNames = append(fileNames, file.GetFilename())
	}

	return fileNames, nil
}
