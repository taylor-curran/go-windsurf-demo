package uat

import (
    "os/exec"
    "testing"
)

func TestDemoBinary(t *testing.T) {
    cmd := exec.Command("go", "run", "../../cmd/demo")
    out, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("binary failed: %v\n%s", err, out)
    }
    if len(out) == 0 {
        t.Fatal("expected some output")
    }
}
