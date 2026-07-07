# Contributing to terrapilot

Thanks for your interest in contributing. Here's how to get started.

## Reporting issues

- Use [GitHub Issues](https://github.com/laurencechau/terrapilot/issues) for bug reports and feature requests
- Check existing issues before opening a new one
- For bugs, include your OS, terrapilot version, and a minimal reproduction

## Pull requests

1. Fork the repo and create a branch from `main`
2. Make your changes
3. Run tests: `go test ./...`
4. Open a pull request with a clear description of the change

## Development setup

Requires Go 1.22+.

```bash
git clone https://github.com/laurencechau/terrapilot.git
cd terrapilot
make build
make test
```

## Design principles

Before contributing a feature, make sure it fits the project's core principles:

- No code generation
- No abstraction layers
- Each stack must remain a valid standalone Terraform/OpenTofu stack
- No new syntax — pure HCL only

If you're unsure whether a feature fits, open an issue to discuss it first.

## Code of conduct

Be respectful. See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
