#!/bin/bash

set -o errexit

set -o pipefail

set -o nounset

# Wait scylla server
/wait-for-it.sh "${SCYLLA_HOST}:${SCYLLA_PORT}" --timeout=100

# Wait for ScyllaDB to be ready
echo ScyllaDB is ready, initializing schema...

exec "$@"