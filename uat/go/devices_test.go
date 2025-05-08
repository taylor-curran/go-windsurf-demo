package uat

import (
    "bytes"
    "os/exec"
    "regexp"
    "strings"
    "testing"
)

// TestPingOnlySummary runs the demo in "ping-only" mode (posts=0) and
// checks that the summary line reports the correct total number of
// operations for the supplied device list. We don't assert how many
// passed versus failed because the mockPing() helper in the demo code
// is deliberately non-deterministic (80 % success rate).
func TestPingOnlySummary(t *testing.T) {
    targets := "10.0.0.1,10.0.0.2" // two dummy devices

    cmd := exec.Command("go", "run", "../../cmd/demo", "-urls", targets, "-posts", "0", "-workers", "2")
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out

    if err := cmd.Run(); err != nil {
        t.Fatalf("demo failed: %v\n%s", err, out.String())
    }

    output := strings.TrimSpace(out.String())
    if output == "" {
        t.Fatal("expected some output from demo")
    }

    lines := strings.Split(output, "\n")
    if len(lines) < 3 { // 2 ping lines + summary line at minimum
        t.Fatalf("unexpectedly short output (%d lines):\n%s", len(lines), output)
    }

    // Check the summary line format and total count.
    summary := lines[len(lines)-1]
    re := regexp.MustCompile(`Summary:\s+(\d+)/(\d+)\s+succeeded`)
    matches := re.FindStringSubmatch(summary)
    if matches == nil {
        t.Fatalf("summary line not found or malformed:\n%s", output)
    }

    total := matches[2]
    if total != "2" { // we passed two device IPs
        t.Fatalf("expected total=2 in summary, got %s. Full output:\n%s", total, output)
    }
}
