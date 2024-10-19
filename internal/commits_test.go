package internal

import (
	"os"
	"reflect"
	"testing"
)

func TestGetInputConfig(t *testing.T) {
	// Set up test environment variables
	os.Setenv("INPUT_ENVIRONMENT", "test")
	os.Setenv("INPUT_COMMIT", "abc123")
	os.Setenv("INPUT_FILES", "file1.txt\nfile2.txt")
	os.Setenv("INPUT_IGNORE_FILES", "ignore1.txt\nignore2.txt")
	os.Setenv("INPUT_DELTA_OUTPUT_PATH_DEPTH", "2")
	os.Setenv("INPUT_GITHUB_TOKEN", "ghp_testtoken")
	os.Setenv("INPUT_BRANCH", "main")

	os.Setenv("GITHUB_SHA", "def456")
	os.Setenv("GITHUB_REF", "refs/heads/main")
	os.Setenv("GITHUB_API_URL", "https://api.github.com")
	os.Setenv("GITHUB_WORKFLOW", "test-workflow")
	os.Setenv("GITHUB_EVENT_NAME", "push")
	os.Setenv("GITHUB_JOB", "build")
	os.Setenv("GITHUB_REPOSITORY", "test/test")

	// Call the function
	ic := GetInputConfig()

	// Check the results
	expectedIC := InputConfig{
		Environment:          "test",
		Commit:               "abc123",
		Files:                "file1.txt\nfile2.txt",
		IgnoreFiles:          "ignore1.txt\nignore2.txt",
		DeltaOutputPathDepth: "2",
		GithubToken:          "ghp_testtoken",
		FilePatterns:         []string{"file1.txt", "file2.txt"},
		IgnoreFilePatterns:   []string{"ignore1.txt", "ignore2.txt"},
		Sha:                  "def456",
		Ref:                  "refs/heads/main",
		ApiUrl:               "https://api.github.com",
		Workflow:             "test-workflow",
		EventName:            "push",
		Job:                  "build",
		Repo:                 "test/test",
		Branch:               "main",
	}

	if !reflect.DeepEqual(ic, expectedIC) {
		t.Errorf("GetInputConfig() = %v, want %v", ic, expectedIC)
	}
}

func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	// Clean up
	os.Unsetenv("INPUT_ENVIRONMENT")
	os.Unsetenv("INPUT_COMMIT")
	os.Unsetenv("INPUT_FILES")
	os.Unsetenv("INPUT_IGNORE_FILES")
	os.Unsetenv("INPUT_DELTA_OUTPUT_PATH_DEPTH")
	os.Unsetenv("INPUT_GITHUB_TOKEN")
	os.Unsetenv("INPUT_BRANCH")
	os.Unsetenv("GITHUB_SHA")
	os.Unsetenv("GITHUB_REF")
	os.Unsetenv("GITHUB_API_URL")
	os.Unsetenv("GITHUB_WORKFLOW")
	os.Unsetenv("GITHUB_EVENT_NAME")
	os.Unsetenv("GITHUB_JOB")
	os.Unsetenv("GITHUB_REPOSITORY")

	os.Exit(code)
}
