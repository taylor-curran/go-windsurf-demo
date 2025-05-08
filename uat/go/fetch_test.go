package uat

import (
    "bytes"
    "os/exec"
    "regexp"
    "strings"
    "testing"
)

// TestFetchOnlySummary runs the demo in "fetch-only" mode (no device IPs,
// two JSON posts). It validates that the summary line reports 2 total
// operations.
func TestFetchOnlySummary(t *testing.T) {
    cmd := exec.Command("go", "run", "../../cmd/demo", "-urls", "", "-posts", "2", "-workers", "2")
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
    if len(lines) < 3 { // 2 fetch lines + summary line at minimum
        t.Fatalf("unexpectedly short output (%d lines):\n%s", len(lines), output)
    }

    summary := lines[len(lines)-1]
    re := regexp.MustCompile(`Summary:\s+(\d+)/(\d+)\s+succeeded`)
    matches := re.FindStringSubmatch(summary)
    if matches == nil {
        t.Fatalf("summary line not found or malformed:\n%s", output)
    }

    total := matches[2]
    if total != "2" {
        t.Fatalf("expected total=2 in summary, got %s. Full output:\n%s", total, output)
    }
}
