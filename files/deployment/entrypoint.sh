#!/bin/sh
# Entrypoint script for application
set -e

# source ./files/deployment/dbmate.sh
# dbmate --migrations-dir $DBMATE_DEFAULT_MIGRATIONS_DIR --migrations-table $DBMATE_DEFAULT_MIGRATIONS_TABLE migrate
# if [ "$AUTO_MIGRATE_POSTDEPLOYMENT_MIGRATIONS" = "true" ]; then
#     dbmate --migrations-dir $DBMATE_POSTDEPLOYMENT_MIGRATIONS_DIR --migrations-table $DBMATE_POSTDEPLOYMENT_MIGRATIONS_TABLE migrate
# fi

# Debugging: Check if the file exists and print out the path
if [ -f "/app/app" ]; then
    echo "Binary /app/app found!"
else
    echo "Binary /app/app not found!"
fi

# List contents of /app to confirm the binary's location
echo "Listing contents of /app:"
ls -al /app

exec "$@"
