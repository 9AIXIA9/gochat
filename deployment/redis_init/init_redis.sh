#!/usr/bin/env sh
set -eu

log() {
  printf '%s\n' "[init_redis] $*"
}

fail() {
  printf '%s\n' "[init_redis] ERROR: $*" >&2
  exit 1
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

hash_password() {
  printf '%s' "$1" | sha256sum | awk '{print $1}'
}

require_value() {
  name="$1"
  value="$2"

  [ -n "$value" ] || fail "$name is required"
}

REDIS_ADMIN_USER=${REDIS_ADMIN_USER:-admin}
REDIS_ADMIN_PASSWORD=${REDIS_ADMIN_PASSWORD:-}
REDIS_USER=${REDIS_USER:-gochat_app}
REDIS_PASSWORD=${REDIS_PASSWORD:-}
ACL_FILE=${REDIS_ACL_FILE:-/data/users.acl}

require_value REDIS_ADMIN_PASSWORD "$REDIS_ADMIN_PASSWORD"
require_value REDIS_PASSWORD "$REDIS_PASSWORD"
validate_identifier REDIS_ADMIN_USER "$REDIS_ADMIN_USER"
validate_identifier REDIS_USER "$REDIS_USER"

ADMIN_HASH=$(hash_password "$REDIS_ADMIN_PASSWORD")
APP_HASH=$(hash_password "$REDIS_PASSWORD")

mkdir -p "$(dirname "$ACL_FILE")"
cat >"$ACL_FILE" <<EOF
user default off
user ${REDIS_ADMIN_USER} on #${ADMIN_HASH} ~* +@all
user ${REDIS_USER} on #${APP_HASH} ~* +@all
EOF
chmod 600 "$ACL_FILE"

log "generated ACL file at ${ACL_FILE}"
log "starting redis-server with generated ACL users"
exec redis-server --appendonly yes --aclfile "$ACL_FILE"

