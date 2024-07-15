package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

type Contact struct {
    ID       int    `json:"id"`
    Email    string `json:"email"`
    Message  string `json:"message"`
    Subject  string `json:"subject"`
    FullName string `json:"full_name"`
    Tel      string `json:"tel"`
}

func main() {
    dsn := "avnadmin:AVNS_wWoRjEZRmFF5NgjGCcY@tcp(mysql-1fb82b3b-boukhar-d756.e.aivencloud.com:20744)/defaultdb?tls=false"

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to the database successfully!")

    rows, err := db.Query("SELECT * FROM contacts")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var contacts []Contact

    for rows.Next() {
        var c Contact
        if err := rows.Scan(&c.ID, &c.Email, &c.Message, &c.Subject, &c.FullName, &c.Tel); err != nil {
            log.Fatal(err)
        }
        contacts = append(contacts, c)
    }

    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }

    jsonData, err := json.Marshal(contacts)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(jsonData))
}
