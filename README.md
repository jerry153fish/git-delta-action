# Git Delta Action

This GitHub Action helps you visualize and manage file changes (deltas) across different environments, branches, or commits. It uses the `delta` tool to show diffs for specific files and folders, providing insights either offline or by querying GitHub via its REST API.

## Features

- Compare file changes (delta) between different Git branches, commits, or environments.
- Specify file patterns to include or exclude from the comparison.
- Run the delta check either offline or by querying the GitHub API.

## Inputs

| Name              | Description                                                                                              | Required | Default      |
|-------------------|----------------------------------------------------------------------------------------------------------|----------|--------------|
| `github_token`    | GitHub token for querying the GitHub REST API (used when comparing against environments).                 | No       | N/A          |
| `environment`     | The environment to compare against (requires GitHub token if used).                                       | No       | N/A          |
| `branch`          | The base branch to compare against.                                                                      | No       | `main`       |
| `commit`          | Specific commit to compare against.                                                                      | No       | N/A          |
| `includes`        | [Shell Glob style](https://teaching.idallen.com/cst8207/18w/notes/190_glob_patterns.html) patterns to include in the delta calculation, separated by newlines (`\n`).                        | No       | `""`         |
| `excludes`        | [Shell Glob style](https://teaching.idallen.com/cst8207/18w/notes/190_glob_patterns.html) patterns to exclude from the delta calculation, separated by newlines (`\n`). Excludes are applied after includes. | No       | `""`         |
| `online`          | Whether to run the delta comparison online using the GitHub API (`true`) or offline (`false`).            | No       | `true`       |

### Example of `includes` and `excludes`

```
includes: |
  live/local/*
  live/stag/ec2/terragrunt.hcl

excludes: |
  *.zip
  */**/README.md
```

## Outputs

| Name            | Description                                                             |
|-----------------|-------------------------------------------------------------------------|
| `delta_files`   | A JSON string with the paths of the files that have a delta (difference).|
| `is_detected`   | A boolean value indicating whether a delta was detected or not.          |

## Usage

Here's an example of how to use this action in your workflow:

```yaml
name: Check Git Delta

on: [push, pull_request]

jobs:
  git-delta:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v2

      - name: Run Git Delta
        uses: jerry153fish/git-delta-action@v0.0.1
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          branch: 'main'
          includes: |
            live/prod/*
            live/stag/*
            modules/*
          excludes: |
            **/*.md
```

In this example, the Git Delta action compares the current branch against the `main` branch, considering only files in `live/prod//` and `live/prod//` while excluding markdown files (`*.md`).

## Docker Image

This action runs inside a Docker container. The latest version of the image used is:

docker://jerry153fish/git-delta:latest

## Development

To contribute to this project or make modifications:

1. Fork the repository and clone it to your local machine.
2. Make your changes or additions to the codebase.
3. Test your changes thoroughly.
4. Create a pull request with a clear description of your modifications.

### Prerequisites
For local development and testing:

1. Ensure you have [Docker](https://www.docker.com/) and [act](https://github.com/nektos/act) installed on your machine.
2. Create a `.secrets` file in the root directory of the project with the following content:


```
GITHUB_TOKEN=your_github_token

# Optional: only i you want to push the docker image
DOCKER_USER=your_docker_username
DOCKER_TOKEN=your_docker_password
```

### Useful commands

1. Run the following commands for help
   
```
make help
```

2. Install dependencies

```
make install
```

3. Lint
   
```
make fmt 
make lint
```
   
4. Tests

```
make test # Unit tests
make vet # Vet checks
make act-test # Runs `test` stage in the workflow of `.github/workflows/ci.yaml`
make act-sanity # Runs `sanity-check` stage in then workflow of `.github/workflows/ci.yaml`
```

5. Build docker build

```
make act-docker # Run the whole workflow of `.github/workflows/ci.yaml`
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.