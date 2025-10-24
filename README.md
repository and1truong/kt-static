# kt-crawler

`kt-crawler` is an event-driven web crawler written in Go, designed for scanning and processing content related to translations, books, and chapters. It utilizes `goquery` for efficient HTML parsing and an internal event dispatcher for managing the crawling workflow.

## Features

- Event-driven architecture for scalable crawling logic.
- Configurable initial URL and crawling parameters.
- Uses `goquery` for robust HTML scraping.
- Pluggable listeners for handling different types of scanned content (Books, Translations, Chapters, Store).
- Integrated caching mechanism.

## Prerequisites

- Go (version 1.25 or higher)

## Getting Started

### Configuration

The crawler requires a configuration file in JSON format. A sample is provided at `config.sample.json`.

To run the crawler, you must specify the path to your configuration file using the `-config` flag. The configuration must include an `InitialURL` with a `v` query parameter (translation code) to start the process.

### Running

Use the provided `Makefile` targets for common operations:

| Command | Description |
| :--- | :--- |
| `make go-run` | Runs the crawler using `config.sample.json`. |
| `make go-test` | Runs all unit tests with race detection. |
| `make go-lint` | Runs the linter (`golangci-lint`). |
| `make go-fmt` | Formats all Go source files. |
| `make go-vest` | Runs `go vet` for static analysis. |

### Example Run

```bash
go run ./cmd/kt-crawler -config config.json
```

## Project Structure

The core logic is organized as follows:

- `cmd/kt-crawler/main.go`: Application entry point and configuration loading.
- `internal/`: Contains core business logic, services, and types.
    - `internal/services/eventdispatcher`: The core event handling mechanism.
    - `internal/services/cache`: Caching service implementation.
    - `internal/listeners`: Implementations of listeners for various scan events.
