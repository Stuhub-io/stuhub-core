#!/usr/bin/env bash

set -o errexit

set -o pipefail

set -o nounset

if [ -z "${DB_USER}" ]; then
  base_postgres_image_default_user='postgres'
  export DB_USER="${base_postgres_image_default_user}"
fi

DBUrl="postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"
# Construct the database URL

echo "${DBUrl}"

/wait-for-it.sh "${DB_HOST}:${DB_PORT}" --timeout=100
>&2 echo "DB is available"

/wait-for-it.sh "${REDIS_HOST}:${REDIS_PORT}" --timeout=100
>&2 echo "REDIS is available"

exec "$@"