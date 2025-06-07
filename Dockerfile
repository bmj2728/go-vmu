# Multi-stage build with hardened approach

# Builder stage: Build the Go application
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o vmu ./cmd/vmu

# Final stage: Create the runtime image
FROM jrottenberg/ffmpeg:4.4-alpine

# Add labels for better documentation
LABEL maintainer="Brian Jipson <brian.jipson@novelgit.com>"
LABEL description="Video Metadata Updater - Embeds metadata from NFO files into video files"
LABEL version="0.7.0"

# Create a non-root user for running the application
# ensuring correct permissions for bind-mounted volumes.
RUN addgroup -S -g 100 appgroup && adduser -S -u 1024 -G appgroup appuser

# Copy application binary to standard location
COPY --from=builder /app/vmu /usr/local/bin/vmu

# Create videos directory as mount point with proper permissions
RUN mkdir -p /videos && chmod 777 /videos

# Declare the volume mount point
VOLUME ["/videos"]

# Switch to non-root user
USER appuser

# Set entrypoint
ENTRYPOINT ["vmu"]