# GitStat — Agent Guide

## Dev
- Backend: `cd backend && go run .` — API-only on :12580
- Frontend: `cd frontend && bun run dev` — Vite proxies `/api` → `http://localhost:12580`
- Version sync: `bun run sync-version` — reads `git describe --tags`

## Test
- `cd backend && go test ./...` — no frontend tests, no lint

## Package
- **Bun** only.

## Architecture
- Single binary: Go + `//go:embed web/dist/*`
- Backend: Go 1.26, chi router, 6 internal packages
- Frontend: Vue 3 + Vite + ECharts, no Vue Router/Pinia, custom i18n (`src/i18n.js`)
- API: REST under `/api/*`, query params: `startDate`, `endDate`, `range`, `repos`

## Gotchas
- No CI. `.gitignore` has `*.exe` (covers binary).
- Lazy loading in `handler/lazy_load.go` — incremental forward/backward git log scan.

## Frontend Conventions

- **API**: `client.js` helpers only. POST body as plain object (auto-serialized). No `JSON.stringify` / manual headers.
- **i18n**: All UI text via `t()`. No hardcoded Chinese. Missing keys warn in DEV.
- **Imports**: Named imports only. No `import * as` for local modules. No `.js` extension.
- **CSS**: `--color-*` / `--font-*` vars from `:root`. Global keyframes in `style.css`. Icons via CSS unicode (`content`), never `v-html`.
- **Cleanup**: `onUnmounted` for timers, resize listeners, AbortController. Toast instead of `alert()`. Global `errorHandler` in `main.js`.
- **Data**: `LRUCache(50)` for caches. `v-for :key` uses stable ID, never index.
- **Security**: `rel="noopener noreferrer"` on all `target="_blank"` links.
