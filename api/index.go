package handler

import (
    "database/sql"
    "fmt"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
)

// Database connection parameters
const (
    dbUser     = "avnadmin"
    dbPassword = "AVNS_wWoRjEZRmFF5NgjGCcY"
    dbHost     = "mysql-1fb82b3b-boukhar-d756.e.aivencloud.com"
    dbPort     = 20744
    dbName     = "defaultdb"
)

func getDBConnection() (*sql.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
    return sql.Open("mysql", dsn)
}

func Handler(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/":
        fmt.Fprintf(w, "Welcome to the home page!")
    case "/about":
        fmt.Fprintf(w, "This is the about page.")
    case "/getContact":
        db, err := getDBConnection()
        if err != nil {
            http.Error(w, "Database connection error", http.StatusInternalServerError)
            return
        }
        defer db.Close()

        rows, err := db.Query("SELECT * FROM contacts")
        if err != nil {
            http.Error(w, "Error executing query", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var contacts []string
        for rows.Next() {
            var id int
            var name string
            var email string
            if err := rows.Scan(&id, &name, &email); err != nil {
                http.Error(w, "Error reading rows", http.StatusInternalServerError)
                return
            }
            contacts = append(contacts, fmt.Sprintf("ID: %d, Name: %s, Email: %s", id, name, email))
        }

        if err := rows.Err(); err != nil {
            http.Error(w, "Error iterating rows", http.StatusInternalServerError)
            return
        }

        for _, contact := range contacts {
            fmt.Fprintln(w, contact)
        }
    default:
        http.NotFound(w, r)
    }
}
