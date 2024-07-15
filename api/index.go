package handler

import (
    "database/sql"
    "encoding/json"
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

// Contact represents the structure of a contact in the database
type Contact struct {
    ID        int    `json:"id"`
    Email     string `json:"email"`
    Message   string `json:"message"`
    Subject   string `json:"subject"`
    FullName  string `json:"full_name"`
    Tel       string `json:"tel"`
}

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
        fmt.Fprintln(w, "Welcome to the home page!")
    case "/about":
        fmt.Fprintln(w, "This is the about page.")
    case "/contacts":
        if r.Method == http.MethodGet {
            db, err := getDBConnection()
            if err != nil {
                http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
                log.Println("Database connection error:", err)
                return
            }
            defer db.Close()

            rows, err := db.Query("SELECT id, email, message, subject, full_name, tel FROM contacts")
            if err != nil {
                http.Error(w, "Error executing query: "+err.Error(), http.StatusInternalServerError)
                log.Println("Query execution error:", err)
                return
            }
            defer rows.Close()

            var contacts []Contact
            for rows.Next() {
                var contact Contact
                if err := rows.Scan(&contact.ID, &contact.Email, &contact.Message, &contact.Subject, &contact.FullName, &contact.Tel); err != nil {
                    http.Error(w, "Error reading rows: "+err.Error(), http.StatusInternalServerError)
                    log.Println("Error reading rows:", err)
                    return
                }
                contacts = append(contacts, contact)
            }

            if err := rows.Err(); err != nil {
                http.Error(w, "Error iterating rows: "+err.Error(), http.StatusInternalServerError)
                log.Println("Error iterating rows:", err)
                return
            }

            w.Header().Set("Content-Type", "application/json")
            if err := json.NewEncoder(w).Encode(contacts); err != nil {
                http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
                log.Println("Error encoding JSON:", err)
            }
        } else if r.Method == http.MethodPost {
            var contact Contact
            if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
                http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
                log.Println("Invalid request payload:", err)
                return
            }

            db, err := getDBConnection()
            if err != nil {
                http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
                log.Println("Database connection error:", err)
                return
            }
            defer db.Close()

            stmt, err := db.Prepare("INSERT INTO contacts (email, message, subject, full_name, tel) VALUES (?, ?, ?, ?, ?)")
            if err != nil {
                http.Error(w, "Error preparing statement: "+err.Error(), http.StatusInternalServerError)
                log.Println("Error preparing statement:", err)
                return
            }
            defer stmt.Close()

            _, err = stmt.Exec(contact.Email, contact.Message, contact.Subject, contact.FullName, contact.Tel)
            if err != nil {
                http.Error(w, "Error executing statement: "+err.Error(), http.StatusInternalServerError)
                log.Println("Error executing statement:", err)
                return
            }

            w.WriteHeader(http.StatusCreated)
            fmt.Fprintln(w, "Contact added successfully")
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    default:
        http.NotFound(w, r)
    }
}
