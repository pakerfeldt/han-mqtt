# Build stage
FROM golang:1.24-alpine AS builder

# Install git and bash for dependency fetching and scripts
RUN apk add --no-cache git bash

WORKDIR /app

# Add Go module files
COPY go.mod go.sum ./
RUN go mod download

# Add the source code
COPY . .

# Build the binary
RUN go build -o han-mqtt .

# Runtime stage
FROM alpine:latest

# Copy only the binary and config
WORKDIR /app
COPY --from=builder /app/han-mqtt .
COPY config.yaml .

# Add certificate bundle for MQTT TLS support
RUN apk add --no-cache ca-certificates

# Run the binary
ENTRYPOINT ["./han-mqtt"]