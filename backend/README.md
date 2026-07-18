# Rimu backend

Go API following Domain-Driven Design: one bounded context per module
(`user`, `groups`, `habits`, `workouts`, `finance`, `notes`, `admin`), each
with `domain` / `application` / `infrastructure` / `interfaces/http` layers.
See `internal/platform/httpserver/build.go` for the composition root.

## Run

```
cp .env.example .env   # edit as needed
go run ./cmd/api
```

Requires a local PostgreSQL reachable at `DATABASE_URL`; migrations in
`internal/platform/db/migrations` run automatically on boot.

## Test

Every bounded context has an `*_e2e_test.go` file next to its
`interfaces/http` package that drives the real router over HTTP against a
dedicated `rimu_test` database (no mocks).

```
createdb rimu_test   # once
TEST_DATABASE_URL=postgres://rimu:rimu@localhost:5432/rimu_test?sslmode=disable \
  go test ./... -p 1
```

`-p 1` is required: every test package truncates the shared `rimu_test`
database on setup, so packages must run one at a time rather than in
parallel.
