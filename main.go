package main

import (
    "fmt"
    "net/http"
)

// VercelHandler is the exported function that Vercel will use
func VercelHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, test!")
}

// Index is an alternative name that Vercel might look for
func Index(w http.ResponseWriter, r *http.Request) {
    VercelHandler(w, r)
}

// Keep the main function for local development
func main() {
    http.HandleFunc("/", VercelHandler)
    fmt.Println("Server is running on localhost:3000")
    http.ListenAndServe(":3000", nil)
}