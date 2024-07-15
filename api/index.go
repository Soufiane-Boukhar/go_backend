package handler

import (
    "database/sql"
    "fmt"
    "log"
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

// getDBConnection establishes a connection to the MySQL database
func getDBConnection() (*sql.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("error opening database connection: %w", err)
    }
    // Check if the connection is valid
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error connecting to the database: %w", err)
    }
    return db, nil
}

// Handler processes HTTP requests and interacts with the database
func Handler(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/":
        fmt.Fprintf(w, "Welcome to the home page!")
    case "/about":
        fmt.Fprintf(w, "This is the about page.")
    case "/getContacts":
        db, err := getDBConnection()
        if err != nil {
            http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
            log.Println("Database connection error:", err)
            return
        }
        defer db.Close()

        rows, err := db.Query("SELECT id, name, email FROM contacts")
        if err != nil {
            http.Error(w, "Error executing query: "+err.Error(), http.StatusInternalServerError)
            log.Println("Query execution error:", err)
            return
        }
        defer rows.Close()

        var contacts []string
        for rows.Next() {
            var id int
            var name string
            var email string
            if err := rows.Scan(&id, &name, &email); err != nil {
                http.Error(w, "Error reading rows: "+err.Error(), http.StatusInternalServerError)
                log.Println("Error reading rows:", err)
                return
            }
            contacts = append(contacts, fmt.Sprintf("ID: %d, Name: %s, Email: %s", id, name, email))
        }

        if err := rows.Err(); err != nil {
            http.Error(w, "Error iterating rows: "+err.Error(), http.StatusInternalServerError)
            log.Println("Error iterating rows:", err)
            return
        }

        for _, contact := range contacts {
            fmt.Fprintln(w, contact)
        }
    default:
        http.NotFound(w, r)
    }
}
