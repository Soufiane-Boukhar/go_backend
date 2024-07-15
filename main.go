package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Server is running on localhost:3000")
    http.ListenAndServe(":3000", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}
