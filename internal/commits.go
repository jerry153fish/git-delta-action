package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"encoding/json"

	"github.com/caarlos0/env/v11"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v66/github"
)

const (
	// FileSeparator is used to split the Files and IgnoreFiles strings into slices
	FileSeparator = "\n"
)

// InputConfig holds the configuration for the Action Inputs
type InputConfig struct {
	Environment          string `env:"INPUT_ENVIRONMENT"`
	Commit               string `env:"INPUT_COMMIT"`
	Files                string `env:"INPUT_FILES"`
	IgnoreFiles          string `env:"INPUT_IGNORE_FILES"`
	DeltaOutputPathDepth string `env:"INPUT_DELTA_OUTPUT_PATH_DEPTH"`
	GithubToken          string `env:"INPUT_GITHUB_TOKEN"`
	Sha                  string `env:"GITHUB_SHA"`
	Ref                  string `env:"GITHUB_REF"`
	ApiUrl               string `env:"GITHUB_API_URL"`
	Workflow             string `env:"GITHUB_WORKFLOW"`
	EventName            string `env:"GITHUB_EVENT_NAME"`
	Job                  string `env:"GITHUB_JOB"`
	Repo                 string `env:"GITHUB_REPOSITORY"`
	Branch               string `env:"INPUT_BRANCH"`
	FilePatterns         []string
	IgnoreFilePatterns   []string
}

// GetInputConfig parses environment variables into an InputConfig struct
// and processes the Files and IgnoreFiles fields
func GetInputConfig() InputConfig {
	// Parse environment variables into InputConfig struct
	c, err := env.ParseAs[InputConfig]()
	if err != nil {
		log.Fatalf("Failed to parse InputConfig: %v", err)
	}

	// If Files is not empty, split it into FilePatterns
	if c.Files != "" {
		c.FilePatterns = strings.Split(c.Files, FileSeparator)
	}

	// If IgnoreFiles is not empty, split it into IgnoreFilePatterns
	if c.IgnoreFiles != "" {
		c.IgnoreFilePatterns = strings.Split(c.IgnoreFiles, FileSeparator)
	}

	return c
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

// isValidRegex checks if the provided string is a valid regular expression.
func isValidRegex(pattern string) (bool, error) {
	_, err := regexp.Compile(pattern)
	if err != nil {
		return false, err // Return false and the error
	}
	return true, nil // Return true if no error
}

// SetGitHubOutput sets a GitHub Actions output variable.
func SetGitHubOutput(name, value string) error {
	// Get the GITHUB_OUTPUT environment variable
	outputFile := os.Getenv("GITHUB_OUTPUT")
	if outputFile == "" {
		return fmt.Errorf("GITHUB_OUTPUT is not set")
	}

	// Write the output to the environment file
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open output file: %v", err)
	}
	defer f.Close()

	// Format the output and write it to the file
	if _, err := f.WriteString(fmt.Sprintf("%s=%s\n", name, value)); err != nil {
		return fmt.Errorf("failed to write output: %v", err)
	}

	return nil
}

func Delta() {
	cfg := GetInputConfig()
	client := GetClient(&cfg)

	var baseSha string
	if cfg.Environment != "" {
		baseSha = GetLatestSuccessfulDeploymentSha(client, &cfg)
	} else {
		baseSha = GetBranchLatestSHA(client, &cfg)
	}

	diffs, err := GetDiffBetweenCommits(".", baseSha, cfg.Sha)
	if err != nil {
		log.Fatalf("Error getting diff between commits: %v", err)
	}

	log.Println(diffs)

	deltas := FilterStrings(diffs, cfg.FilePatterns, cfg.IgnoreFilePatterns)

	log.Println(deltas)

	if len(deltas) > 0 {
		SetGitHubOutput("is_detected", "true")
		jsonData, err := json.Marshal(deltas)
		if err != nil {
			log.Println("Error marshalling to JSON: %v", err)
		}
		SetGitHubOutput("delta_files", string(jsonData))
	} else {
		SetGitHubOutput("is_detected", "false")
	}
}
