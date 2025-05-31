# gogobox

A powerful command-line tool collection for bash environments, providing utilities for file management, time formatting, and more.

## Overview

gogobox is a versatile CLI toolkit built in Go that combines multiple useful utilities into a single binary. It features an intuitive command structure and includes both command-line interfaces and interactive Terminal User Interfaces (TUI) powered by Bubble Tea.

## Features

- **MinIO Operations**: Upload, download, and manage files with MinIO/S3-compatible storage
- **Time Formatting**: Convert between various time formats, timestamps, and timezones
- **Interactive TUI Components**: Rich terminal user interfaces for enhanced user experience
- **Cross-platform**: Works on Linux, macOS, and Windows
- **Extensible**: Modular architecture for easy addition of new commands

## Installation

### From Source

```bash
git clone https://github.com/gogodjzhu/gogobox.git
cd gogobox
go build -o gogobox ./cmd/main.go
```

### Using Go Install

```bash
go install github.com/gogodjzhu/gogobox/cmd@latest
```

## Commands

### Version

```bash
gogobox version
```

Get the current version of gogobox.

### MinIO Operations

Upload files to MinIO/S3-compatible storage:

```bash
gogobox minio upload [files...] [flags]
```

**Flags:**
- `--endpoint`: MinIO server endpoint (e.g., "localhost:9000")
- `--access-key-id`: Access key for authentication
- `--secret-access-key`: Secret key for authentication
- `--bucket-name`: Target bucket name
- `--use-ssl`: Use HTTPS for connection
- `--region`: Server region (default: "us-east-1")

**Example:**
```bash
gogobox minio upload file1.jpg file2.png \
  --endpoint localhost:9000 \
  --access-key-id minioadmin \
  --secret-access-key minioadmin \
  --bucket-name my-bucket
```

### Time Formatting

Convert between various time formats and timestamps:

```bash
gogobox timefmt [input] [flags]
```

**Flags:**
- `-i, --input-format`: Input time format (auto-detected if not specified)
- `-o, --output-format`: Output time format (default: "2006-01-02 15:04:05")
- `-t, --timestamp`: Output as timestamp instead of formatted string
- `-z, --timezone`: Timezone for output (e.g., 'UTC', 'America/New_York')
- `-u, --unit`: Timestamp unit: 'ms' (milliseconds) or 's' (seconds)

**Examples:**

Parse timestamp and format as date:
```bash
gogobox timefmt 1640995200000
# Output: 2022-01-01 00:00:00
```

Parse date string and convert to timestamp:
```bash
gogobox timefmt "2022-01-01 00:00:00" --timestamp
# Output: 1640995200000
```

Convert between different string formats:
```bash
gogobox timefmt "2022-01-01" --output-format "Jan 2, 2006"
# Output: Jan 1, 2022
```

Parse with specific input format:
```bash
gogobox timefmt "01/01/2022" --input-format "01/02/2006"
# Output: 2022-01-01 00:00:00
```

## TUI Components

gogobox includes several interactive terminal user interface components:

- **List**: Interactive lists with custom callbacks and key bindings
- **Text Input**: Prompt users for text input with placeholders
- **Result Selection**: Present choices to users with keyboard navigation

These components are built using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Bubbles](https://github.com/charmbracelet/bubbles) for a rich terminal experience.

## Architecture

The project follows a modular architecture:

- `cmd/`: Main application entry point
- `pkg/cmd/`: Individual command implementations
- `pkg/cmdutil/`: Shared utilities and TUI components
- `internal/`: Internal utilities and helpers
- `demo/`: Example applications and demos

## Dependencies

- [Cobra](https://github.com/spf13/cobra): CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea): TUI framework
- [MinIO Go Client](https://github.com/minio/minio-go): MinIO/S3 operations
- [Logrus](https://github.com/sirupsen/logrus): Structured logging

## Development

### Building

```bash
go build -o gogobox ./cmd/main.go
```

### Running Tests

```bash
go test ./...
```

### Adding New Commands

1. Create a new package under `pkg/cmd/`
2. Implement the command using Cobra
3. Register the command in `pkg/cmd/root/root.go`

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

## License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.

## Website

Visit [www.gogobox.xyz](http://www.gogobox.xyz) for more information and documentation.
