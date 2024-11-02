// sensor.go

package main

// SensorData represents the data sent by the sensor
type SensorData struct {
	Temperature   float64 `json:"temperature"`
	Humidity      int     `json:"humidity"`
	WindDirection string  `json:"wind_direction"`
}
