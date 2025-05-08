package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "math/rand"
    "net/http"
    "strings"
    "sync"
    "time"
)

// Result holds the outcome of an operation (ping or fetch).
type Result struct {
    Target string
    Kind   string // "ping" or "fetch"
    OK     bool
    Err    error
}

func main() {
    // CLI flags.
    var (
        urls    string
        posts   int
        workers int
    )
    flag.StringVar(&urls, "urls", "127.0.0.1,192.168.0.1", "comma-separated list of device IPs to ping")
    flag.IntVar(&posts, "posts", 3, "how many JSONPlaceholder posts to fetch concurrently")
    flag.IntVar(&workers, "workers", 4, "goroutines per worker pool")
    flag.Parse()

    ipList := splitCSV(urls)

    rand.Seed(time.Now().UnixNano())

    results := make(chan Result, len(ipList)+posts)
    var wg sync.WaitGroup

    // Ping pool
    wg.Add(1)
    go func() {
        defer wg.Done()
        pool(ipList, workers, func(ip string) {
            ok := mockPing(ip)
            results <- Result{Target: ip, Kind: "ping", OK: ok, Err: nil}
        })
    }()

    // Fetch pool
    wg.Add(1)
    go func() {
        defer wg.Done()
        targets := make([]string, posts)
        for i := 1; i <= posts; i++ {
            targets[i-1] = fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", i)
        }
        pool(targets, workers, func(url string) {
            ok := fetchJSON(url)
            results <- Result{Target: url, Kind: "fetch", OK: ok, Err: nil}
        })
    }()

    // Close results when done.
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect + print summary.
    passed := 0
    total := 0
    for r := range results {
        total++
        status := "FAIL"
        if r.OK {
            status = "OK"
            passed++
        }
        fmt.Printf("%-5s %-40s %s\n", r.Kind, r.Target, status)
    }
    fmt.Printf("\nSummary: %d/%d succeeded\n", passed, total)
}

// pool processes targets with a bounded number of workers.
func pool(targets []string, workers int, fn func(string)) {
    ch := make(chan string)
    var wg sync.WaitGroup
    // Spawn workers.
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for t := range ch {
                fn(t)
            }
        }()
    }
    // Feed targets.
    for _, t := range targets {
        ch <- t
    }
    close(ch)
    wg.Wait()
}

// --- Helpers ---

func splitCSV(s string) []string {
    var out []string
    for _, part := range strings.Split(s, ",") {
        p := strings.TrimSpace(part)
        if p != "" {
            out = append(out, p)
        }
    }
    return out
}

func mockPing(ip string) bool {
    // Simulate latency 50–150 ms.
    time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)
    // 80% success rate.
    return rand.Intn(100) < 80
}

func fetchJSON(url string) bool {
    resp, err := http.Get(url)
    if err != nil {
        return false
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        return false
    }
    var js map[string]any
    return json.NewDecoder(resp.Body).Decode(&js) == nil
}
