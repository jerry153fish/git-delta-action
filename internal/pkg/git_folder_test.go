package internal

import (
	"reflect"
	"testing"
)

func TestCompareGitFolderSHAs(t *testing.T) {
	t.Parallel()
	sha1 := "c6023e778dac2c67e7ec0c42889e349a76414294"
	sha2 := "839bc7c55038951cfd3fed884617fd80d02ddbd5"

	result, err := CompareGitFolderSHAs("../../", sha1, sha2)
	if err != nil {
		t.Fatalf("Error getting diff between commits: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result, got nil")
	}

	expectedDiff := []string{"Dockerfile", "go.mod", "go.sum"}

	if !reflect.DeepEqual(result, expectedDiff) {
		t.Errorf("CompareGitFolderSHAs() = %v, want %v", result, expectedDiff)
	}
}

// func TestGetGitFolderBranchLatestSHA(t *testing.T) {
// 	t.Parallel()
// 	branch := "main"
// 	result, err := GetGitFolderBranchLatestSHA("../../", branch)
// 	if err != nil {
// 		t.Fatalf("Error getting latest sha on %v: %v", branch, err)
// 	}

// 	if result == "" {
// 		t.Fatal("Expected non-nil result, got nil")
// 	}
// }
