package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type ValueRequest struct {
	Weight float32 `json:"weight"`
	Diastolic int  `json:"diastolic"`
	Systolic  int  `json:"systolic"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	// Initialize database
	initDB()

	// Setup routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/post", handlePost)
	http.HandleFunc("/list", handleList) 

	// Start server
	fmt.Println("Health Tracker API server starting on :9001...")
	fmt.Println("Open https://tom-rose.de/healthtracker/ in your browser")
	log.Fatal(http.ListenAndServe(":9001", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Method not allowed. Use POST to echo a value.",
		})
		return
	}

	var req ValueRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Invalid JSON format",
		})
		return
	}

	db, err := sql.Open("sqlite3", "diary.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Database error"})
		return
	}
	defer db.Close()

	if req.Weight != 0 {
		_, err = db.Exec("INSERT INTO weights (timestamp, weight) VALUES (CURRENT_TIMESTAMP, ?)", req.Weight)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to save weight"})
			return
		}
	}

	if req.Diastolic != 0 {
		_, err = db.Exec("INSERT INTO blood (timestamp, diastolic, systolic) VALUES (CURRENT_TIMESTAMP, ?, ?)", req.Diastolic, req.Systolic)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to save diastolic"})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ValueRequest{Weight: req.Weight, Diastolic: req.Diastolic, Systolic: req.Systolic})
}

func handleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Method not allowed. Use GET to list values.",
		})
		return
	}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	db, err := sql.Open("sqlite3", "diary.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Database error"})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT timestamp, weight FROM weights ORDER BY timestamp DESC")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to query weights"})
		return
	}
	defer rows.Close()

	type WeightEntry struct {
		Timestamp string  `json:"timestamp"`
		Weight    float32 `json:"weight"`
	}

	var entries []WeightEntry
	for rows.Next() {
		var entry WeightEntry
		err := rows.Scan(&entry.Timestamp, &entry.Weight)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to read weight entry"})
			return
		}
		entries = append(entries, entry)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entries)
}

func initDB() {
	db, err := sql.Open("sqlite3", "diary.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS weights (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		weight REAL NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	fmt.Println("Database initialized successfully")
}
