package main

import (
    "fmt"
    "net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

func main() {
    http.HandleFunc("/", Handler)
    fmt.Println("Server is running on localhost:3000")
    http.ListenAndServe(":3000", nil)
}