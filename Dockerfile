# Builder stage: use an Alpine-based Go image to build the binary
FROM golang:1.24.0-alpine3.21 AS builder

# Set the working directory for the Go project
WORKDIR /app

# Copy dependency files and download modules
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# Copy the entire source code into the container
COPY . .

# Build the binary from the main package in cmd/market-rest-api
RUN CGO_ENABLED=0 go build -o market-rest-api ./cmd/market-rest-api

# Final stage: use a minimal Alpine image to run the application
FROM alpine:3.21

# Set the working directory in the final image
WORKDIR /

# Copy the statically compiled binary from the builder stage
COPY --from=builder /app/market-rest-api .

# Copy the docs folder (needed for serving openapi.yaml)
COPY --from=builder /app/docs ./docs

# Run the application
CMD ["/market-rest-api"]
