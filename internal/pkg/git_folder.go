package internal

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// CompareGitFolderSHAs retrieves the list of files that have changed between two commits identified by their SHAs.
// It takes the repository path and the two commit SHAs as input parameters.
// Returns a slice of strings containing the names of the changed files and an error if any occurs.
func CompareGitFolderSHAs(repoPath, sha1, sha2 string) ([]string, error) {
	// Open the repository at the given path
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("could not open repository: %v", err)
	}

	// Get the commits corresponding to the given SHAs
	commit1, err := repo.CommitObject(plumbing.NewHash(sha1))
	if err != nil {
		return nil, fmt.Errorf("could not find commit for SHA %s: %v", sha1, err)
	}

	commit2, err := repo.CommitObject(plumbing.NewHash(sha2))
	if err != nil {
		return nil, fmt.Errorf("could not find commit for SHA %s: %v", sha2, err)
	}

	// Get the tree objects for both commits
	tree1, err := commit1.Tree()
	if err != nil {
		return nil, fmt.Errorf("could not get tree for commit %s: %v", sha1, err)
	}

	tree2, err := commit2.Tree()
	if err != nil {
		return nil, fmt.Errorf("could not get tree for commit %s: %v", sha2, err)
	}

	// Get the diff between the two trees
	changes, err := object.DiffTree(tree1, tree2)
	if err != nil {
		return nil, fmt.Errorf("could not get diff between trees: %v", err)
	}

	// Collect the filenames of changed files
	var diffFiles []string
	for _, change := range changes {
		// Append the file name (NewName) to the list of diff files
		diffFiles = append(diffFiles, change.To.Name)
	}

	return diffFiles, nil
}

// GetGitFolderBranchLatestSHA retrieves the latest commit of a given branch in a Git repository.
// It takes the repository path and the branch name as input parameters.
// Returns the commit hash as a string and an error if any occurs.
func GetGitFolderBranchLatestSHA(repoPath, branchName string) (string, error) {
	// Open the repository at the given path
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("could not open repository: %v", err)
	}

	// Get the reference for the specified branch
	ref, err := repo.Reference(plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branchName)), true)
	if err != nil {
		return "", fmt.Errorf("could not find reference for branch %s: %v", branchName, err)
	}

	// Get the commit object for the reference
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return "", fmt.Errorf("could not get commit object for branch %s: %v", branchName, err)
	}

	// Return the commit hash as a string
	return commit.Hash.String(), nil
}
