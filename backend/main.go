package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
	School   string `json:"school"`
}

var db *sql.DB

func getItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // CHo reacr tu domain khac goi API nay

	// STH to test
	rows, err := db.Query("select id, name, birthday, school from users order by name asc")
	if err != nil {
		log.Printf("Er when querying database: %v", err)
		http.Error(w, "Server err", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Birthday,
			&user.School,
		); err != nil {
			log.Printf("Err when scanning : %v", err)
			continue
		}
		users = append(users, user)
	}
	log.Printf("Retrieved %d users from database", len(users))

	json.NewEncoder(w).Encode(users)
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Cho phép các domain khác gọi API này

	response := map[string]string{"message": "Hello, World!"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func main() {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("DB_HOST"),
	)

	log.Printf("Connecting to database with connection string: %s", connStr)

	var err error // Connect to the database
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// check connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// Create router
	r := mux.NewRouter()
	r.HandleFunc("/api/users", getItem).Methods("GET")
	r.HandleFunc("/test", test).Methods("GET")

	log.Printf("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
