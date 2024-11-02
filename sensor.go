// sensor.go

package main

import (
	"encoding/json"
	"errors"
	"math"
)

// SensorData represents the data sent by the sensor
type SensorData struct {
	Temperature   float64 `json:"temperature"`
	Humidity      int     `json:"humidity"`
	WindDirection string  `json:"wind_direction"`
}

// Encode converts SensorData into a 3-byte payload.
// Format (24 bits):
// - Temperature: 14 bits (0-16383) representing 0.0-110.0°C
// - Humidity: 7 bits (0-127) representing 0-100%
// - Wind Direction: 3 bits (0-7) representing 8 possible directions
func (s *SensorData) Encode() ([]byte, error) {
	// Validate humidity
	if s.Humidity < 0 || s.Humidity > 100 {
		return nil, errors.New("humidity out of range (0-100)")
	}

	// Map wind direction to 3 bits
	dirMap := map[string]uint8{
		"N":  0,
		"NW": 1,
		"W":  2,
		"SW": 3,
		"S":  4,
		"SE": 5,
		"E":  6,
		"NE": 7,
	}

	dir, exists := dirMap[s.WindDirection]
	if !exists {
		return nil, errors.New("invalid wind direction")
	}

	// Encode temperature: 0.0-110.0°C mapped to 0-16383
	tempScaled := int(math.Round((s.Temperature / 110.0) * 16383))
	if tempScaled < 0 {
		tempScaled = 0
	}
	if tempScaled > 16383 {
		tempScaled = 16383
	}

	// Assemble the 24 bits
	var payload uint32
	payload |= uint32(tempScaled&0x3FFF) << 10 // 14 bits for temperature
	payload |= uint32(s.Humidity&0x7F) << 3    // 7 bits for humidity
	payload |= uint32(dir & 0x07)              // 3 bits for wind direction

	// Convert to 3 bytes
	bytes := []byte{
		byte((payload >> 16) & 0xFF),
		byte((payload >> 8) & 0xFF),
		byte(payload & 0xFF),
	}

	return bytes, nil
}

// Decode decodes a 3-byte payload into SensorData and returns raw values for validation.
func Decode(data []byte) (*SensorData, int, int, int, error) {
	if len(data) != 3 {
		return nil, 0, 0, 0, errors.New("invalid payload size")
	}

	// Initialize raw values
	comingTemp := 0
	comingHumidity := 0
	comingWindDirection := 0

	// Assemble the 24 bits
	payload := uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2])

	// Extract temperature
	tempScaled := (payload >> 10) & 0x3FFF // 14 bits
	comingTemp = int(tempScaled)
	temperature := (float64(tempScaled) / 16383.0) * 110.0

	// Extract humidity
	humidity := int((payload >> 3) & 0x7F) // 7 bits
	comingHumidity = humidity

	// Extract wind direction
	dirCode := (payload) & 0x07 // 3 bits
	comingWindDirection = int(dirCode)

	// Map back wind direction
	dirMap := []string{"N", "NW", "W", "SW", "S", "SE", "E", "NE"}
	if dirCode >= uint32(len(dirMap)) {
		return nil, comingTemp, comingHumidity, comingWindDirection, errors.New("invalid wind direction code")
	}
	windDirection := dirMap[dirCode]

	return &SensorData{
		Temperature:   temperature,
		Humidity:      humidity,
		WindDirection: windDirection,
	}, comingTemp, comingHumidity, comingWindDirection, nil
}

// ToJSON converts SensorData to JSON string
func (s *SensorData) ToJSON() (string, error) {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// FromJSON populates SensorData from JSON string
func (s *SensorData) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), s)
}
