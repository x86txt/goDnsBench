## goDnsBench Multi-Agent Development Plan

### Goals (next milestone)

- **GUI MVP**: replace frontend mocks with real Wails bindings; run benchmark; show progress/results/chart; import/export.
- **Single binary**: `goDnsBench` launches GUI by default; `--tui` runs the terminal UI.

### Working agreements (avoid conflicts)

- **Single-writer files** (one agent at a time): `app.go`, `main.go`, `gui.go`, `Makefile`, `cmd/goDnsBench/main.go`.
- **Directory boundaries**:
- Frontend agents: `frontend/src/**` only
- Backend agents: `internal/**` (+ `app.go`/`gui.go`)
- TUI agents: `internal/tui/**`
- **Prefer additive changes**: create new helpers/modules instead of rewriting shared files.
- **Integration**: one Integrator merges PRs at the end of each stage and resolves conflicts.

### Stage 0 - Contracts + seams (parallel)

**Definition of done**: FE has stable API surface; frontend code is modular enough for parallel FE work.

#### Agent A (Backend contract + event names)

- **Owns**: `app.go` (+ new `internal/*` helper(s) if needed)
- **Tasks**:
- Define GUI-safe DTOs for **Results**, **Server**, **Settings**, **Progress** (durations as ms numbers; timestamps as ISO strings).
- Decide event names and payload shape (example: `benchmark:progress`, `benchmark:done`, `benchmark:error`).
- Add conversion helpers (internal structs -> DTO).

#### Agent B (Frontend refactor for parallelism)

- **Owns**: `frontend/src/scripts/**`
- **Tasks**:
- Split `frontend/src/scripts/app.ts` into modules (example: `api/*`, `ui/*`, `charts/*`).
- Keep a single `init()` entrypoint that wires events + UI.

#### Agent C (Docs + dev UX)

- **Owns**: `README.md` (and `CONTEXT.md` if needed)
- **Tasks**:
- Update dev/build docs to match unified binary behavior and GUI dev flow (Bun + Wails).

### Stage 1 - GUI MVP (parallel)

**Definition of done**: GUI runs real benchmark (no mocks) and displays table+chart; import/export works.

#### Agent A (Backend: selection + progress)

- **Owns**: `app.go`, `internal/benchmark/runner.go` (only if needed)
- **Tasks**:
- Make benchmark run respect GUI server selection (not only `a.servers`).
- Wire `internal/benchmark/runner.go` `ProgressCallback` to Wails event emission.
- Ensure `GetResults()` returns DTOs.

#### Agent B (Frontend: real bindings + results/charts)

- **Owns**: `frontend/src/scripts/api/**`, `frontend/src/scripts/charts/**`
- **Tasks**:
- Replace `mock*()` calls with real `wailsjs` bindings.
- Render results table and Chart.js from DTO results.

#### Agent C (Frontend: dialogs + settings UI)

- **Owns**: `frontend/src/scripts/ui/**`, `frontend/src/pages/index.astro`
- **Tasks**:
- Implement import/export dialogs and wire to `LoadServersFromFile`, `ExportResultsCSV`, `ExportResultsJSON`.
- Implement minimal settings UI (timeout/concurrency/protocols/domains) backed by backend settings.

#### Agent D (Build: deterministic Wails bindings)

- **Owns**: `wails.json`, `frontend/**`
- **Tasks**:
- Decide whether `frontend/src/wailsjs` is committed or generated; document the rule.
- Ensure Astro build output matches Wails embed path (`frontend/dist`).

### Stage 2 - Unify entrypoints + Makefile (sequential gate)

**Definition of done**: there is one canonical entrypoint and one primary artifact.

#### Agent D (Integrator)

- **Owns**: `main.go`, `cmd/goDnsBench/main.go`, `Makefile`
- **Tasks**:
- Remove/redirect the stale CLI entrypoint (`cmd/goDnsBench/main.go` currently prints GUI not implemented).
- Update build targets so the unified `goDnsBench` binary is the primary artifact.
- Confirm flags are consistent across GUI and TUI paths.

### Stage 3 - TUI + headless exports (parallel)

**Definition of done**: `--tui` runs a real benchmark; export flags work headlessly.

#### Agent E (TUI)

- **Owns**: `internal/tui/app.go`
- **Tasks**:
- Implement benchmark run + progress view + results rendering (replace current stubs).

#### Agent F (CLI exports)

- **Owns**: new helper file(s) + minimal `main.go` changes (via PR reviewed by Agent D)
- **Tasks**:
- Implement `--export-json` and `--export-csv` using `internal/export/*`.

### Stage 4 - Quality + release (parallel)

**Definition of done**: CI green; basic tests; reproducible builds.

#### Agent G (Tests)

- Add unit tests for `internal/benchmark`, `internal/config`, `internal/export`.

#### Agent H (CI)

- Add GitHub Actions for `go test ./...`, frontend build, and Wails build smoke checks.

### Optional Stage 5 - Enhancements (parallel)

- GUI filtering/sorting, history, richer charts, protocol/IPv6 hardening.