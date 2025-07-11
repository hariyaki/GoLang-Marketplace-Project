package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "ok")
    })
    fmt.Println("Server on :8080")
    http.ListenAndServe(":8080", nil)
}
