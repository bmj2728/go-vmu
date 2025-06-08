#!/bin/sh

# Default PUID and PGID if not set
PUID=${PUID:-1000}
PGID=${PGID:-1000}

echo "Starting with PUID: $PUID, PGID: $PGID"

# Check if group 'appgroup' exists and its GID matches PGID
# If not, try to create/modify it
if ! getent group $PGID >/dev/null; then
    echo "Group with GID $PGID does not exist. Attempting to add group 'appgroup' with GID $PGID."
    addgroup -S -g "$PGID" appgroup || echo "Warning: Could not add group 'appgroup' with GID $PGID. It might already exist with a different name."
elif [ "$(getent group $PGID | cut -d: -f1)" != "appgroup" ]; then
    echo "GID $PGID is used by another group. Attempting to modify group 'appgroup' GID to $PGID."
    groupmod -g "$PGID" appgroup || echo "Warning: Could not modify group 'appgroup' GID. Group might be in use or name conflict."
fi

# Check if user 'appuser' exists and its UID matches PUID
# If not, try to create/modify it
if ! getent passwd $PUID >/dev/null; then
    echo "User with UID $PUID does not exist. Attempting to add user 'appuser' with UID $PUID."
    adduser -S -u "$PUID" -G appgroup appuser || echo "Warning: Could not add user 'appuser' with UID $PUID."
elif [ "$(getent passwd $PUID | cut -d: -f1)" != "appuser" ]; then
    echo "UID $PUID is used by another user. Attempting to modify user 'appuser' UID to $PUID."
    usermod -u "$PUID" appuser || echo "Warning: Could not modify user 'appuser' UID. User might be in use or name conflict."
fi

# Ensure user 'appuser' is part of 'appgroup'
if ! id -nG appuser | grep -qw appgroup; then
    echo "Adding appuser to appgroup if not already a member."
    adduser appuser appgroup || echo "Warning: Could not add appuser to appgroup."
fi

# Ensure proper permissions for the /videos mount point
# This is crucial for new files created by the application
chown -R "$PUID":"$PGID" /videos

# Execute the main application as the appuser, explicitly calling 'vmu'
# "$@" passes all command-line arguments (e.g., /videos --workers 4) to vmu
exec su-exec appuser /usr/local/bin/vmu "$@"