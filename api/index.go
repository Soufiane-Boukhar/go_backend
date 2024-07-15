package handler

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
    "regexp"
    "strings"
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
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error connecting to the database: %w", err)
    }
    return db, nil
}

// validateEmail checks if the email address is valid
func validateEmail(email string) bool {
    const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    re := regexp.MustCompile(emailRegex)
    return re.MatchString(email)
}

// validateInput sanitizes and validates the contact input
func validateInput(contact *Contact) error {
    if !validateEmail(contact.Email) {
        return fmt.Errorf("invalid email address")
    }
    // Example: Ensure message length does not exceed 500 characters
    if len(contact.Message) > 500 {
        return fmt.Errorf("message too long")
    }
    // Additional validations can be added here
    return nil
}

// Handler processes HTTP requests and interacts with the database
func Handler(w http.ResponseWriter, r *http.Request) {
    // Set security headers
    w.Header().Set("Content-Security-Policy", "default-src 'self'")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.Header().Set("X-Frame-Options", "DENY")
    w.Header().Set("X-XSS-Protection", "1; mode=block")

    switch r.URL.Path {
    case "/":
        fmt.Fprintln(w, "Welcome to the home page!")
    case "/about":
        fmt.Fprintln(w, "This is the about page.")
    case "/contacts":
        switch r.Method {
        case http.MethodGet:
            db, err := getDBConnection()
            if err != nil {
                http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
                log.Println("Database connection error:", err)
                return
            }
            defer db.Close()

            rows, err := db.Query("SELECT id, email, message, subject, full_name, tel FROM contactsTours")
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
        case http.MethodPost:
            var contact Contact
            if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
                http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
                log.Println("Invalid request payload:", err)
                return
            }

            if err := validateInput(&contact); err != nil {
                http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
                return
            }

            db, err := getDBConnection()
            if err != nil {
                http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
                log.Println("Database connection error:", err)
                return
            }
            defer db.Close()

            stmt, err := db.Prepare("INSERT INTO contactsTours (email, message, subject, full_name, tel) VALUES (?, ?, ?, ?, ?)")
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
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    default:
        http.NotFound(w, r)
    }
}
