FROM golang:1.23-alpine

# Add git for fetching dependencies
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /web-analyzer ./cmd/api

# Create a minimal production image
FROM alpine:3.18

WORKDIR /app

# Copy binary from builder
COPY --from=0 /web-analyzer .

EXPOSE ${PORT:-8080}

# Run the application
CMD ["./web-analyzer"]