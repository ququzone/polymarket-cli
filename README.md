# polymarket-cli

A CLI tool for Polymarket.

## Features

- Command-line interface built with [Cobra](https://github.com/spf13/cobra)
- Configuration management with [Viper](https://github.com/spf13/viper)
- Easy to extend with new commands

## Installation

```bash
make build
```

The binary will be created in `bin/polymarket-cli`.

## Usage

```bash
# Show help
polymarket-cli --help

# Run example command
polymarket-cli example

# Show version
polymarket-cli version

# Use custom config file
polymarket-cli --config /path/to/config.yaml example
```

## Configuration

Create a configuration file at `~/.polymarket-cli.yaml`:

```yaml
api_key: "your-api-key"
api_secret: "your-api-secret"
debug: false
output_format: "table"
```

## Development

```bash
# Run the application
make run

# Run tests
make test

# Format code
make fmt

# Install dependencies
make deps
```

## Project Structure

```
.
├── cmd/              # CLI commands
│   ├── root.go       # Root command
│   ├── version.go    # Version command
│   └── example.go    # Example subcommand
├── internal/         # Private application code
│   ├── config/       # Configuration management
│   └── client/       # API client
├── main.go           # Application entry point
├── Makefile          # Build commands
└── go.mod            # Go modules
```

## Adding New Commands

1. Create a new file in `cmd/` directory:

```go
package cmd

import (
    "github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
    Use:   "newcmd",
    Short: "Short description",
    Run: func(cmd *cobra.Command, args []string) {
        // Command logic here
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
    // Add flags here
    newCmd.Flags().StringP("flag", "f", "", "Flag description")
}
```

2. Rebuild the application:

```bash
make build
```
