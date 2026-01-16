# Contributing to Squaremind

First off, thank you for considering contributing to Squaremind! It's people like you that make Squaremind such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by the [Squaremind Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the problem**
- **Provide specific examples** (code snippets, config files)
- **Describe the behavior you observed and what you expected**
- **Include logs and error messages**
- **Specify your environment** (OS, Go version, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide a detailed description of the proposed functionality**
- **Explain why this enhancement would be useful**
- **List any alternatives you've considered**

### Pull Requests

1. **Fork the repo** and create your branch from `main`
2. **Follow the coding style** (run `make fmt` and `make lint`)
3. **Add tests** for any new functionality
4. **Ensure all tests pass** (`make test`)
5. **Update documentation** as needed
6. **Write a clear commit message**

## Development Setup

### Prerequisites

- Go 1.21 or later
- Make
- Git

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/squaremind.git
cd squaremind

# Add upstream remote
git remote add upstream https://github.com/square-mind/squaremind.git

# Install dependencies
make deps

# Build
make build

# Run tests
make test

# Run linter
make lint
```

### Project Structure

```
squaremind/
├── cmd/sqm/           # CLI application
├── pkg/
│   ├── identity/      # Cryptographic identity
│   ├── agent/         # Agent runtime
│   ├── coordination/  # Gossip, market, consensus
│   ├── collective/    # Collective management
│   └── llm/           # LLM providers
├── sdk/               # TypeScript SDK
├── examples/          # Examples
└── docs/              # Documentation
```

### Coding Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Write descriptive variable and function names
- Add comments for exported functions
- Keep functions focused and small

### Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters
- Reference issues and pull requests when relevant

Examples:
```
Add reputation decay to agent lifecycle

Implement automatic reputation decay over time to prevent
stale reputation scores. Decay rate is configurable via
CollectiveConfig.ReputationDecay.

Fixes #123
```

### Testing

- Write unit tests for new functionality
- Ensure existing tests pass
- Aim for good test coverage
- Use table-driven tests where appropriate

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

## Architecture Guidelines

### Adding a New Capability Type

1. Add the constant to `pkg/identity/capabilities.go`
2. Update the CLI help text in `cmd/sqm/main.go`
3. Add tests
4. Update documentation

### Adding a New LLM Provider

1. Implement the `Provider` interface in `pkg/llm/`
2. Add configuration support in the CLI
3. Add tests with mocked responses
4. Update documentation

### Modifying the Gossip Protocol

1. Ensure backwards compatibility
2. Add migration path if message format changes
3. Update tests thoroughly
4. Document protocol changes

## Questions?

Feel free to open an issue with the "question" label or reach out on [Discord](https://discord.gg/squaremind).

## Recognition

Contributors will be recognized in our README and release notes. Thank you for helping make Squaremind better!
