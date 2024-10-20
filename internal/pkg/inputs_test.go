package internal

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInputConfigWithEmptyEnvironment(t *testing.T) {
	os.Setenv("INPUT_ENVIRONMENT", "")
	os.Setenv("INPUT_COMMIT", "abc123")
	os.Setenv("INPUT_INCLUDES", "file1.txt")
	os.Setenv("INPUT_EXCLUDES", "ignore1.txt")
	os.Setenv("INPUT_GITHUB_TOKEN", "ghp_testtoken")
	os.Setenv("INPUT_BRANCH", "develop")

	os.Setenv("GITHUB_SHA", "def456")
	os.Setenv("GITHUB_REF", "refs/heads/develop")
	os.Setenv("GITHUB_API_URL", "https://api.github.com")
	os.Setenv("GITHUB_WORKFLOW", "test-workflow")
	os.Setenv("GITHUB_EVENT_NAME", "pull_request")
	os.Setenv("GITHUB_JOB", "test")
	os.Setenv("GITHUB_REPOSITORY", "test/repo")

	ic := GetInputConfig()

	if ic.Environment != "" {
		t.Errorf("Expected empty Environment, got %s", ic.Environment)
	}
	if ic.Commit != "abc123" {
		t.Errorf("Expected Commit abc123, got %s", ic.Commit)
	}
	if ic.Branch != "develop" {
		t.Errorf("Expected Branch develop, got %s", ic.Branch)
	}
	if ic.EventName != "pull_request" {
		t.Errorf("Expected EventName pull_request, got %s", ic.EventName)
	}
}

func TestGetInputConfigWithMultipleIncludesAndExcludes(t *testing.T) {
	os.Setenv("INPUT_ENVIRONMENT", "staging")
	os.Setenv("INPUT_COMMIT", "")
	os.Setenv("INPUT_INCLUDES", "file1.txt\nfile2.go\nfile3.js")
	os.Setenv("INPUT_EXCLUDES", "ignore1.txt\nignore2.log\nignore3.tmp")
	os.Setenv("INPUT_GITHUB_TOKEN", "ghp_testtoken")
	os.Setenv("INPUT_BRANCH", "")

	os.Setenv("GITHUB_SHA", "abc789")
	os.Setenv("GITHUB_REF", "refs/heads/feature/new-feature")
	os.Setenv("GITHUB_API_URL", "https://api.github.com")
	os.Setenv("GITHUB_WORKFLOW", "feature-workflow")
	os.Setenv("GITHUB_EVENT_NAME", "push")
	os.Setenv("GITHUB_JOB", "build")
	os.Setenv("GITHUB_REPOSITORY", "org/repo")

	ic := GetInputConfig()

	if ic.Environment != "staging" {
		t.Errorf("Expected Environment staging, got %s", ic.Environment)
	}
	if len(ic.IncludesPatterns) != 3 {
		t.Errorf("Expected 3 IncludesPatterns, got %d", len(ic.IncludesPatterns))
	}
	if len(ic.ExcludesPatterns) != 3 {
		t.Errorf("Expected 3 ExcludesPatterns, got %d", len(ic.ExcludesPatterns))
	}
	if ic.Sha != "abc789" {
		t.Errorf("Expected Sha abc789, got %s", ic.Sha)
	}
	if ic.Ref != "refs/heads/feature/new-feature" {
		t.Errorf("Expected Ref refs/heads/feature/new-feature, got %s", ic.Ref)
	}
}

func TestGetInputConfigWithEmptyOptionalFields(t *testing.T) {
	os.Setenv("INPUT_ENVIRONMENT", "production")
	os.Setenv("INPUT_COMMIT", "")
	os.Setenv("INPUT_INCLUDES", "")
	os.Setenv("INPUT_EXCLUDES", "")
	os.Setenv("INPUT_GITHUB_TOKEN", "ghp_testtoken")
	os.Setenv("INPUT_BRANCH", "")

	os.Setenv("GITHUB_SHA", "xyz123")
	os.Setenv("GITHUB_REF", "refs/tags/v1.0.0")
	os.Setenv("GITHUB_API_URL", "https://api.github.com")
	os.Setenv("GITHUB_WORKFLOW", "release-workflow")
	os.Setenv("GITHUB_EVENT_NAME", "release")
	os.Setenv("GITHUB_JOB", "deploy")
	os.Setenv("GITHUB_REPOSITORY", "company/project")

	ic := GetInputConfig()

	if ic.Environment != "production" {
		t.Errorf("Expected Environment production, got %s", ic.Environment)
	}
	if ic.Commit != "" {
		t.Errorf("Expected empty Commit, got %s", ic.Commit)
	}
	if len(ic.IncludesPatterns) != 0 {
		t.Errorf("Expected 0 IncludesPatterns, got %d", len(ic.IncludesPatterns))
	}
	if len(ic.ExcludesPatterns) != 0 {
		t.Errorf("Expected 0 ExcludesPatterns, got %d", len(ic.ExcludesPatterns))
	}
	if ic.Ref != "refs/tags/v1.0.0" {
		t.Errorf("Expected Ref refs/tags/v1.0.0, got %s", ic.Ref)
	}
	if ic.EventName != "release" {
		t.Errorf("Expected EventName release, got %s", ic.EventName)
	}
}

func TestValidatePatterns(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		patterns  []string
		wantPanic bool
	}{
		{
			name:      "Valid patterns",
			patterns:  []string{"*.txt", "file?.log", "[a-z]*.go"},
			wantPanic: false,
		},
		{
			name:      "Empty patterns",
			patterns:  []string{},
			wantPanic: false,
		},
		{
			name:      "Single invalid pattern",
			patterns:  []string{"[invalid"},
			wantPanic: true,
		},
		{
			name:      "Mix of valid and invalid patterns",
			patterns:  []string{"*.txt", "[invalid", "file?.log"},
			wantPanic: true,
		},
		{
			name:      "Complex valid patterns",
			patterns:  []string{"**/*.{js,ts,jsx,tsx}", "src/[a-z]*/**.go"},
			wantPanic: false,
		},
		{
			name:      "Patterns with escaped characters",
			patterns:  []string{"file\\*.txt", "log\\?.dat"},
			wantPanic: false,
		},
		{
			name:      "character range",
			patterns:  []string{"[a-z]*.txt"},
			wantPanic: false,
		},
		{
			name:      "Valid patterns with whitespace",
			patterns:  []string{"  *.txt  ", "  file?.log  "},
			wantPanic: false,
		},
		{
			name:      "Valid patterns with subdirectories",
			patterns:  []string{"live/prod/*", "live/local/*"},
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() { validatePatterns(tt.patterns) })
			} else {
				assert.NotPanics(t, func() { validatePatterns(tt.patterns) })
			}
		})
	}
}
func TestMain(m *testing.M) {
	// Save current environment
	oldEnv := os.Environ()

	// Run tests
	code := m.Run()

	// Restore environment
	for _, envVar := range oldEnv {
		pair := strings.SplitN(envVar, "=", 2)
		os.Setenv(pair[0], pair[1])
	}

	os.Exit(code)
}

func TestInputConfigValidate(t *testing.T) {
	os.Setenv("INPUT_ENVIRONMENT", "production")
	tests := []struct {
		name        string
		inputConfig InputConfig
		wantPanic   bool
	}{
		{
			name: "Valid config with environment and github token",
			inputConfig: InputConfig{
				Environment: "production",
				GithubToken: "ghp_validtoken",
				Repo:        "test/repo",
				Sha:         "abc123",
			},
			wantPanic: false,
		},
		{
			name: "Invalid config with environment but no github token",
			inputConfig: InputConfig{
				Environment: "staging",
				GithubToken: "",
				Repo:        "test/repo",
				Sha:         "def456",
			},
			wantPanic: true,
		},
		{
			name: "Valid config with online mode and github token",
			inputConfig: InputConfig{
				online:      "true",
				GithubToken: "ghp_validtoken",
				Repo:        "test/repo",
				Sha:         "ghi789",
			},
			wantPanic: false,
		},
		{
			name: "Invalid config with online mode but no github token",
			inputConfig: InputConfig{
				online:      "true",
				GithubToken: "",
				Repo:        "test/repo",
				Sha:         "jkl012",
			},
			wantPanic: true,
		},
		{
			name: "Valid config with offline mode",
			inputConfig: InputConfig{
				online: "false",
				Repo:   "test/repo",
				Sha:    "mno345",
			},
			wantPanic: false,
		},
		{
			name: "Invalid config with empty repo",
			inputConfig: InputConfig{
				Repo: "",
				Sha:  "pqr678",
			},
			wantPanic: true,
		},
		{
			name: "Invalid config with empty sha",
			inputConfig: InputConfig{
				Repo: "test/repo",
				Sha:  "",
			},
			wantPanic: true,
		},
		{
			name: "Valid config with includes and excludes patterns",
			inputConfig: InputConfig{
				Repo:             "test/repo",
				Sha:              "stu901",
				IncludesPatterns: []string{"*.go", "*.js"},
				ExcludesPatterns: []string{"vendor/*", "node_modules/*"},
			},
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() { tt.inputConfig.Validate() })
			} else {
				assert.NotPanics(t, func() { tt.inputConfig.Validate() })
			}

		})
	}
}
