// consumer.go

package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func Consumer() {
	// Kafka Reader (Consumer) configuration
	brokerAddress := "164.92.76.15:9092" // IP and port of the Kafka server
	topic := "21881"
	groupID := "weather_group"

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerAddress},
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		StartOffset:    kafka.FirstOffset,
		CommitInterval: time.Second, // Interval for committing offsets
	})
	defer reader.Close()

	// Open or create the CSV file
	file, err := os.OpenFile("sensor_data.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error opening CSV file: %v\n", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers if the file is empty
	fi, err := file.Stat()
	if err != nil {
		log.Fatalf("Error getting file information: %v\n", err)
	}
	if fi.Size() == 0 {
		writer.Write([]string{"Timestamp", "Temperature (°C)", "Humidity (%)", "Wind Direction"})
	}

	fmt.Println("Kafka Consumer started. Listening for messages and writing to sensor_data.csv...")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			continue
		}

		// Parse the JSON message
		var data SensorData
		err = json.Unmarshal(m.Value, &data)
		if err != nil {
			log.Printf("Error parsing JSON: %v\n", err)
			continue
		}

		// Get the current timestamp
		timestamp := time.Now().Format(time.RFC3339)

		// Write the data to the CSV file
		record := []string{
			timestamp,
			fmt.Sprintf("%.2f", data.Temperature),
			fmt.Sprintf("%d", data.Humidity),
			data.WindDirection,
		}
		err = writer.Write(record)
		if err != nil {
			log.Printf("Error writing to CSV file: %v\n", err)
			continue
		}
		writer.Flush()

		// Print the received data
		fmt.Printf("Message received - Timestamp: %s, Temperature: %.2f°C, Humidity: %d%%, Wind Direction: %s\n",
			timestamp, data.Temperature, data.Humidity, data.WindDirection)
	}
}
