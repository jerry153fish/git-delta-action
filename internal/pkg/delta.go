package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
)

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
			log.Fatalf("Error getting diff between commits: %v", err)
		}
	} else {
		diffs, err = CompareGitFolderSHAs(repoPath, baseSha, cfg.Sha)
		if err != nil {
			log.Fatalf("Error getting diff between commits: %v", err)
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
