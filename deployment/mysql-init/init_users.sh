#!/usr/bin/env sh
# NOTE: Use LF line endings. CRLF will break execution inside linux container.
set -euo pipefail

# This script runs during first-time MySQL initialization inside the container.
# It creates the application user and (optionally) the binlog CDC user if env variables are present.
# Env variables consumed (provided via docker-compose env_file .env.db):
#   MYSQL_ROOT_PASSWORD
#   MYSQL_DATABASE
#   MYSQL_USER / MYSQL_PASSWORD (app user)
#   BINLOG_USER / BINLOG_PASSWORD (binlog user)
# Safe to re-run: uses CREATE USER IF NOT EXISTS and idempotent GRANTs.

APP_DB=${MYSQL_DATABASE:-gochat}
APP_USER=${MYSQL_USER:-gochat_app}
APP_PASS=${MYSQL_PASSWORD:-change_me_app}
BINLOG_USER=${BINLOG_USER:-}
BINLOG_PASSWORD=${BINLOG_PASSWORD:-}

mysql=( mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" )

# Create application user & grant privileges on its schema
"${mysql[@]}" <<SQL
CREATE USER IF NOT EXISTS '${APP_USER}'@'%' IDENTIFIED WITH mysql_native_password BY '${APP_PASS}';
GRANT ALL PRIVILEGES ON \`${APP_DB}\`.* TO '${APP_USER}'@'%';
FLUSH PRIVILEGES;
SQL

# Create binlog user only if both env vars are non-empty
if [ -n "${BINLOG_USER}" ] && [ -n "${BINLOG_PASSWORD}" ]; then
  "${mysql[@]}" <<SQL
CREATE USER IF NOT EXISTS '${BINLOG_USER}'@'%' IDENTIFIED WITH mysql_native_password BY '${BINLOG_PASSWORD}';
-- Replication related privileges (for row-based binlog consumption)
GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO '${BINLOG_USER}'@'%';
-- Needed for consistent snapshot / listing databases
GRANT RELOAD, LOCK TABLES, SHOW DATABASES ON *.* TO '${BINLOG_USER}'@'%';
-- Read access to application database (initial snapshot)
GRANT SELECT ON \`${APP_DB}\`.* TO '${BINLOG_USER}'@'%';
FLUSH PRIVILEGES;
SQL
fi

# Log created users (host only) for debugging
"${mysql[@]}" -e "SELECT user, host, plugin FROM mysql.user WHERE user IN ('${APP_USER}', '${BINLOG_USER}');"

echo "[init_users] Completed user initialization for app='${APP_USER}' binlog='${BINLOG_USER}'"
