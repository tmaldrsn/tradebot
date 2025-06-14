# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o detector .

# Stage 2: Minimal runtime image
FROM alpine:latest

# Install ca-certificates (for HTTPS support)
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from the builder
COPY --from=builder /app/detector .

# Set entrypoint
ENTRYPOINT ["./detector"]