package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"encoding/json"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v66/github"
)

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

// GetDiffBetweenCommits retrieves the list of files that have changed between two commits identified by their SHAs.
// It takes the repository path and the two commit SHAs as input parameters.
// Returns a slice of strings containing the names of the changed files and an error if any occurs.
func GetDiffBetweenCommits(repoPath, sha1, sha2 string) ([]string, error) {
	// Open the repository at the given path
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("could not open repository: %v", err)
	}

	// Get the commits corresponding to the given SHAs
	commit1, err := repo.CommitObject(plumbing.NewHash(sha1))
	if err != nil {
		return nil, fmt.Errorf("could not find commit for SHA %s: %v", sha1, err)
	}

	commit2, err := repo.CommitObject(plumbing.NewHash(sha2))
	if err != nil {
		return nil, fmt.Errorf("could not find commit for SHA %s: %v", sha2, err)
	}

	// Get the tree objects for both commits
	tree1, err := commit1.Tree()
	if err != nil {
		return nil, fmt.Errorf("could not get tree for commit %s: %v", sha1, err)
	}

	tree2, err := commit2.Tree()
	if err != nil {
		return nil, fmt.Errorf("could not get tree for commit %s: %v", sha2, err)
	}

	// Get the diff between the two trees
	changes, err := object.DiffTree(tree1, tree2)
	if err != nil {
		return nil, fmt.Errorf("could not get diff between trees: %v", err)
	}

	// Collect the filenames of changed files
	var diffFiles []string
	for _, change := range changes {
		// Append the file name (NewName) to the list of diff files
		diffFiles = append(diffFiles, change.To.Name)
	}

	return diffFiles, nil
}

// filterStrings filters the input strings based on inclusion and exclusion regex patterns.
func FilterStrings(input []string, includePatterns, excludePatterns []string) []string {
	var result []string

	// Compile inclusion regex patterns
	var includeRegexes []*regexp.Regexp
	for _, pattern := range includePatterns {
		if pattern != "" {
			valid, err := isValidRegex(pattern)
			if valid {
				re, _ := regexp.Compile(pattern)
				includeRegexes = append(includeRegexes, re)
			} else {
				fmt.Printf("The included pattern '%s' is an illegal regex: %v\n", pattern, err)
			}
		}
	}

	// Compile exclusion regex patterns
	var excludeRegexes []*regexp.Regexp
	for _, pattern := range excludePatterns {
		if pattern != "" {
			valid, err := isValidRegex(pattern)
			if valid {
				re, _ := regexp.Compile(pattern)
				excludeRegexes = append(excludeRegexes, re)
			} else {
				fmt.Printf("The excluded pattern '%s' is an illegal regex: %v\n", pattern, err)
			}
		}
	}

	// Filter the input strings
	for _, str := range input {
		included := false
		for _, re := range includeRegexes {
			if re != nil && re.MatchString(str) {
				included = true
				break
			}
		}

		if included {
			excluded := false
			for _, re := range excludeRegexes {
				if re.MatchString(str) {
					excluded = true
					break
				}
			}

			if !excluded {
				result = append(result, str)
			}
		}
	}

	return result
}

// SetGitHubOutput sets a GitHub Actions output variable.
func SetGitHubOutput(name, value string) {
	// Get the GITHUB_OUTPUT environment variable
	outputFile := os.Getenv("GITHUB_OUTPUT")
	if outputFile == "" {
		log.Println("GITHUB_OUTPUT is not set")
		return
	}

	// Write the output to the environment file
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Printf("failed to open output file: %v", err)
		return
	}
	defer f.Close()

	// Format the output and write it to the file
	if _, err := f.WriteString(fmt.Sprintf("%s=%s\n", name, value)); err != nil {
		log.Printf("failed to write output: %v", err)
	}
}

// Delta calculates the difference between the latest successful deployment SHA and the current SHA,
// and sets GitHub Actions output variables with the results. If there are any changes detected,
// the "is_detected" output is set to "true" and the "delta_files" output is set to a JSON-encoded
// list of the changed files. If there are no changes, the "is_detected" output is set to "false".
func Delta(repoPath string) {
	cfg := GetInputConfig()
	cfg.Validate()
	client := GetClient(&cfg)

	var baseSha string
	if cfg.Environment != "" {
		baseSha = GetLatestSuccessfulDeploymentSha(client, &cfg)
	} else {
		baseSha = GetBranchLatestSHA(client, &cfg)
	}

	diffs, err := GetDiffBetweenCommits(repoPath, baseSha, cfg.Sha)
	if err != nil {
		log.Fatalf("Error getting diff between commits: %v", err)
	}

	deltas := FilterStrings(diffs, cfg.IncludesPatterns, cfg.ExcludesPatterns)

	if len(deltas) > 0 {
		SetGitHubOutput("is_detected", "true")
		jsonData, err := json.Marshal(deltas)
		if err != nil {
			log.Printf("Error marshalling to JSON: %v", err)
		}
		SetGitHubOutput("delta_files", string(jsonData))
	} else {
		SetGitHubOutput("is_detected", "false")
	}
}
