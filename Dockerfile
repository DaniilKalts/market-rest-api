# Build stage
FROM golang:1.24.0-alpine3.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o market-rest-api ./cmd/market-rest-api

# Final stage
FROM scratch
WORKDIR /app
COPY --from=builder /app/market-rest-api .
COPY --from=builder /app/docs ./docs
EXPOSE 8080
ENTRYPOINT ["/app/market-rest-api"]
