#!/usr/bin/env bash

set -euo pipefail

until psql -h localhost -U "$POSTGRES_USER" -c '\q' -p "$POSTGRES_PORT"; do
  >&2 echo "postgres is unavailable - sleeping ..."
  sleep 1
done
