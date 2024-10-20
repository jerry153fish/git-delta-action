package internal

import (
	"os"
	"testing"
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
func TestIsValidRegex(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		pattern string
		want    bool
		wantErr bool
	}{
		{
			name:    "Valid simple regex",
			pattern: "^[a-z]+$",
			want:    true,
			wantErr: false,
		},
		{
			name:    "Valid simple regex",
			pattern: "live/prod/*",
			want:    true,
			wantErr: false,
		},
		{
			name:    "Valid complex regex",
			pattern: "^(https?:\\/\\/)?([\\da-z\\.-]+)\\.([a-z\\.]{2,6})([\\/\\w \\.-]*)*\\/?$",
			want:    true,
			wantErr: false,
		},
		{
			name:    "Invalid regex - unmatched parenthesis",
			pattern: "([a-z]+",
			want:    false,
			wantErr: true,
		},
		{
			name:    "Empty string",
			pattern: "",
			want:    true,
			wantErr: false,
		},
		{
			name:    "Invalid regex - unescaped special character",
			pattern: "a+*",
			want:    false,
			wantErr: true,
		},
		{
			name:    "Invalid regex",
			pattern: "*.txt",
			want:    false,
			wantErr: true,
		},
		{
			name:    "Valid regex with quantifiers",
			pattern: "a{2,4}",
			want:    true,
			wantErr: false,
		},
		{
			name:    "Valid regex with character classes",
			pattern: "[0-9a-fA-F]{6}",
			want:    true,
			wantErr: false,
		},
		{
			name:    "Invalid regex with lookahead",
			pattern: "(?=.*[A-Z])(?=.*[a-z])(?=.*\\d)[A-Za-z\\d]{8,}",
			want:    false,
			wantErr: true,
		},
		{
			name:    "Invalid regex - unmatched square bracket",
			pattern: "[a-z",
			want:    false,
			wantErr: true,
		},
		{
			name:    "Valid regex with escaped special characters",
			pattern: "\\[.*\\]",
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isValidRegex(tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("isValidRegex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isValidRegex() = %v, want %v", got, tt.want)
			}
		})
	}
}
