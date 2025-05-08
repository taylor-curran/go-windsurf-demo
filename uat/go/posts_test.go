package uat

import (
	"bytes"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// TestPostsFetching runs the demo with a specified number of posts
// and verifies that the correct number of fetch operations are reported
// in the output. This test has an intentional bug that needs fixing.
func TestPostsFetching(t *testing.T) {
	// Number of posts to fetch
	postCount := 5
	
	// Run the demo with the specified number of posts
	cmd := exec.Command("go", "run", "../../cmd/demo", "-posts", strconv.Itoa(postCount), "-urls", "127.0.0.1")
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
	
	// Count the number of fetch operations in the output
	fetchCount := 0
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "fetch") {
			fetchCount++
		}
	}
	
	// BUG: This check is incorrect - we're checking against postCount-1 instead of postCount
	// This will fail when the actual count is correct
	if fetchCount != postCount-1 {
		t.Fatalf("expected %d fetch operations, got %d\nOutput:\n%s", 
			postCount-1, fetchCount, output)
	}
	
	// Check the summary line format and success count
	summary := lines[len(lines)-1]
	re := regexp.MustCompile(`Summary:\s+(\d+)/(\d+)\s+succeeded`)
	matches := re.FindStringSubmatch(summary)
	if matches == nil {
		t.Fatalf("summary line not found or malformed:\n%s", output)
	}
	
	// Extract total operations count from summary
	total, err := strconv.Atoi(matches[2])
	if err != nil {
		t.Fatalf("could not parse total count: %v", err)
	}
	
	// Verify total operations count (1 ping + postCount fetches)
	expectedTotal := 1 + postCount // One ping operation plus fetch operations
	if total != expectedTotal {
		t.Fatalf("expected total=%d in summary, got %d. Full output:\n%s", 
			expectedTotal, total, output)
	}
}
