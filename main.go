package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

// Book a book struct
type Book struct {
	ID    string `json:"id"`
	Isbn  string `json:"isbn"`
	Title string `json:"title"`
}

func getBooksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(books)
	}
}

func createBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var book Book
		_ = json.NewDecoder(r.Body).Decode(&book)
		row := db.QueryRow("INSERT INTO books(isbn, title) VALUES($1, $2) returning id, isbn, title;", book.Isbn, book.Title)
		switch err := row.Scan(&book.ID, &book.Isbn, &book.Title); err {
		case nil:
			json.NewEncoder(w).Encode(book)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func getBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		var book Book
		row := db.QueryRow("SELECT id, isbn, title FROM books WHERE id=$1", params["id"])
		switch err := row.Scan(&book.ID, &book.Isbn, &book.Title); err {
		case nil:
			json.NewEncoder(w).Encode(book)
		case sql.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func updateBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		var book Book
		book.ID = params["id"]
		_ = json.NewDecoder(r.Body).Decode(&book)
		row := db.QueryRow("UPDATE books SET isbn=$2, title=$3 WHERE id = $1 returning id, isbn, title;", book.ID, book.Isbn, book.Title)
		switch err := row.Scan(&book.ID, &book.Isbn, &book.Title); err {
		case nil:
			json.NewEncoder(w).Encode(book)
		case sql.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func deleteBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		switch _, err := db.Exec("DELETE FROM books where id = $1", params["id"]); err {
		case nil:
			w.WriteHeader(http.StatusOK)
		case sql.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func createDB() *sql.DB {
	dbStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", dbStr)

	if err != nil {
		panic(err)
	}

	return db
}

func main() {
	db := createDB()
	router := mux.NewRouter()

	router.HandleFunc("/api/books", getBooksHandler(db)).Methods(http.MethodGet)
	router.HandleFunc("/api/books", createBookHandler(db)).Methods(http.MethodPost)
	router.HandleFunc("/api/books/{id}", getBookHandler(db)).Methods(http.MethodGet)
	router.HandleFunc("/api/books/{id}", updateBookHandler(db)).Methods(http.MethodPut)
	router.HandleFunc("/api/books/{id}", deleteBookHandler(db)).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", router))
}
