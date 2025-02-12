# Web Page Analyzer Backend

A Go service that analyzes web pages for HTML version, headings, links, and more.

## Prerequisites

Ensure you have the following installed:

- **Go**: 1.21 or higher
- **make** (optional)

## Setup

1. Clone the repository:
   ```sh
   git clone <repository-url>
   cd webpage-analyzer-backend
   ```

2. Install dependencies:
   ```sh
   go mod download
   ```

3. Create a `.env` file in the root directory:
   ```env
   PORT=8080
   GIN_MODE=release
   ALLOWED_ORIGINS=http://localhost:3000
   ```

## Running the Application

Start the server:
```sh
go run cmd/api/main.go
```
Access the application at: [http://localhost:8080](http://localhost:8080)

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
│   └── api/          # Application entry point
├── internal/
│   ├── core/         # Business logic
│   ├── handlers/     # HTTP handlers
│   ├── middleware/   # HTTP middleware
│   └── infrastructure/  # External implementations
└── .env              # Environment variables
```

## License

This project is licensed under the MIT License.

