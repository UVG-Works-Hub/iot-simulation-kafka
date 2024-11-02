// producer.go

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

// Generates a temperature following a normal distribution centered at 55°C
func generateTemperature() float64 {
	mean := 55.0
	stddev := 10.0
	// Box-Muller transform
	u1 := rand.Float64()
	u2 := rand.Float64()
	z0 := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2*math.Pi*u2)
	temp := mean + z0*stddev
	// Ensure that the temperature is within the range [0, 110]
	if temp < 0 {
		temp = 0
	}
	if temp > 110 {
		temp = 110
	}
	// Round to two decimal places
	temp, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", temp), 64)
	return temp
}

// Generates relative humidity following a normal distribution centered at 55%
func generateHumidity() int {
	mean := 55.0
	stddev := 15.0
	// Box-Muller transform
	u1 := rand.Float64()
	u2 := rand.Float64()
	z0 := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2*math.Pi*u2)
	hum := mean + z0*stddev
	// Round and convert to integer
	hum = math.Round(hum)
	// Ensure that the humidity is within the range [0, 100]
	if hum < 0 {
		hum = 0
	}
	if hum > 100 {
		hum = 100
	}
	return int(hum)
}

// Generates a random wind direction
func generateWindDirection() string {
	directions := []string{"N", "NW", "W", "SW", "S", "SE", "E", "NE"}
	return directions[rand.Intn(len(directions))]
}

// Generates sensor data in JSON format
func generateSensorData() (string, error) {
	data := SensorData{
		Temperature:   generateTemperature(),
		Humidity:      generateHumidity(),
		WindDirection: generateWindDirection(),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// Producer sends sensor data to the specified Kafka topic
func Producer(brokerAddress, topic string) {
	// Initialize the seed for random number generation
	rand.Seed(time.Now().UnixNano())

	// Kafka Writer (Producer) configuration
	writer := kafka.Writer{
		Addr:         kafka.TCP(strings.Split(brokerAddress, ",")[0]),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false, // Synchronous writes to catch errors immediately
	}
	defer writer.Close()

	fmt.Println("Kafka Producer started. Sending data every 15-30 seconds...")

	for {
		sensorData, err := generateSensorData()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating sensor data: %v\n", err)
			continue
		}

		// Create the Kafka message
		msg := kafka.Message{
			Key:   []byte("sensor1"), // You can customize the key as needed
			Value: []byte(sensorData),
		}

		// Send the message
		err = writer.WriteMessages(context.Background(), msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message to Kafka: %v\n", err)
		} else {
			fmt.Printf("Message sent: %s\n", sensorData)
		}

		// Wait between 15 and 30 seconds
		waitTime := time.Duration(15+rand.Intn(16)) * time.Second
		time.Sleep(waitTime)
	}
}
