package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbUser     = "avnadmin"
	dbPassword = "AVNS_wWoRjEZRmFF5NgjGCcY"
	dbHost     = "mysql-1fb82b3b-boukhar-d756.e.aivencloud.com"
	dbPort     = 20744
	dbName     = "defaultdb"
)

type Contact struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Departure   string    `json:"departure"`
	Destination string    `json:"destination"`
	Number      string    `json:"number"`
	Tour        string    `json:"tour"`
	Comments    string    `json:"comments"`
}

type Reservation struct {
	ID              int       `json:"id"`
	Tour            string    `json:"tour"`
	DateReservation time.Time `json:"date_reservation"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Tel             string    `json:"tel"`
	Transport       string    `json:"transport"`
}

type Review struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Quality   int    `json:"quality"`
	Location  int    `json:"location"`
	Services  int    `json:"services"`
	Team      int    `json:"team"`
	Price     int    `json:"price"`
	Message   string `json:"message"`
	Image     string `json:"image"`
}

type Payment struct {
	ID            int       `json:"id"`
	Amount        float64   `json:"amount"`
	PaymentDate   time.Time `json:"payment_date"`
	ReservationID int       `json:"reservation_id"`
	Status        string    `json:"status"`
}

const AllowedOrigin = "http://localhost:5500"

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

func handleContacts(w http.ResponseWriter, r *http.Request) {
	db, err := getDBConnection()
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT id, first_name, last_name, start_date, end_date, departure, destination, number, tour, comments FROM contactsTours")
		if err != nil {
			http.Error(w, "Error executing query: "+err.Error(), http.StatusInternalServerError)
			log.Println("Query execution error:", err)
			return
		}
		defer rows.Close()

		var contacts []Contact
		for rows.Next() {
			var contact Contact
			if err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.StartDate, &contact.EndDate, &contact.Departure, &contact.Destination, &contact.Number, &contact.Tour, &contact.Comments); err != nil {
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

		stmt, err := db.Prepare("INSERT INTO contactsTours (first_name, last_name, start_date, end_date, departure, destination, number, tour, comments) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			http.Error(w, "Error preparing statement: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error preparing statement:", err)
			return
		}
		defer stmt.Close()

		if _, err = stmt.Exec(contact.FirstName, contact.LastName, contact.StartDate, contact.EndDate, contact.Departure, contact.Destination, contact.Number, contact.Tour, contact.Comments); err != nil {
			http.Error(w, "Error executing statement: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error executing statement:", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Contact added successfully")

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleReservation(w http.ResponseWriter, r *http.Request) {
	db, err := getDBConnection()
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT id, tour, date_reservation, name, email, tel, transport FROM reservations")
		if err != nil {
			http.Error(w, "Error executing query: "+err.Error(), http.StatusInternalServerError)
			log.Println("Query execution error:", err)
			return
		}
		defer rows.Close()

		var reservations []Reservation
		for rows.Next() {
			var (
				id              int
				tour            string
				dateReservation time.Time
				name            string
				email           string
				tel             string
				transport       string
			)

			if err := rows.Scan(&id, &tour, &dateReservation, &name, &email, &tel, &transport); err != nil {
				http.Error(w, "Error reading rows: "+err.Error(), http.StatusInternalServerError)
				log.Println("Error reading rows:", err)
				return
			}

			reservations = append(reservations, Reservation{
				ID:              id,
				Tour:            tour,
				DateReservation: dateReservation,
				Name:            name,
				Email:           email,
				Tel:             tel,
				Transport:       transport,
			})
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating rows: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error iterating rows:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(reservations); err != nil {
			http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error encoding JSON:", err)
		}

	case http.MethodPost:
		var reservation Reservation
		if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
			http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			log.Println("Invalid request payload:", err)
			return
		}

		stmt, err := db.Prepare("INSERT INTO reservations (tour, date_reservation, name, email, tel, transport) VALUES (?, ?, ?, ?, ?, ?)")
		if err != nil {
			http.Error(w, "Error preparing statement: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error preparing statement:", err)
			return
		}
		defer stmt.Close()

		if _, err = stmt.Exec(reservation.Tour, reservation.DateReservation, reservation.Name, reservation.Email, reservation.Tel, reservation.Transport); err != nil {
			http.Error(w, "Error executing statement: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error executing statement:", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Reservation added successfully")

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleReviews(w http.ResponseWriter, r *http.Request) {
	db, err := getDBConnection()
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT id, first_name, last_name, quality, location, services, team, price, message, image FROM reviews")
		if err != nil {
			http.Error(w, "Error executing query: "+err.Error(), http.StatusInternalServerError)
			log.Println("Query execution error:", err)
			return
		}
		defer rows.Close()

		var reviews []Review
		for rows.Next() {
			var review Review
			if err := rows.Scan(&review.ID, &review.FirstName, &review.LastName, &review.Quality, &review.Location, &review.Services, &review.Team, &review.Price, &review.Message, &review.Image); err != nil {
				http.Error(w, "Error reading rows: "+err.Error(), http.StatusInternalServerError)
				log.Println("Error reading rows:", err)
				return
			}
			reviews = append(reviews, review)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating rows: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error iterating rows:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(reviews); err != nil {
			http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error encoding JSON:", err)
		}

	case http.MethodPost:
		var review Review
		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			log.Println("Invalid request payload:", err)
			return
		}

		stmt, err := db.Prepare("INSERT INTO reviews (first_name, last_name, quality, location, services, team, price, message, image) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			http.Error(w, "Error preparing statement: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error preparing statement:", err)
			return
		}
		defer stmt.Close()

		if _, err = stmt.Exec(review.FirstName, review.LastName, review.Quality, review.Location, review.Services, review.Team, review.Price, review.Message, review.Image); err != nil {
			http.Error(w, "Error executing statement: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error executing statement:", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Review added successfully")

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePayments(w http.ResponseWriter, r *http.Request) {
	db, err := getDBConnection()
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT id, amount, payment_date, reservation_id, status FROM payments")
		if err != nil {
			http.Error(w, "Error executing query: "+err.Error(), http.StatusInternalServerError)
			log.Println("Query execution error:", err)
			return
		}
		defer rows.Close()

		var payments []Payment
		for rows.Next() {
			var payment Payment
			if err := rows.Scan(&payment.ID, &payment.Amount, &payment.PaymentDate, &payment.ReservationID, &payment.Status); err != nil {
				http.Error(w, "Error reading rows: "+err.Error(), http.StatusInternalServerError)
				log.Println("Error reading rows:", err)
				return
			}
			payments = append(payments, payment)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Error iterating rows: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error iterating rows:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(payments); err != nil {
			http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error encoding JSON:", err)
		}

	case http.MethodPost:
		var payment Payment
		if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
			http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
			log.Println("Invalid request payload:", err)
			return
		}

		stmt, err := db.Prepare("INSERT INTO payments (amount, payment_date, reservation_id, status) VALUES (?, ?, ?, ?)")
		if err != nil {
			http.Error(w, "Error preparing statement: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error preparing statement:", err)
			return
		}
		defer stmt.Close()

		if _, err = stmt.Exec(payment.Amount, payment.PaymentDate, payment.ReservationID, payment.Status); err != nil {
			http.Error(w, "Error executing statement: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error executing statement:", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Payment added successfully")

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/contacts", handleContacts)
	http.HandleFunc("/reservations", handleReservation)
	http.HandleFunc("/reviews", handleReviews)
	http.HandleFunc("/payments", handlePayments)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
