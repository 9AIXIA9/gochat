# Database Migrations (Goose)

We use [pressly/goose](https://github.com/pressly/goose) to manage MySQL schema versions.

## Layout
- Migrations folder: `db/migrations`
- Files: SQL migrations (preferred). Use names like `YYYYMMDDHHMMSS_add_users.sql` or sequential `0002_add_users.sql`.

## DSN
Goose expects a DSN. In docker-compose, MySQL is accessible at:
```
mysql://gochat_app:gochat@tcp(mysql:3306)/gochat?parseTime=true&charset=utf8mb4
```
Alternatively, use environment variables from `config/service/.env.app`.

## Commands (Makefile)
- `make migrate-status` — show migration status
- `make migrate-up` — apply pending migrations
- `make migrate-down` — roll back the last migration
- `make migrate-reset` — roll back all migrations
- `make migrate-create name=add_table` — create a new SQL migration stub

## Create a migration
```
make migrate-create name=add_users
```
Then edit the generated SQL file under `db/migrations`.

## Running with Docker
Ensure services are up:
```
make up
```
Then run migrations:
```
make migrate-up
```

## Notes
- Keep migrations idempotent and reversible.
- Prefer SQL migrations for portability.
- Goose stores its own schema version in table `goose_db_version`.

