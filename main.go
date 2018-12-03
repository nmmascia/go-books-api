package main

import (
	"database/sql"
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
		// var book Book
		// _ = json.NewDecoder(r.Body).Decode(&book)
		// json.NewEncoder(w).Encode()
	}
}

func getBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// params := mux.Vars(r)
		// json.NewEncoder(w).Encode()
	}
}

func updateBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// params := mux.Vars(r)
		// json.NewEncoder(w).Encode()
	}
}

func deleteBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// params := mux.Vars(r)
		// json.NewEncoder(w).Encode()
	}
}

func createDB() *sql.DB {
	dbStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", dbStr)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	db := createDB()
	router := mux.NewRouter()

	router.HandleFunc("/api/books", getBooksHandler(db)).Methods(http.MethodGet)
	router.HandleFunc("/api/books", createBookHandler(db)).Methods(http.MethodPost)
	router.HandleFunc("/api/books/{id}", getBookHandler(db)).Methods(http.MethodGet)
	router.HandleFunc("/api/books/{id}", updateBookHandler(db)).Methods(http.MethodPost)
	router.HandleFunc("/api/books/{id}", deleteBookHandler(db)).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", router))
}
