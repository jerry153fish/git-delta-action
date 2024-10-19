package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-github/v66/github"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints. See issue #752.
	baseURLPath = "/api-v3"
)

func setup(t *testing.T) (client *github.Client, mux *http.ServeMux, serverURL string) {
	t.Helper()
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		fmt.Fprintln(os.Stderr, "\tSee https://github.com/google/go-github/issues/752 for information.")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the GitHub client being tested and is
	// configured to use test server.
	client = github.NewClient(nil)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url
	client.UploadURL = url

	t.Cleanup(server.Close)

	return client, mux, server.URL
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func TestGetLatestSuccessfulDeploymentSha(t *testing.T) {
	client, mux, _ := setup(t)

	// Mock the ListDeployments endpoint
	mux.HandleFunc("/repos/owner/repo/deployments", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id": 1, "sha": "abc123"}]`)
	})

	// Mock the ListDeploymentStatuses endpoint
	mux.HandleFunc("/repos/owner/repo/deployments/1/statuses", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"state": "success"}]`)
	})

	// Create a test InputConfig
	testConfig := &InputConfig{
		Environment: "production",
		GithubToken: "test-token",
		Repo:        "owner/repo",
	}

	// Call the function under test
	sha := GetLatestSuccessfulDeploymentSha(client, testConfig)

	// Assert the results
	if sha == "" {
		t.Fatal("Expected a deployment, got nil")
	}
	if sha != "abc123" {
		t.Errorf("Expected deployment SHA abc123, got %s", sha)
	}
}

func TestGetBranchLatestSHA(t *testing.T) {
	client, mux, _ := setup(t)

	// Mock the ListDeployments endpoint
	mux.HandleFunc("/repos/owner/repo/git/ref/heads/main", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
		{
			"ref": "refs/heads/b",
			"url": "https://api.github.com/repos/o/r/git/refs/heads/b",
			"object": {
				"type": "commit",
				"sha": "abc123",
				"url": "https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"
			}
		}`)
	})

	// Create a test InputConfig
	testConfig := &InputConfig{
		Branch:      "main",
		GithubToken: "test-token",
		Repo:        "owner/repo",
	}

	// Call the function under test
	sha := GetBranchLatestSHA(client, testConfig)

	// Assert the results
	if sha == "" {
		t.Fatal("Expected a deployment, got nil")
	}
	if sha != "abc123" {
		t.Errorf("Expected deployment SHA abc123, got %s", sha)
	}
}

// c6023e778dac2c67e7ec0c42889e349a76414294 839bc7c55038951cfd3fed884617fd80d02ddbd5

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

	if !contains(result, "Dockerfile") {
		t.Errorf("Expected Dockerfile to be inside result %v", result)
	}
}
