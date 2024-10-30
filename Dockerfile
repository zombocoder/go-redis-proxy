# First stage: Build the Go binary
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod  ./

# Copy the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o redis-proxy cmd/server/main.go

# Second stage: Create a minimal image with the binary
FROM alpine:latest

# Set up a working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/redis-proxy .

# Set up a volume to allow mounting external configuration files
VOLUME /config

# Set entrypoint to run the proxy, passing in config file path as argument
ENTRYPOINT ["./redis-proxy"]
