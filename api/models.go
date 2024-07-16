package handler

import "time"

type Contact struct {
    ID        int    `json:"id"`
    Email     string `json:"email"`
    Message   string `json:"message"`
    Subject   string `json:"subject"`
    FullName  string `json:"full_name"`
    Tel       string `json:"tel"`
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
    Rating    string `json:"rating"`
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
