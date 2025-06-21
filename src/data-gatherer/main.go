package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
)

type SensorPayload struct {
	Room        string  `json:"room"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Timestamp   string  `json:"timestamp"`
}

func insertToDB(db *sql.DB, data SensorPayload) {
	parsedTime, err := time.Parse(time.RFC3339, data.Timestamp)
	if err != nil {
		log.Printf("Time parse error: %v", err)
		return
	}

	formatted := parsedTime.Format("2006-01-02 15:04:05")

	_, err = db.Exec(
		`INSERT INTO sensor_data (room, temperature, humidity, timestamp) VALUES (?, ?, ?, ?)`,
		data.Room, data.Temperature, data.Humidity, formatted,
	)
	if err != nil {
		log.Printf("DB insert error: %v", err)
	}
}

func main() {
	// Load environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	mqttHost := os.Getenv("MQTT_HOST")
	mqttPort := os.Getenv("MQTT_PORT")
	mqttUser := os.Getenv("MQTT_USER")
	mqttPass := os.Getenv("MQTT_PASS")

	// Connect to MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	// Connect to MQTT broker
	broker := fmt.Sprintf("tcp://%s:%s", mqttHost, mqttPort)
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID("go-sensor-consumer").
		SetUsername(mqttUser).
		SetPassword(mqttPass)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("MQTT connection error:", token.Error())
	}

	// Subscribe to incoming sensor data
	client.Subscribe("home/+", 1, func(client mqtt.Client, msg mqtt.Message) {
		var data SensorPayload
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			log.Printf("JSON unmarshal error: %v", err)
			return
		}
		insertToDB(db, data)
	})

	log.Println("Consumer is running and subscribed to home/+")
	select {} // Keep the app alive
}
