#!/bin/bash
set -e

# Ensure required environment variables are set
if [[ -z "$SCYLLA_KEYSPACE" || -z "$SCYLLA_USER" || -z "$SCYLLA_PASSWORD" ]]; then
  echo "Error: Missing required environment variables SCYLLA_KEYSPACE, SCYLLA_USER, or SCYLLA_PASSWORD."
  exit 1
fi

# Generate CQL script dynamically
cat <<EOF > temp_init.cql
-- Initialize the keyspace
CREATE KEYSPACE IF NOT EXISTS $SCYLLA_KEYSPACE
WITH REPLICATION = {
    'class': 'SimpleStrategy',
    'replication_factor': 1
}
AND durable_writes = true;
EOF

# Run the CQL script
cqlsh $SCYLLA_HOST 9042 -u cassandra -p cassandra -f temp_init.cql


echo "Keyspace $SCYLLA_KEYSPACE created successfully."
# Clean up
rm temp_init.cql
