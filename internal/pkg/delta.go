package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bmatcuk/doublestar/v4"
)

// matchPatterns checks if the string matches any of the patterns using filepath.Match.
func matchPatterns(str string, include bool, patterns []string) bool {
	if len(patterns) == 0 {
		return include
	}

	for _, pattern := range patterns {
		if pattern != "" {
			matched, err := doublestar.Match(pattern, str)
			if err != nil {
				log.Printf("Error matching pattern '%s': %v", pattern, err)
				continue
			}
			if matched {
				return true
			}
		}
	}
	return false
}

// FilterStrings filters the input strings based on inclusion and exclusion patterns using filepath.Match.
func FilterStrings(input []string, includePatterns, excludePatterns []string) []string {
	var result []string

	// Filter the input strings
	for _, str := range input {
		// Check if the string matches any include pattern
		if matchPatterns(str, true, includePatterns) {
			// If it matches an include pattern, check that it doesn't match any exclude pattern
			if !matchPatterns(str, false, excludePatterns) {
				if str != "" {
					result = append(result, str)
				}
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
	var diffs []string
	var err error
	if cfg.Environment != "" {
		baseSha = GetLatestSuccessfulDeploymentSha(client, &cfg)
	} else {
		baseSha = GetGitHubBranchLatestSHA(client, &cfg)
	}

	if cfg.online == "true" {
		diffs, err = CompareGithubSHAs(client, &cfg, baseSha)
		if err != nil {
			log.Panicf("Error getting diff between commits: %v", err)
		}
	} else {
		diffs, err = CompareGitFolderSHAs(repoPath, baseSha, cfg.Sha)
		if err != nil {
			log.Panicf("Error getting diff between commits: %v", err)
		}
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
