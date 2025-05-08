package uat

import (
    "bytes"
    "fmt"
    "os/exec"
    "strings"
    "testing"
    "time"
)

// TestWorkerPoolConcurrency ensures that using multiple workers finishes
// significantly faster than a single worker. Because the ping helper in
// the demo sleeps 50–150 ms, running 20 pings serially should take roughly
// ≥1 s whereas with 10 workers it should drop well below half that time.
// We keep the assertion loose (factor 2×) to avoid flakiness.
func TestWorkerPoolConcurrency(t *testing.T) {
    // Build 20 dummy IPs.
    var ips []string
    for i := 1; i <= 20; i++ {
        ips = append(ips, fmt.Sprintf("10.0.0.%d", i))
    }
    ipList := strings.Join(ips, ",")

    // Helper to run demo and return duration + output.
    runDemo := func(workers string) (time.Duration, string, error) {
        cmd := exec.Command("go", "run", "../../cmd/demo", "-urls", ipList, "-posts", "0", "-workers", workers)
        var out bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &out
        start := time.Now()
        err := cmd.Run()
        return time.Since(start), out.String(), err
    }

    serialDur, serialOut, err := runDemo("1")
    if err != nil {
        t.Fatalf("serial run failed: %v\n%s", err, serialOut)
    }

    concurrDur, concurrOut, err := runDemo("10")
    if err != nil {
        t.Fatalf("concurrent run failed: %v\n%s", err, concurrOut)
    }

    // Expect at least 2× speed-up.
    if concurrDur >= serialDur/2 {
        t.Fatalf("expected concurrency speed-up ≥2×\nserial=%v\nconcurrent=%v", serialDur, concurrDur)
    }
}
