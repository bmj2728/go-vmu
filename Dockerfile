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

# Install su-exec for safe user switching (used by entrypoint.sh)
# Create a default non-root user and group with placeholder IDs
# These IDs will be changed by the entrypoint.sh script at runtime
RUN apk add --no-cache su-exec \
    && addgroup -S appgroup \
    && adduser -S -G appgroup appuser

# Copy application binary to standard location
COPY --from=builder /app/vmu /usr/local/bin/vmu

# Copy the entrypoint script and make it executable
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Create videos directory as mount point with proper permissions
RUN mkdir -p /videos && chmod 777 /videos

# Declare the volume mount point
VOLUME ["/videos"]

# Set the entrypoint to our script
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

# Set default command if no arguments are provided to docker run
CMD ["/videos"]