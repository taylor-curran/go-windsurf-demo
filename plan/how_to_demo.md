# How to Demo the Project

## Prerequisites
* Go ≥ 1.22 in `$PATH`
* Python 3 with `pytest` installed (`pip install pytest`)

## Quick Steps

### 1. Build & run the demo binary
```bash
go build -o demo ./cmd/demo
./demo                       # or: go run ./cmd/demo
```

### 2. Start the mock device server (separate terminal)
```bash
go run ./mockdevice          # listens on http://localhost:8080
```

### 3. Run UAT tests
```bash
go test ./...                # Go tests
pytest -q uat/python         # Python tests
```

---

### Optional Flags for the binary
```
-urls     comma-separated device IPs  (default "127.0.0.1,192.168.0.1")
-posts    number of JSONPlaceholder posts to fetch (default 3)
-workers  goroutines per worker pool  (default 4)
```
Example:
```bash
./demo -urls="10.0.0.1,10.0.0.2" -posts=5 -workers=8
```

## Expected Output
A sample run prints something like:
```
ping  127.0.0.1                           OK
ping  192.168.0.1                         FAIL
fetch https://jsonplaceholder.typicode.com/posts/1 OK
...

Summary: 5/6 succeeded
```

---

## CI / Further Steps
A simple GitHub Actions workflow can be added to run `go test ./...` and `pytest -q uat/python`; let Cascade know whenever you’d like that.
