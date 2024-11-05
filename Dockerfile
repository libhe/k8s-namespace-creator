# Stage 1: Build the Go application
FROM quay.io/flacatus/go-test-tools:latest AS builder

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o namespace-creator main.go

# Stage 2: Create a lightweight image
FROM registry.access.redhat.com/ubi9/ubi:latest

WORKDIR /app

COPY --from=builder /app/namespace-creator .

# Set the entry point for the application
ENTRYPOINT ["/app/namespace-creator"]
