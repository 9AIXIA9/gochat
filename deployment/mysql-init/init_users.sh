set -eu

# This script runs during first-time MySQL initialization inside the container

APP_DB=${MYSQL_DATABASE:-gochat}
APP_USER=${MYSQL_USER:-gochat_app}
APP_PASS=${MYSQL_PASSWORD:-change_me_app}

BINLOG_USER=${BINLOG_USER:-}
BINLOG_PASSWORD=${BINLOG_PASSWORD:-}

mysql=( mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" )

# App 用户（用于应用连接）
"${mysql[@]}" <<SQL
CREATE USER IF NOT EXISTS '${APP_USER}'@'%' IDENTIFIED WITH mysql_native_password BY '${APP_PASS}';
GRANT ALL PRIVILEGES ON \`${APP_DB}\`.* TO '${APP_USER}'@'%';
FLUSH PRIVILEGES;
SQL

# Binlog 用户（用于 CDC/增量读取）
if [ -n "${BINLOG_USER}" ] && [ -n "${BINLOG_PASSWORD}" ]; then
  "${mysql[@]}" <<SQL
CREATE USER IF NOT EXISTS '${BINLOG_USER}'@'%' IDENTIFIED WITH mysql_native_password BY '${BINLOG_PASSWORD}';

-- 二进制日志与复制相关（MySQL 8.0.22+ 可用 REPLICATION REPLICA；此处使用兼容别名）
GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO '${BINLOG_USER}'@'%';

-- 快照/一致性读（FLUSH TABLES WITH READ LOCK、列出数据库）
GRANT RELOAD, LOCK TABLES, SHOW DATABASES ON *.* TO '${BINLOG_USER}'@'%';

-- 读取业务库内容用于初始快照
GRANT SELECT ON \`${APP_DB}\`.* TO '${BINLOG_USER}'@'%';

FLUSH PRIVILEGES;
SQL
fi

# 可选：查看当前用户（仅用于日志）
"${mysql[@]}" -e "SELECT user, host FROM mysql.user;"
