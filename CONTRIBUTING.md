# Contributing to Butler Coffee CLI

Thank you for considering contributing to Butler Coffee CLI! We welcome contributions from the community, whether it's bug reports, feature requests, documentation improvements, or code contributions.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Features](#suggesting-features)
  - [Submitting Pull Requests](#submitting-pull-requests)
- [Development Guidelines](#development-guidelines)
  - [Code Style](#code-style)
  - [Testing](#testing)
  - [Commit Messages](#commit-messages)
- [Project Structure](#project-structure)
- [Questions](#questions)

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please be respectful and constructive in all interactions.

## Getting Started

Before you begin contributing, please:

1. **Read the documentation**: Familiarize yourself with the [README.md](README.md) and [CLAUDE.md](CLAUDE.md) files
2. **Check existing issues**: Look for existing [issues](https://github.com/butlercoffee/bc-cli/issues) or discussions related to your contribution
3. **Fork the repository**: Create your own fork to work on changes

## Development Setup

### Prerequisites

- Go 1.25.4 or later
- Git
- Make (optional, but recommended)
- Python 3.14+ (for development tools)

### Setup Instructions

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/bc-cli.git
cd bc-cli

# Install dependencies and setup pre-commit hooks
make install

# Build the project
make compile

# Run the CLI locally
./bc-cli --help

# Run tests
go test ./...

# Run tests with coverage
go test ./... -cover
```

### Environment Variables for Local Development

```bash
# Point to local backend API (if running backend locally)
export BASE_HOSTNAME=http://localhost:8000
```

## How to Contribute

### Reporting Bugs

When reporting bugs, please include:

1. **Clear description**: Describe what you expected to happen and what actually happened
2. **Steps to reproduce**: Provide detailed steps to reproduce the issue
3. **Environment details**: OS, Go version, CLI version
4. **Error messages**: Include full error messages or stack traces
5. **Screenshots**: If applicable, add screenshots to help explain the problem

**Example bug report:**

```
**Bug**: Login fails with "invalid credentials" despite correct username/password

**Steps to reproduce:**
1. Run `bc-cli login`
2. Enter email: user@example.com
3. Enter password: correct_password
4. Error appears: "invalid credentials"

**Environment:**
- OS: macOS 14.6.0
- Go version: 1.25.4
- bc-cli version: 1.1.2

**Expected behavior:** User should be logged in successfully

**Actual behavior:** Login fails with error message

**Additional context:** This started happening after updating to version 1.1.2
```

### Suggesting Features

We welcome feature suggestions! When proposing new features:

1. **Check existing issues**: Search for similar feature requests
2. **Describe the use case**: Explain why this feature would be useful
3. **Provide examples**: Show how the feature would work
4. **Consider alternatives**: Mention any alternative solutions you've considered

### Submitting Pull Requests

1. **Create an issue first**: For significant changes, create an issue to discuss the approach
2. **Fork and branch**: Create a feature branch from `main` (e.g., `feature/add-coffee-timer`)
3. **Make your changes**: Follow the [development guidelines](#development-guidelines)
4. **Write tests**: Add tests for new functionality
5. **Update documentation**: Update README.md or CLAUDE.md if needed
6. **Run pre-commit hooks**: Ensure all checks pass
7. **Submit the PR**: Create a pull request with a clear description

**Pull Request Guidelines:**

- Keep changes focused and atomic (one feature/fix per PR)
- Write clear, descriptive commit messages
- Reference related issues (e.g., "Fixes #123")
- Ensure all tests pass
- Update CHANGELOG.md if appropriate

## Development Guidelines

### Code Style

This project follows Go's standard coding conventions:

- Run `gofmt` to format your code (automatically done by pre-commit hooks)
- Follow Go best practices and idioms
- Use meaningful variable and function names
- Add comments for complex logic

**TUI Components:**

- Place reusable components in `tui/components/`
- Use composed models in `tui/models/` for specific use cases
- Follow the existing pattern: components should be self-contained with their own Update/View methods
- Use centralized styles from `tui/styles/`

**API Layer:**

- Keep API methods in appropriate files (`api/auth.go`, `api/subscriptions.go`, etc.)
- Handle errors gracefully with user-friendly messages
- Use structured error handling with Django REST format

### Testing

- Write unit tests for new functionality
- Place tests alongside the code in `*_test.go` files
- Use `httptest.NewServer` for API endpoint tests
- Test both success and error cases
- Aim for good test coverage (run `go test ./... -cover`)

**Example test:**

```go
func TestListSubscriptions(t *testing.T) {
    // Create mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/api/core/v1/subscriptions" {
            t.Errorf("Expected path /api/core/v1/subscriptions, got %s", r.URL.Path)
        }
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "data": []interface{}{},
        })
    }))
    defer server.Close()

    // Test your function
    // ...
}
```

### Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/) for commit messages:

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `test:` Adding or updating tests
- `refactor:` Code refactoring
- `chore:` Maintenance tasks

**Examples:**

```
feat: add coffee timer command for pour-over brewing
fix: resolve token refresh loop on 401 errors
docs: update installation instructions for Windows
test: add tests for subscription management API
```

## Project Structure

```
bc-cli/
├── api/              # API client and endpoint methods
├── cmd/              # Cobra command definitions
├── config/           # Configuration management
├── templates/        # Text templates for UI
├── tui/              # Terminal UI components
│   ├── components/   # Reusable TUI components
│   ├── models/       # Composed TUI models
│   ├── prompts/      # Prompt wrapper functions
│   └── styles/       # Centralized styling
├── utils/            # Utility functions
├── main.go           # Entry point
├── CLAUDE.md         # AI assistant guidance
└── README.md         # Project documentation
```

For detailed architecture information, see [CLAUDE.md](CLAUDE.md).

## Questions

If you have questions about contributing:

1. Check the [README.md](README.md) and [CLAUDE.md](CLAUDE.md)
2. Search [existing issues](https://github.com/butlercoffee/bc-cli/issues)
3. Create a new issue with the `question` label

---

Thank you for contributing to Butler Coffee CLI! Your efforts help make specialty coffee more accessible to everyone.
