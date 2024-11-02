// producer.go

package main

import (
	"context"
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

// Generates sensor data
func generateSensorData() SensorData {
	return SensorData{
		Temperature:   generateTemperature(),
		Humidity:      generateHumidity(),
		WindDirection: generateWindDirection(),
	}
}

// Producer sends sensor data to the specified Kafka topic with configurable intervals and mode
func Producer(brokerAddress, topic string, minInterval, maxInterval int, mode string) {
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

	fmt.Printf("Kafka Producer started in %s mode. Sending data every %d-%d seconds...\n", mode, minInterval, maxInterval)

	for {
		sensorData := generateSensorData()

		var payload []byte
		var err error

		if mode == "json" {
			// JSON Mode
			payloadStr, err := sensorData.ToJSON()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating JSON sensor data: %v\n", err)
				continue
			}
			payload = []byte(payloadStr)
		} else if mode == "compact" {
			// Show original data
			fmt.Printf("Original data: %v\n", sensorData)

			// Compact Mode
			payload, err = sensorData.Encode()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error encoding sensor data: %v\n", err)
				continue
			}
		} else {
			fmt.Fprintf(os.Stderr, "Unknown mode: %s. Exiting.\n", mode)
			os.Exit(1)
		}

		// Create the Kafka message
		msg := kafka.Message{
			Key:   []byte("sensor1"), // You can customize the key as needed
			Value: payload,
		}

		// Send the message
		err = writer.WriteMessages(context.Background(), msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message to Kafka: %v\n", err)
		} else {
			if mode == "json" {
				fmt.Printf("Message sent: %s\n", string(payload))
			} else {
				fmt.Printf("Message sent: %v\n---\n", payload)
			}
		}

		// Validate and adjust intervals if necessary
		if minInterval > maxInterval {
			fmt.Printf("Warning: min-interval (%d) is greater than max-interval (%d). Adjusting max-interval to %d.\n", minInterval, maxInterval, minInterval)
			maxInterval = minInterval
		}

		// Wait between minInterval and maxInterval seconds
		waitTime := time.Duration(minInterval+rand.Intn(maxInterval-minInterval+1)) * time.Second
		time.Sleep(waitTime)
	}
}
