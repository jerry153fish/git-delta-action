package internal

import (
	"log"
	"strings"

	"github.com/caarlos0/env/v11"
)

const (
	// FileSeparator is used to split the Files and IgnoreFiles strings into slices
	FileSeparator = "\n"
)

// InputConfig holds the configuration for the Action Inputs
type InputConfig struct {
	Environment          string `env:"INPUT_ENVIRONMENT"`
	Commit               string `env:"INPUT_COMMIT"`
	Includes             string `env:"INPUT_INCLUDES"`
	Excludes             string `env:"INPUT_EXCLUDES"`
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
	IncludesPatterns     []string
	ExcludesPatterns     []string
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
