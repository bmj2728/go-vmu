#!/bin/sh

# Default PUID and PGID if not set
PUID=${PUID:-1000}
PGID=${PGID:-1000}

echo "Starting with PUID: $PUID, PGID: $PGID"

# Get the group name for the specified GID, or create one
if ! getent group $PGID >/dev/null; then
    echo "Creating group with GID $PGID"
    addgroup -S -g "$PGID" appgroup
    GROUP_NAME="appgroup"
else
    GROUP_NAME=$(getent group $PGID | cut -d: -f1)
    echo "Group with GID $PGID already exists: $GROUP_NAME"
fi

# Create user with specified UID if it doesn't exist
if ! getent passwd $PUID >/dev/null; then
    echo "Creating user with UID $PUID in group $GROUP_NAME"
    adduser -S -u "$PUID" -G "$GROUP_NAME" appuser
else
    echo "User with UID $PUID already exists"
fi

# Ensure proper permissions for the /videos mount point
# This is crucial for new files created by the application
chown -R "$PUID":"$PGID" /videos

# Debug: Show what's in the /videos directory
echo "Contents of /videos directory:"
ls -la /videos
echo "Directory permissions:"
ls -ld /videos

# Execute the main application with the correct user
# Use the actual PUID instead of assuming 'appuser' exists
exec su-exec "$PUID":"$PGID" /usr/local/bin/vmu "$@"