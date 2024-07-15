package handler

import (
    "fmt"
    "net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/":
        fmt.Fprintf(w, "Welcome to the home page!")
    case "/about":
        fmt.Fprintf(w, "This is the about page.")
    default:
        http.NotFound(w, r)
    }
}