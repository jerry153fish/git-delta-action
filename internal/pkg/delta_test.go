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
	os.Setenv("INPUT_INCLUDES", "file1.txt\nfile2.txt")
	os.Setenv("INPUT_EXCLUDES", "ignore1.txt\nignore2.txt")
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
		Environment:      "test",
		Commit:           "abc123",
		Includes:         "file1.txt\nfile2.txt",
		Excludes:         "ignore1.txt\nignore2.txt",
		GithubToken:      "ghp_testtoken",
		IncludesPatterns: []string{"file1.txt", "file2.txt"},
		ExcludesPatterns: []string{"ignore1.txt", "ignore2.txt"},
		Sha:              "def456",
		Ref:              "refs/heads/main",
		ApiUrl:           "https://api.github.com",
		Workflow:         "test-workflow",
		EventName:        "push",
		Job:              "build",
		Repo:             "test/test",
		Branch:           "main",
	}

	if !reflect.DeepEqual(ic, expectedIC) {
		t.Errorf("GetInputConfig() = %v, want %v", ic, expectedIC)
	}
}

func TestFilterStrings(t *testing.T) {
	tests := []struct {
		input           []string
		includePatterns []string
		excludePatterns []string
		expectedResult  []string
		expectedErrors  []string
	}{
		{
			input:           []string{"file1.txt", "file2.txt", "file3.log"},
			includePatterns: []string{".*.txt$"},
			excludePatterns: []string{"file2.txt"},
			expectedResult:  []string{"file1.txt"},
			expectedErrors:  []string{},
		},
		{
			input:           []string{"image.png", "document.pdf", "notes.txt"},
			includePatterns: []string{".*\\.(png|pdf)$"},
			excludePatterns: []string{"document.pdf"},
			expectedResult:  []string{"image.png"},
			expectedErrors:  []string{},
		},
		{
			input:           []string{"test.go", "main.go", "script.js"},
			includePatterns: []string{".*\\.go$"},
			excludePatterns: []string{"main.go"},
			expectedResult:  []string{"test.go"},
			expectedErrors:  []string{},
		},
		{
			input:           []string{"data.csv", "data.json"},
			includePatterns: []string{".*\\.csv$"},
			excludePatterns: []string{".*"},
			expectedResult:  []string{},
			expectedErrors:  []string{},
		},
		{
			input:           []string{"aa/bb/data.csv", "data.json", "aa/cc/data.csv", "aa/data.json"},
			includePatterns: []string{"aa/*"},
			excludePatterns: []string{"aa/cc/*"},
			expectedResult:  []string{"aa/bb/data.csv", "aa/data.json"},
			expectedErrors:  []string{},
		},
		{
			input:           []string{"aa/bb/data.csv", "data.json", "aa/cc/data.csv", "aa/data.json"},
			includePatterns: []string{"aa/*"},
			excludePatterns: []string{"*.zip"},
			expectedResult:  []string{"aa/bb/data.csv", "aa/cc/data.csv", "aa/data.json"},
			expectedErrors:  []string{},
		},
		{
			input:           []string{"aa/bb/data.csv", "data.json", "aa/cc/data.csv", "aa/data.json"},
			includePatterns: []string{".*"},
			excludePatterns: []string{"aa/*", "*.zip"},
			expectedResult:  []string{"data.json"},
			expectedErrors:  []string{},
		},
	}

	for _, test := range tests {
		result := FilterStrings(test.input, test.includePatterns, test.excludePatterns)

		if len(result) != len(test.expectedResult) {
			t.Errorf("expected %v, got %v", test.expectedResult, result)
		}

		for i, res := range result {
			if res != test.expectedResult[i] {
				t.Errorf("expected %v, got %v", test.expectedResult[i], res)
			}
		}

	}
}

func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	// Clean up
	os.Unsetenv("INPUT_ENVIRONMENT")
	os.Unsetenv("INPUT_COMMIT")
	os.Unsetenv("INPUT_INCLUDES")
	os.Unsetenv("INPUT_EXCLUDES")
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
