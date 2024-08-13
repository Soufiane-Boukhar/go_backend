package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"time"
)

const (
	dbUser     = "avnadmin"
	dbPassword = "AVNS_wWoRjEZRmFF5NgjGCcY"
	dbHost     = "mysql-1fb82b3b-boukhar-d756.e.aivencloud.com"
	dbPort     = 20744
	dbName     = "defaultdb"
)

type Contact struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Departure   string `json:"departure"`
	Destination string `json:"destination"`
	Number      string `json:"number"`
	Tour        string `json:"tour"`
	Comments    string `json:"comments"`
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

func Handler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin != AllowedOrigin && origin != "" {
		http.Error(w, "Origin not allowed", http.StatusForbidden)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", AllowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

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

			stmt, err := db.Prepare("INSERT INTO contactsTours (first_name, last_name, start_date, end_date, departure, destination, number, tour, comments) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				http.Error(w, "Error preparing statement: "+err.Error(), http.StatusInternalServerError)
				log.Println("Error preparing statement:", err)
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(contact.FirstName, contact.LastName, contact.StartDate, contact.EndDate, contact.Departure, contact.Destination, contact.Number, contact.Tour, contact.Comments)
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
	case "/reservation":
		if r.Method == http.MethodGet {
			db, err := getDBConnection()
			if err != nil {
				http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
				log.Println("Database connection error:", err)
				return
			}
			defer db.Close()

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
					dateReservation []byte
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

				dateStr := string(dateReservation)
				dateTime, err := time.Parse("2006-01-02 15:04:05", dateStr)
				if err != nil {
					http.Error(w, "Error parsing date: "+err.Error(), http.StatusInternalServerError)
					log.Println("Error parsing date:", err)
					return
				}

				reservations = append(reservations, Reservation{
					ID:              id,
					Tour:            tour,
					DateReservation: dateTime,
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
		} else if r.Method == http.MethodPost {
			var reservation Reservation
			if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
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

			stmt, err := db.Prepare("INSERT INTO reservations (tour, date_reservation, name, email, tel, transport) VALUES (?, ?, ?, ?, ?, ?)")
			if err != nil {
				http.Error(w, "Error preparing statement: "+err.Error(), http.StatusInternalServerError)
				log.Println("Error preparing statement:", err)
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(reservation.Tour, reservation.DateReservation, reservation.Name, reservation.Email, reservation.Tel, reservation.Transport)
			if err != nil {
				http.Error(w, "Error executing statement: "+err.Error(), http.StatusInternalServerError)
				log.Println("Error executing statement:", err)
				return
			}

			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, "Reservation added successfully")
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/reviews":
		if r.Method == http.MethodGet {
			db, err := getDBConnection()
			if err != nil {
				http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
				log.Println("Database connection error:", err)
				return
			}
			defer db.Close()

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
		} else if r.Method == http.MethodPost {
			var review Review
			if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
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

			stmt, err := db.Prepare("INSERT INTO reviews (first_name, last_name, quality, location, services, team, price, message, image) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				http.Error(w, "Error preparing statement: "+err.Error(), http.StatusInternalServerError)
				log.Println("Error preparing statement:", err)
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(review.FirstName, review.LastName, review.Quality, review.Location, review.Services, review.Team, review.Price, review.Message, review.Image)
			if err != nil {
				http.Error(w, "Error executing statement: "+err.Error(), http.StatusInternalServerError)
				log.Println("Error executing statement:", err)
				return
			}

			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, "Review added successfully")
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/payment":
		if r.Method == http.MethodGet {
			db, err := getDBConnection()
			if err != nil {
				http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
				log.Println("Database connection error:", err)
				return
			}
			defer db.Close()

			rows, err := db.Query("SELECT id, amount, payment_date, reservation_id, status FROM payments")
			if err != nil {
				http.Error(w, "Error executing query: "+err.Error(), http.StatusInternalServerError)
				log.Println("Query execution error:", err)
				return
			}
			defer rows.Close()

			var payments []Payment
			for rows.Next() {
				var (
					id            int
					amount        float64
					paymentDate   []byte
					reservationID int
					status        string
				)

				if err := rows.Scan(&id, &amount, &paymentDate, &reservationID, &status); err != nil {
					http.Error(w, "Error reading rows: "+err.Error(), http.StatusInternalServerError)
					log.Println("Error reading rows:", err)
					return
				}

				dateStr := string(paymentDate)
				paymentTime, err := time.Parse("2006-01-02 15:04:05", dateStr)
				if err != nil {
					http.Error(w, "Error parsing date: "+err.Error(), http.StatusInternalServerError)
					log.Println("Error parsing date:", err)
					return
				}

				payments = append(payments, Payment{
					ID:            id,
					Amount:        amount,
					PaymentDate:   paymentTime,
					ReservationID: reservationID,
					Status:        status,
				})
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
		} else if r.Method == http.MethodPost {
			var payment Payment
			if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
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

			stmt, err := db.Prepare("INSERT INTO payments (amount, payment_date, reservation_id, status) VALUES (?, ?, ?, ?)")
			if err != nil {
				http.Error(w, "Error preparing statement: "+err.Error(), http.StatusInternalServerError)
				log.Println("Error preparing statement:", err)
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(payment.Amount, payment.PaymentDate, payment.ReservationID, payment.Status)
			if err != nil {
				http.Error(w, "Error executing statement: "+err.Error(), http.StatusInternalServerError)
				log.Println("Error executing statement:", err)
				return
			}

			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, "Payment added successfully")
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.NotFound(w, r)
	}
}
