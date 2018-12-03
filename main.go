package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

// Book a book struct
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author an author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var books []Book

func getBooksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
	}
}

func createBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var book Book
		_ = json.NewDecoder(r.Body).Decode(&book)
		book.ID = strconv.Itoa(rand.Intn(1000000))
		books = append(books, book)
		json.NewEncoder(w).Encode(book)
	}
}

func getBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		for _, item := range books {
			if item.ID == params["id"] {
				json.NewEncoder(w).Encode(item)
				return
			}
		}

		json.NewEncoder(w).Encode(&Book{})
	}
}

func updateBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var book Book
		params := mux.Vars(r)
		for index, item := range books {
			if item.ID == params["id"] {
				books = append(books[:index], books[index+1:]...)
				_ = json.NewDecoder(r.Body).Decode(&book)
				book.ID = item.ID
				books = append(books, book)
				json.NewEncoder(w).Encode(book)
				break
			}
		}
	}
}

func deleteBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		for index, item := range books {
			if item.ID == params["id"] {
				books = append(books[:index], books[index+1:]...)
				break
			}
		}
		json.NewEncoder(w).Encode(books)
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

	router.HandleFunc("/api/books", getBooksHandler(db)).Methods("GET")
	router.HandleFunc("/api/books", createBookHandler(db)).Methods("POST")
	router.HandleFunc("/api/books/{id}", getBookHandler(db)).Methods("GET")
	router.HandleFunc("/api/books/{id}", updateBookHandler(db)).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBookHandler(db)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
