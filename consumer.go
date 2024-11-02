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

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/segmentio/kafka-go"
)

// Consumer reads messages from a Kafka topic, writes them to a CSV file,
// and plots temperature and humidity in real-time.
func Consumer(brokerAddress, topic, groupID string) {
	// Initialize TermUI
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Create line charts for Temperature and Humidity
	tempChart := widgets.NewPlot()
	tempChart.Title = "Temperature (°C)"
	tempChart.Data = [][]float64{}
	tempChart.SetRect(0, 0, 120, 15)
	tempChart.AxesColor = ui.ColorWhite
	tempChart.LineColors[0] = ui.ColorRed
	tempChart.Marker = widgets.MarkerDot

	humChart := widgets.NewPlot()
	humChart.Title = "Humidity (%)"
	humChart.Data = [][]float64{}
	humChart.SetRect(0, 16, 120, 30)
	humChart.AxesColor = ui.ColorWhite
	humChart.LineColors[0] = ui.ColorBlue
	humChart.Marker = widgets.MarkerDot

	// Create a log box for displaying messages within TermUI
	logBox := widgets.NewParagraph()
	logBox.Title = "Log"
	logBox.SetRect(0, 32, 120, 40)
	logBox.Text = "Kafka Consumer started. Listening for messages...\nPress 'q' to quit."

	// Render the UI
	ui.Render(tempChart, humChart, logBox)

	// Kafka Reader (Consumer) configuration
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
		err = writer.Write([]string{"Timestamp", "Temperature (°C)", "Humidity (%)", "Wind Direction"})
		if err != nil {
			log.Fatalf("Error writing CSV headers: %v\n", err)
		}
		writer.Flush()
	}

	// Slices to hold data points for plotting
	var tempData []float64
	var humData []float64
	var timestamps []string

	// Slice to store log messages for inverse trim
	var logMessages []string
	const maxLogLines = 20 // Maximum number of lines to display in the log box

	// Context for handling cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to handle UI events
	uiEvents := ui.PollEvents()

	// Goroutine to handle UI events (like quitting)
	go func() {
		for e := range uiEvents {
			if e.Type == ui.KeyboardEvent {
				switch e.ID {
				case "q", "C-c":
					cancel()
					return
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			// Exit the loop if context is canceled
			return
		default:
			// Read message from Kafka
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					return
				}
				log.Printf("\nError reading message: %v\n", err)
				continue
			}

			// Parse the JSON message
			var data SensorData
			err = json.Unmarshal(m.Value, &data)
			if err != nil {
				log.Printf("\nError parsing JSON: %v\n", err)
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

			// Add the new log entry at the beginning of logMessages
			logEntry := fmt.Sprintf("Message received - Timestamp: %s, Temperature: %.2f°C, Humidity: %d%%, Wind Direction: %s",
				timestamp, data.Temperature, data.Humidity, data.WindDirection)
			logMessages = append([]string{logEntry}, logMessages...)

			// Limit the number of lines in logMessages to maxLogLines
			if len(logMessages) > maxLogLines {
				logMessages = logMessages[:maxLogLines]
			}

			// Set logBox text to the concatenated log messages
			logBox.Text = ""
			for _, message := range logMessages {
				logBox.Text += message + "\n"
			}

			// Append data to slices for plotting
			tempData = append(tempData, data.Temperature)
			humData = append(humData, float64(data.Humidity))
			timestamps = append(timestamps, timestamp)

			// Limit data points to the last 20 for clarity
			if len(tempData) > 20 {
				tempData = tempData[len(tempData)-20:]
				humData = humData[len(humData)-20:]
				timestamps = timestamps[len(timestamps)-20:]
			}

			// Only plot if there are at least two data points
			if len(tempData) >= 2 && len(humData) >= 2 {
				// Update the charts
				tempChart.Data = [][]float64{tempData}
				humChart.Data = [][]float64{humData}

				// Refresh the UI
				ui.Render(tempChart, humChart, logBox)
			}
		}
	}
}
