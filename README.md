# gmapcli

A Go CLI and library for the Google Maps Platform **Places API (New)** and **Routes API**.

The command-line binary is named **`googlemapscli`**.

## Features

- Text place search with keyword, category, open-now, rating, price, location-bias, and pagination options.
- Place autocomplete suggestions with optional session-token support.
- Nearby place search.
- Place details with hours, contact information, reviews, and photos.
- Photo media URLs.
- Location lookup from free-form text.
- Route search for places along a path.
- Directions with walking, driving, bicycling, or transit modes.
- Human-readable terminal output or `--json`.
- Reusable Go library under `pkg/maps`.

## Requirements

- Go 1.26 or newer.
- Git is not required.
- A Google Maps Platform API key with the **Places API (New)** and **Routes API** enabled.

## Install

Install the CLI directly:

```bash
go install github.com/krishnashahane/gmapcli/cmd/googlemapscli@latest
```

The binary will be installed as `googlemapscli` in your Go bin directory.

The CLI defaults Google Maps `languageCode` to `en`; use `--lang` to request another supported language. Google may still return local-language/transliterated address content by design.

Or build locally:

```bash
git clone https://github.com/krishnashahane/gmapcli.git
cd gmapcli
make build
```

## Configure

Set the API key through the environment:

macOS/Linux:

```bash
export GOOGLE_PLACES_API_KEY="your-key"
```

Windows PowerShell:

```powershell
$env:GOOGLE_PLACES_API_KEY="your-key"
```

Restrict the key in Google Cloud to the APIs this tool uses and to the applications/environments where the key is actually needed.

## Usage

```text
googlemapscli [global flags] <command>
```

Commands:

```text
search
suggest
nearby
route
directions
info
photo
lookup
```

Examples:

```bash
googlemapscli search "coffee" --min-score 4 --only-open --limit 5 \
  --lat 40.8065 --lng -73.9719 --radius 3000

googlemapscli suggest "cof" --limit 5

googlemapscli nearby --lat 47.6062 --lng -122.3321 --radius 1500 --category cafe

googlemapscli route "coffee" --from "Seattle, WA" --to "Portland, OR" --stops 5

googlemapscli directions --from "Pike Place Market" --to "Space Needle" --steps

googlemapscli info ChIJN1t_tDeuEmsRUsoyG83frY4 --reviews --photos

googlemapscli lookup "Riverside Park, New York" --limit 5

googlemapscli search "sushi" --json
```

Use `googlemapscli <command> --help` for command-specific options.

## Library

```go
package main

import (
    "context"
    "os"
    "time"

    "github.com/krishnashahane/gmapcli/pkg/maps"
)

func main() {
    gm := maps.NewGoogleMaps(maps.Settings{
        Key:     os.Getenv("GOOGLE_PLACES_API_KEY"),
        Timeout: 8 * time.Second,
    })

    result, err := gm.TextSearch(context.Background(), maps.TextSearchInput{
        Text:       "italian restaurant",
        MaxResults: 10,
        Vicinity: &maps.Area{
            Latitude:  40.8065,
            Longitude: -73.9719,
            Radius:    3000,
        },
    })
    _ = result
    _ = err
}
```

## Security and reliability

- API keys are sent in the Google API header, not query strings.
- CLI API requests default to English (`--lang` overrides the language); user-generated Google content may still appear in its original language.
- GitHub-style or shell-style command injection is avoided because CLI inputs are passed to the HTTP client as structured values.
- API base URLs are restricted to Google's official HTTPS hosts; HTTP is allowed only for loopback local testing.
- HTTP redirects are disabled so the API key is not forwarded to an unexpected redirect target.
- Response bodies are capped at 1 MiB.
- Inputs such as coordinates, result counts, route samples, photo sizes, and navigation parameters are validated.
- Place IDs and photo resource names are validated before being inserted into URL paths.
- The default HTTP client has a bounded timeout.
- The CLI never prints the configured API key as part of normal error handling.

## Development and tests

Run the test suite:

```bash
make test
```

Run static analysis:

```bash
make lint
```

Generate coverage:

```bash
make coverage
```

Format code before committing:

```bash
gofmt -w cmd internal pkg
```

For a dependency vulnerability scan, use the official Go vulnerability tooling:

```bash
govulncheck ./...
```

## Project layout

```text
gmapcli/
├── cmd/googlemapscli/       # CLI entry point
├── internal/terminal/       # CLI commands, formatting, and execution
├── pkg/maps/                # Public Google Maps client library
├── Makefile
├── go.mod
├── go.sum
└── LICENSE
```

## License

MIT
