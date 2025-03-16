#!/usr/bin/env bash

set -o errexit

set -o pipefail

set -o nounset

echo "Starting database sync..."

# Export environment variables
export PGPASSWORD=$REMOTE_DB_PASS

# Dump the remote database (without schema changes)
pg_dump -h $REMOTE_DB_HOST -p $REMOTE_DB_PORT -U $REMOTE_DB_USER -d $REMOTE_DB_NAME -a --column-inserts --data-only > /backups/staging_data.sql

echo "Remote database dump completed."

# Restore to local database
export PGPASSWORD=$LOCAL_DB_PASS
psql -h $LOCAL_DB_HOST -p $LOCAL_DB_PORT -U $LOCAL_DB_USER -d $LOCAL_DB_NAME -f /backups/staging_data.sql

echo "Database sync completed!"