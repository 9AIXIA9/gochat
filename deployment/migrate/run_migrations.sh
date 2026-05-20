#!/usr/bin/env sh
set -eu

log() {
  printf '%s\n' "[migrate] $*"
}

fail() {
  printf '%s\n' "[migrate] ERROR: $*" >&2
  exit 1
}

require_value() {
  name="$1"
  value="$2"

  [ -n "$value" ] || fail "$name is required"
}

GOOSE_BIN=${GOOSE_BIN:-/usr/local/bin/goose}
MIGRATIONS_DIR=${GOOSE_MIGRATIONS_DIR:-/app/db/migrations}
DB_HOST=${DB_HOST:-mysql}
DB_PORT=${DB_PORT:-3306}
DB_USERNAME=${DB_USERNAME:-gochat_app}
DB_PASSWORD=${DB_PASSWORD:-}
DB_NAME=${DB_NAME:-gochat}

require_value GOOSE_BIN "$GOOSE_BIN"
require_value DB_PASSWORD "$DB_PASSWORD"

[ -x "$GOOSE_BIN" ] || fail "goose binary not found or not executable: $GOOSE_BIN"
[ -d "$MIGRATIONS_DIR" ] || fail "migrations directory not found: $MIGRATIONS_DIR"

DSN="${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=true&charset=utf8mb4"

log "running goose up against ${DB_HOST}:${DB_PORT}/${DB_NAME}"
log "using migrations dir: ${MIGRATIONS_DIR}"
"$GOOSE_BIN" -dir "$MIGRATIONS_DIR" mysql "$DSN" up
log "migration completed successfully"

