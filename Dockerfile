# Build stage
FROM golang:1.21-alpine AS builder

# Install git for go mod download
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pomodoro .

# Final stage
FROM alpine:latest

# Install ca-certificates for any HTTPS calls
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 pomodoro && \
    adduser -D -s /bin/sh -u 1001 -G pomodoro pomodoro

WORKDIR /home/pomodoro

# Copy the binary from builder stage
COPY --from=builder /app/pomodoro .

# Change ownership to non-root user
RUN chown pomodoro:pomodoro pomodoro

# Switch to non-root user
USER pomodoro

# Set the binary as entrypoint
ENTRYPOINT ["./pomodoro"]

# Default arguments (25 min work, 5 min break)
CMD ["25", "5"]

# Metadata
LABEL maintainer="Your Name <your.email@example.com>"
LABEL description="Production-ready Pomodoro CLI timer"
LABEL version="1.0.0"
