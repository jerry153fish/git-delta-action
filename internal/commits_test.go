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

func TestGetDiffBetweenCommits(t *testing.T) {
	sha1 := "c6023e778dac2c67e7ec0c42889e349a76414294"
	sha2 := "839bc7c55038951cfd3fed884617fd80d02ddbd5"

	result, err := GetDiffBetweenCommits("../", sha1, sha2)
	if err != nil {
		t.Fatalf("Error getting diff between commits: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result, got nil")
	}

	expectedDiff := []string{"Dockerfile", "go.mod", "go.sum"}

	if !reflect.DeepEqual(result, expectedDiff) {
		t.Errorf("GetDiffBetweenCommits() = %v, want %v", result, expectedDiff)
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
