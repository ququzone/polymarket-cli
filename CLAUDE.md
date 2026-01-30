# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A CLI tool for Polymarket built with Go, using [Cobra](https://github.com/spf13/cobra) for command structure and [Viper](https://github.com/spf13/viper) for configuration management.

## Architecture

### Command Structure (Cobra)

- **Root command** ([cmd/root.go](cmd/root.go)): Base command that handles configuration initialization and persistent flags
- **Subcommands**: Each subcommand lives in its own file in `cmd/` (e.g., [cmd/example.go](cmd/example.go))

Adding a new command:
1. Create a new file in `cmd/` with a `cobra.Command` variable
2. In `init()`, register with `rootCmd.AddCommand(newCmd)`
3. Add any flags in the `init()` function

### Configuration Flow

1. `cobra.OnInitialize(initConfig)` is registered in [cmd/root.go:26](cmd/root.go#L26)
2. `initConfig()` reads config from:
   - `--config` flag if provided
   - Otherwise `~/.polymarket-cli.yaml`
3. Environment variables are automatically loaded via `viper.AutomaticEnv()`
4. `config.Init()` ([internal/config/config.go](internal/config/config.go)) initializes the global `AppCfg` variable

### Build Configuration

Version info is injected via ldflags during build:
- `main.Version`
- `main.BuildTime`
- `main.GitCommit`

## Code Style

- Remove unnecessary comments that simply restate what the code already makes obvious (e.g., `// versionCmd represents the version command`)
- Keep code self-documenting through clear naming
- Comments explaining business logic or non-obvious implementation details are acceptable

### After Generating Code

After writing or modifying code, always run the following commands to ensure code quality:

```bash
go mod tidy
go vet ./...
go fmt ./...
```
