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

// DeltaConfig holds the configuration for the delta processing
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
