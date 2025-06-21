package main

import (
	"database/sql"
	"log"
	"math/rand"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Hardcoded DB connection string (update as needed)
	dsn := "exampleuser:examplepass@tcp(localhost:3306)/exampledb"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("MySQL connection error:", err)
	}
	defer db.Close()

	// Delete all existing records before seeding
	_, err = db.Exec("DELETE FROM sensor_data")
	if err != nil {
		log.Fatalf("Failed to delete existing sensor_data: %v", err)
	}
	log.Println("All existing records deleted from sensor_data.")

	rooms := []string{"diningroom", "kitchen", "livingroom", "bedroom 1", "bedroom 2", "bedroom 3"}
	year, month, thisday := time.Now().Year(), time.Now().Month(), time.Now().Day()

	for day := thisday; day >= 1; day-- {
		for hour := 0; hour < 24; hour++ {
			log.Printf("Seeding data: day %d, hour %02d...", day, hour)
			// Build batch insert
			query := "INSERT INTO sensor_data (room, temperature, humidity, timestamp) VALUES "
			vals := []interface{}{}
			placeholders := []string{}
			for _, room := range rooms {
				ts := time.Date(year, month, day, hour, rand.Intn(60), 0, 0, time.UTC)
				temp := 17 + rand.Float64()*10 // 17–27 °C
				hum := 35 + rand.Float64()*25  // 35–60 %
				placeholders = append(placeholders, "(?, ?, ?, ?)")
				vals = append(vals, room, temp, hum, ts.Format("2006-01-02 15:04:05"))
			}
			query += strings.Join(placeholders, ",")
			_, err := db.Exec(query, vals...)
			if err != nil {
				log.Printf("Batch insert error (day %d hour %d): %v", day, hour, err)
			}
		}
	}

	log.Println("✅ Seeding complete: June 2024, 6 rooms, 1-hour intervals.")
}
