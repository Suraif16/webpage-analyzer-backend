# Web Page Analyzer Backend

A Go service that analyzes web pages for HTML version, title, headings with different levels, internal & external links, inaccessible links and whether it contains a login form or not.

Check out demo video of application:  
[Demo URL](https://youtu.be/WPlctqzm0u8)

## Prerequisites

Ensure you have the following installed:

- **Go**: 1.21 or higher
- **make** (optional)

## Setup

1. Install dependencies:
   ```sh
   go mod download
   ```

2. Running the Server
   ```sh
   go run cmd/api/main.go
   ```
Access the server at: [http://localhost:8080](http://localhost:8080)

## API Documentation

After starting the server, visit:
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## Testing

Run tests:
```sh
go test ./... -v
```

## Project Structure

```
webpage-analyzer-backend/
├── cmd/
│   └── api/                    # Application entry point
│       └── main.go
├── internal/
│   ├── core/         # Business logic
│   │   ├── domain/            # Business entities and errors
│   │   │   ├── types.go       # Data structures
│   │   │   └── errors.go      # Custom error definitions
│   │   ├── ports/             # Interfaces
│   │   │   └── analyzer.go    # Core interfaces
│   │   └── services/          # Business logic implementation
│   │       └── analyzer.go    # Main analysis service
│   ├── handlers/     # HTTP handlers
│   ├── middleware/   # HTTP middleware
│   └── infrastructure/  # External implementations
│       ├── http/             # HTTP client
│       │   └── client/
│       │       └── client.go
│       └── parser/           # HTML parsing
│           └── html_parser.go
├── tests/                     # Integration tests
├── docs/                      # Generated Swagger docs
└── .env              # Environment variables
```

## Libraries used & their purpose
- gin-gonic/gin - Web framework with high performance, middleware support
- uber-go/zap - Structured Logging
- stretchr/testify - Testing (with support for assertions, mocking)
- swaggo/swag - API documentation
- PuerkitoBio/goquery - HTML parsing

## Challenges faced and how I overcome them
- Mocking External Services for testing- Used mocks
- Writing repetitive test cases for different URL scenarios - Implemented table-driven tests
- Deploying to Microsoft Azure using docker.compose file - (Work in Progress)

## Additional Information

For more details, visit the documentation:  
[Documentation URL](https://docs.google.com/document/d/18IrcFGb_ur-Axp3A0NRtFfond7CdH8vCVmjz4spNSyg/edit?tab=t.0#heading=h.vwi3xxoqbucr)

## License

This project is licensed under the MIT License.

