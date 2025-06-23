package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type DataPoint struct {
	Period      string  `json:"period"`
	MaxTemp     float64 `json:"max_temp"`
	MaxHumidity float64 `json:"max_humidity"`
}

func getRoomStats(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		room := vars["roomname"]
		period := r.URL.Query().Get("period")

		now := time.Now()
		var start, end time.Time
		var dateFormat string

		switch period {
		case "hour":
			start = now.Add(-1 * time.Hour)
			end = now
			dateFormat = "%Y-%m-%d %H:%i"
		case "day":
			start = now.AddDate(0, 0, -1)
			end = now
			dateFormat = "%Y-%m-%d %H"
		case "week":
			start = now.AddDate(0, 0, -7)
			end = now
			dateFormat = "%Y-%m-%d"
		case "month":
			start = now.AddDate(0, -1, 0)
			end = now
			dateFormat = "%Y-%m-%d"
		case "year":
			start = now.AddDate(-1, 0, 0)
			end = now
			dateFormat = "%Y-%m-%d"
		default:
			http.Error(w, "invalid period", http.StatusBadRequest)
			return
		}

		query := `
			SELECT 
				DATE_FORMAT(timestamp, ?) as period,
				MAX(temperature) as max_temp,
				MAX(humidity) as max_humidity
			FROM sensor_data
			WHERE room = ?
			  AND timestamp >= ?
			  AND timestamp < ?
			GROUP BY period
			ORDER BY period ASC
		`

		rows, err := db.Query(query, dateFormat, room, start, end)
		if err != nil {
			http.Error(w, "DB query error", http.StatusInternalServerError)
			log.Printf("DB query error: %v", err)
			return
		}
		defer rows.Close()

		var results []DataPoint
		for rows.Next() {
			var dp DataPoint
			if err := rows.Scan(&dp.Period, &dp.MaxTemp, &dp.MaxHumidity); err != nil {
				http.Error(w, "row scan error", http.StatusInternalServerError)
				return
			}
			results = append(results, dp)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

func listRooms(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT DISTINCT room FROM sensor_data ORDER BY room ASC`)
		if err != nil {
			http.Error(w, "DB query error", http.StatusInternalServerError)
			log.Printf("DB query error: %v", err)
			return
		}
		defer rows.Close()

		var rooms []string
		for rows.Next() {
			var room string
			if err := rows.Scan(&room); err != nil {
				http.Error(w, "row scan error", http.StatusInternalServerError)
				return
			}
			rooms = append(rooms, room)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rooms)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("MySQL connection error:", err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/room", listRooms(db)).Methods("GET")
	r.HandleFunc("/api/room/{roomname}", getRoomStats(db)).Methods("GET")

	log.Println("🌐 API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(r)))
}
