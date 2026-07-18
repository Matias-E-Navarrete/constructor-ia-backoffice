# rimu-like app

A personal system combining Hábitos, Entrenamientos, Finanzas and an
Obsidian-style Notes vault, with family/coaching multi-user sharing and an
internal admin platform with feature kill switches.

- `backend/` — Go API, Domain-Driven Design, PostgreSQL. See `backend/README.md`.
- `frontend/` — React + Vite + TypeScript SPA, Tailwind CSS.

## Quick start

```bash
# 1. Postgres
createdb rimu

# 2. Backend
cd backend
cp .env.example .env   # edit DATABASE_URL/JWT_SECRET/ADMIN_EMAILS as needed
go run ./cmd/api        # migrations run automatically, listens on :8080

# 3. Frontend (new terminal)
cd frontend
npm install
npm run dev             # proxies /api to :8080, serves on :5173
```

Register the first account with the email you set in `ADMIN_EMAILS` to get
admin access to `/admin` (roadmap, feature flags, users, test runner).

## Testing

- Backend e2e: `cd backend && go test ./... -p 1` against a `rimu_test`
  database (see `backend/README.md`).
- Frontend e2e: `cd frontend && npx playwright test` — boots the real Go
  backend against `rimu_test` and the built frontend, and drives both with
  Chromium.
