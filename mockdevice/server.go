package main

import (
    "encoding/json"
    "log"
    "math/rand"
    "net/http"
    "time"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
    })

    mux.HandleFunc("/unstable", func(w http.ResponseWriter, r *http.Request) {
        if rand.Intn(100) < 50 { // 50% failure
            http.Error(w, "device error", http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(map[string]string{"status": "sometimes"})
    })

    srv := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
    rand.Seed(time.Now().UnixNano())
    log.Println("Mock device running on http://localhost:8080")
    log.Fatal(srv.ListenAndServe())
}
