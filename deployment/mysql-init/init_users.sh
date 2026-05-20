#!/usr/bin/env sh
# NOTE: Use LF line endings. CRLF will break execution inside linux container.
set -eu

log() {
  printf '%s\n' "[init_users] $*"
}

fail() {
  printf '%s\n' "[init_users] ERROR: $*" >&2
  exit 1
}

sql_escape() {
  printf '%s' "$1" | sed "s/'/''/g"
}

validate_identifier() {
  name="$1"
  value="$2"

  case "$value" in
    ''|*[!A-Za-z0-9_]* )
      fail "$name must match [A-Za-z0-9_]+, got '$value'"
      ;;
  esac
}

require_value() {
  name="$1"
  value="$2"

  [ -n "$value" ] || fail "$name is required"
}

APP_DB=${MYSQL_DATABASE:-gochat}
APP_USER=${MYSQL_USER:-gochat_app}
APP_PASS=${MYSQL_PASSWORD:-}
BINLOG_USER=${BINLOG_USER:-}
BINLOG_PASSWORD=${BINLOG_PASSWORD:-}
ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD:-}

require_value MYSQL_ROOT_PASSWORD "$ROOT_PASSWORD"
require_value MYSQL_PASSWORD "$APP_PASS"
validate_identifier MYSQL_DATABASE "$APP_DB"
validate_identifier MYSQL_USER "$APP_USER"

if [ -n "$BINLOG_USER" ] || [ -n "$BINLOG_PASSWORD" ]; then
  require_value BINLOG_USER "$BINLOG_USER"
  require_value BINLOG_PASSWORD "$BINLOG_PASSWORD"
  validate_identifier BINLOG_USER "$BINLOG_USER"
fi

export MYSQL_PWD="$ROOT_PASSWORD"

APP_DB_ESC=$(sql_escape "$APP_DB")
APP_USER_ESC=$(sql_escape "$APP_USER")
APP_PASS_ESC=$(sql_escape "$APP_PASS")

log "creating schema and application user: ${APP_USER}@%"
mysql -uroot <<SQL
SET NAMES utf8mb4;
CREATE DATABASE IF NOT EXISTS \`${APP_DB_ESC}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS '${APP_USER_ESC}'@'%' IDENTIFIED BY '${APP_PASS_ESC}';
GRANT ALL PRIVILEGES ON \`${APP_DB_ESC}\`.* TO '${APP_USER_ESC}'@'%';
FLUSH PRIVILEGES;
SQL

if [ -n "$BINLOG_USER" ]; then
  BINLOG_USER_ESC=$(sql_escape "$BINLOG_USER")
  BINLOG_PASSWORD_ESC=$(sql_escape "$BINLOG_PASSWORD")

  log "creating binlog user: ${BINLOG_USER}@%"
  mysql -uroot <<SQL
SET NAMES utf8mb4;
CREATE USER IF NOT EXISTS '${BINLOG_USER_ESC}'@'%' IDENTIFIED BY '${BINLOG_PASSWORD_ESC}';
GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO '${BINLOG_USER_ESC}'@'%';
GRANT RELOAD, LOCK TABLES, SHOW DATABASES ON *.* TO '${BINLOG_USER_ESC}'@'%';
GRANT SELECT ON \`${APP_DB_ESC}\`.* TO '${BINLOG_USER_ESC}'@'%';
FLUSH PRIVILEGES;
SQL
fi

log "verifying created users"
mysql -uroot -e "SELECT user, host, plugin FROM mysql.user WHERE user IN ('${APP_USER_ESC}'${BINLOG_USER:+, '${BINLOG_USER_ESC}'});"

log "completed user initialization for app='${APP_USER}' binlog='${BINLOG_USER:-none}'"
