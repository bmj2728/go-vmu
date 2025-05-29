# Multi-stage build with hardened approach
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o vmu ./cmd/vmu

FROM jrottenberg/ffmpeg:4.4-alpine AS ffmpeg-builder

# Final with alpine
FROM alpine:latest

# Add labels for better documentation
LABEL maintainer="Brian Jipson <brian.jipson@novelgit.com>"
LABEL description="Video Metadata Updater - Embeds metadata from NFO files into video files"
LABEL version="0.7.0"

# Create a non-root user for running the application
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy application binary to standard location
COPY --from=builder /app/vmu /usr/local/bin/vmu

# Create videos directory as mount point with proper permissions
RUN mkdir -p /videos && chmod 766 /videos

# Copy only the FFmpeg and FFprobe binaries
COPY --from=ffmpeg-builder /usr/local/bin/ffmpeg /usr/local/bin/ffmpeg
COPY --from=ffmpeg-builder /usr/local/bin/ffprobe /usr/local/bin/ffprobe

# Declare the volume mount point
VOLUME ["/videos"]

# Switch to non-root user
USER appuser

# Set entrypoint
ENTRYPOINT ["vmu"]