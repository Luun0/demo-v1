#!/bin/sh
set -eu

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
  -v app_password="$DB_APP_PASSWORD" <<'SQL'
SELECT format('CREATE ROLE ssp_app LOGIN PASSWORD %L', :'app_password') \gexec
SQL
