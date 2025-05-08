# Cascade + Samsung Demo ‑ High-Level Plan

## 1  Goals
* Showcase Cascade’s ability to:
  1. Generate a **basic Go program** that demonstrates lightweight concurrency.
  2. Write and later refactor **UAT-style tests** (both Go & Python) against that program.
* Simulate “device” interactions via a minimal HTTP server—no real hardware required.
* Keep everything small, readable, and industry-standard so a Python-first audience can follow.

---

## 2  Repository Layout (proposed)

```
.
├── cmd/
│   └── demo/         # `main.go` – entry point
├── internal/
│   ├── pinger/       # concurrent ICMP-like pings (mocked)
│   └── fetcher/      # concurrent JSON downloads
├── mockdevice/       # tiny HTTP server to emulate a device (success + failure modes)
├── uat/
│   ├── go/           # Go test files (using `testing` + `httptest`)
│   └── python/       # Python tests (pytest)
├── plan/
│   ├── demo.md        # high-level plan
│   └── how_to_demo.md # run instructions
├── Makefile          # common tasks (build, run, test)
└── go.mod            # Go ≥1.22
```

Minimalism: no third-party libs unless unavoidable.

---

## 3  Feature Walk-through

### 3.1 Go Demo (`cmd/demo`)
* Accepts a list of target **device URLs/IPs** via CLI flags.
* Spins up **two goroutine pools**:
  * **Ping workers**  
    – For each IP, perform a dummy “ping” (simulated latency + random success/failure).  
  * **Fetch workers**  
    – Downloads JSON from example endpoints (`https://jsonplaceholder.typicode.com/posts/{id}`), again mixing success/failure cases.
* Aggregates results; prints a small summary table.

Key concurrency pieces:
* `sync.WaitGroup`
* Bounded worker pool via buffered channels (idiomatic & easy to grok).

### 3.2 Mock Device (`mockdevice`)
* Runs on `localhost:<port>`.
* Two endpoints:
  * `/health` → returns 200 with JSON `{status:"ok"}`.
  * `/unstable` → randomly returns 500 to emulate failure.

### 3.3 UAT Test Suites
| Suite | Purpose | Highlights |
|-------|---------|------------|
| `uat/go/*.go` | Hits the Go binary as a black-box (exec + `httptest`) | Shows Go’s standard `testing` patterns |
| `uat/python/test_*.py` | Same checks using **pytest** | Lets Samsung teams see parity between Go & Python |

Success + failure assertions included.

---

## 4  Cascade Demo Scenarios

1. **Initial generation**  
   Cascade writes app + tests; all green.
2. **Requirement change** – e.g., new endpoint `/metrics` or altered JSON schema.  
   Existing tests break → Cascade updates code *and* fixes both UAT suites.

We can script this live or record beforehand.

---

## 5  Tooling & Commands

* **Go version:** `>=1.22` (widely available).
* **Makefile targets**
  * `make build` – compile binary
  * `make run` – run demo with default args
  * `make test-go` – run Go UATs
  * `make test-py` – run Python UATs (requires `pytest`)
  * `make device` – start mock device server
* Optional GitHub Actions workflow later (simple matrix: Go tests + Py tests).

---

## 6  Next Steps

1. Initialize repo with the layout above (`go mod init`, etc.).
2. Implement `mockdevice` server.
3. Implement `pinger` and `fetcher` packages.
4. Wire them in `cmd/demo/main.go`.
5. Add initial UAT tests in Go & Python.
6. Prepare “refactor” branch for live Cascade demo.

---

**Let me know if this plan looks good** (or any tweaks), and I’ll start scaffolding the codebase.
