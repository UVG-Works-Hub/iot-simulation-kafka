// main.go

package main

import (
	"fmt"
	"os"
)

func main() {
	// Read the arguments
	if len(os.Args) < 2 {
		fmt.Println("Error: No argument provided")
		fmt.Println("Usage: go run main.go [producer|consumer]")
		return
	}
	arg := os.Args[1]

	switch arg {
	case "producer":
		Producer()
	case "consumer":
		Consumer()
	default:
		fmt.Println("Error: Invalid argument. Use 'producer' or 'consumer'")
	}
}
