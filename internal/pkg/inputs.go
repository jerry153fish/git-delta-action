package internal

import (
	"log"
	"regexp"
	"strings"

	"github.com/caarlos0/env/v11"
)

const (
	// FileSeparator is used to split the Files and IgnoreFiles strings into slices
	FileSeparator = "\n"
)

// InputConfig holds the configuration for the Action Inputs
type InputConfig struct {
	Environment      string `env:"INPUT_ENVIRONMENT"`
	Commit           string `env:"INPUT_COMMIT"`
	Includes         string `env:"INPUT_INCLUDES"`
	Excludes         string `env:"INPUT_EXCLUDES"`
	GithubToken      string `env:"INPUT_GITHUB_TOKEN"`
	Sha              string `env:"GITHUB_SHA"`
	Ref              string `env:"GITHUB_REF"`
	ApiUrl           string `env:"GITHUB_API_URL"`
	Workflow         string `env:"GITHUB_WORKFLOW"`
	EventName        string `env:"GITHUB_EVENT_NAME"`
	Job              string `env:"GITHUB_JOB"`
	Repo             string `env:"GITHUB_REPOSITORY"`
	Branch           string `env:"INPUT_BRANCH"`
	Offline          string `env:"INPUT_OFFLINE"`
	IncludesPatterns []string
	ExcludesPatterns []string
}

// GetInputConfig parses environment variables into an InputConfig struct
// and processes the Files and IgnoreFiles fields
func GetInputConfig() InputConfig {
	// Parse environment variables into InputConfig struct
	c, err := env.ParseAs[InputConfig]()
	if err != nil {
		log.Fatalf("Failed to parse InputConfig: %v", err)
	}

	// If Includes is not empty, split it into IncludesPatterns
	if c.Includes != "" {
		c.IncludesPatterns = strings.Split(c.Includes, FileSeparator)
	}

	// If Excludes is not empty, split it into ExcludesPatterns
	if c.Excludes != "" {
		c.ExcludesPatterns = strings.Split(c.Excludes, FileSeparator)
	}
	return c
}

// Validate checks if the required fields in InputConfig are set
func (c *InputConfig) Validate() {
	if c.Environment != "" && c.GithubToken == "" {
		log.Fatal("github_token must be specific when the environment is given")
	}

	if c.Repo == "" {
		log.Fatal("Unexpected error for retrieve repository from runner, please contact developer")
	}
	if c.Sha == "" {
		log.Fatal("Unexpected error for retrieve current commit from runner, please contact developer")
	}

	validatePatterns(c.IncludesPatterns)
	validatePatterns(c.ExcludesPatterns)
}

// validatePatterns checks that the provided patterns are valid regular expressions.
// If any pattern is invalid, it logs a fatal error with the invalid pattern and error.
func validatePatterns(patterns []string) {
	if len(patterns) > 0 {
		for _, pattern := range patterns {
			valid, err := isValidRegex(pattern)
			if !valid {
				log.Fatalf("Invalid regex pattern in IncludesPatterns: %s. Error: %v", pattern, err)
			}
		}
	}
}

// isValidRegex checks if the provided string is a valid regular expression.
func isValidRegex(pattern string) (bool, error) {
	_, err := regexp.Compile(pattern)
	if err != nil {
		return false, err // Return false and the error
	}
	return true, nil // Return true if no error
}
