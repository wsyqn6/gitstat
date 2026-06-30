# GitStat — Agent Guide

## Principle
追求代码精简，实现高效，维护方便 — applies to both frontend and backend.

## Dev
- Backend: `cd backend && go run .` — API-only on :12580
- Frontend: `cd frontend && bun run dev` — Vite proxies `/api` → `http://localhost:12580`
- Version sync: `bun run sync-version` — reads `git describe --tags`
- Dev scan path: env var `GITSTAT_DEV_PATH` (default: `D:/work/ems`)

## Git
- Commits via `git-commit` skill.
- Message: single line, concise, conventional commit format.

## Test
- `cd backend && go test ./...` — no frontend tests, no lint

## Package
- **Bun** only.

## Architecture
- Single binary: Go + `//go:embed web/dist/*`
- Backend: Go 1.26, chi router, 6 internal packages
- Frontend: Vue 3 + Vite + ECharts, no Vue Router/Pinia, custom i18n (`src/i18n.js`)
- API: REST under `/api/*`, query params: `startDate`, `endDate`, `range`, `repos`
- **Data computation boundary**: ALL computation → backend. Frontend renders only. Never derive data (同比, 环比, rankings, totals) in JS.
- Pages: 概览(Dashboard) / 趋势(Trends) / 仓库(RepoSection) / 对比(Compare)
  - `Trends.vue` replaces old `Analytics.vue` — time-series charts only, no cross-repo compare
  - `Compare.vue` (future) — cross-repo radar/table comparison

## Gotchas
- No CI. `.gitignore` has `*.exe` (covers binary).
- Lazy loading in `handler/lazy_load.go` — incremental forward/backward git log scan.
- Page rename: old `analytics` localStorage key auto-migrates to `trends` in App.vue

## Frontend Conventions

- **API**: `client.js` helpers only. POST body as plain object (auto-serialized). No `JSON.stringify` / manual headers.
- **i18n**: All UI text via `t()`. No hardcoded Chinese. Missing keys warn in DEV.
- **Imports**: Named imports only. No `import * as` for local modules. No `.js` extension.
- **CSS**: Prefer `var(--color-*)` / `var(--font-*)` over hardcoded values. Global keyframes in `style.css`. Icons via CSS `content`, never `v-html`. Inline styles only for dynamic values. Scoped CSS per component; share only truly reusable styles globally.
- **Cleanup**: `onUnmounted` for timers, resize listeners, AbortController. Toast instead of `alert()`. Global `errorHandler` in `main.js`.
- **Data**: `LRUCache(50)` for caches. `v-for :key` uses stable ID, never index.
- **Architecture & Size**: Start thin, compose early. View's job is data flow + layout; each child component owns one functional unit. Extract >400 line views to `src/components/`. Never write monolithic and split later.
- **Component contract**: Props typed via `defineProps({...})`. All emits declared via `defineEmits`. Remove unused props/imports before commit.
- **Store discipline**: Mutate global state only through store functions. Never `state.x = y` from components.
- **Security**: `rel="noopener noreferrer"` on all `target="_blank"` links.
- **Accessibility**: Icon-only buttons must have `aria-label`. Decorative Unicode icons use `aria-hidden="true"`.

## Go Conventions
- **No duplication**: Extract repeated patterns into helpers. If code block appears in 2+ places, factor it out. This applies to handlers, routing, startup logic, error handling. "light and fast" means small, composable functions.
