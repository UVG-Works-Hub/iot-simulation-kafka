// main.go

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// ModeType defines the type for mode
type ModeType string

const (
	JSONMode    ModeType = "json"
	CompactMode ModeType = "compact"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: No subcommand provided")
		printUsage()
		os.Exit(1)
	}

	subcommand := os.Args[1]

	switch subcommand {
	case "producer":
		producerCmd := flag.NewFlagSet("producer", flag.ExitOnError)
		broker := producerCmd.String("broker", "localhost:9092", "Kafka broker address")
		topic := producerCmd.String("topic", "", "Kafka topic name (required)")
		minInterval := producerCmd.Int("min-interval", 15, "Minimum interval between messages (in seconds)")
		maxInterval := producerCmd.Int("max-interval", 30, "Maximum interval between messages (in seconds)")
		mode := producerCmd.String("mode", "json", "Mode of operation: 'json' or 'compact'")

		producerCmd.Parse(os.Args[2:])

		if *topic == "" {
			fmt.Println("Error: --topic is required for producer")
			producerCmd.Usage()
			os.Exit(1)
		}

		if *mode != string(JSONMode) && *mode != string(CompactMode) {
			fmt.Println("Error: --mode must be either 'json' or 'compact'")
			producerCmd.Usage()
			os.Exit(1)
		}

		Producer(*broker, *topic, *minInterval, *maxInterval, *mode)

	case "consumer":
		consumerCmd := flag.NewFlagSet("consumer", flag.ExitOnError)
		broker := consumerCmd.String("broker", "localhost:9092", "Kafka broker address")
		topic := consumerCmd.String("topic", "", "Kafka topic name (required)")
		group := consumerCmd.String("group", "weather_group", "Kafka consumer group ID")
		mode := consumerCmd.String("mode", "json", "Mode of operation: 'json' or 'compact'")

		consumerCmd.Parse(os.Args[2:])

		if *topic == "" {
			fmt.Println("Error: --topic is required for consumer")
			consumerCmd.Usage()
			os.Exit(1)
		}

		if *mode != string(JSONMode) && *mode != string(CompactMode) {
			fmt.Println("Error: --mode must be either 'json' or 'compact'")
			consumerCmd.Usage()
			os.Exit(1)
		}

		Consumer(*broker, *topic, *group, *mode)

	case "create_topic":
		createCmd := flag.NewFlagSet("create_topic", flag.ExitOnError)
		broker := createCmd.String("broker", "localhost:9092", "Kafka broker address")
		topic := createCmd.String("topic", "", "Kafka topic name (required)")
		partitions := createCmd.Int("partitions", 1, "Number of partitions")
		replicas := createCmd.Int("replicas", 1, "Replication factor")

		createCmd.Parse(os.Args[2:])

		if *topic == "" {
			fmt.Println("Error: --topic is required for create_topic")
			createCmd.Usage()
			os.Exit(1)
		}

		err := CreateTopic(*broker, *topic, *partitions, *replicas)
		if err != nil {
			log.Fatalf("Failed to create topic: %v", err)
		}

	case "delete_topic":
		deleteCmd := flag.NewFlagSet("delete_topic", flag.ExitOnError)
		broker := deleteCmd.String("broker", "localhost:9092", "Kafka broker address")
		topic := deleteCmd.String("topic", "", "Kafka topic name (required)")

		deleteCmd.Parse(os.Args[2:])

		if *topic == "" {
			fmt.Println("Error: --topic is required for delete_topic")
			deleteCmd.Usage()
			os.Exit(1)
		}

		err := DeleteTopic(*broker, *topic)
		if err != nil {
			log.Fatalf("Failed to delete topic: %v", err)
		}

	default:
		fmt.Printf("Error: Unknown subcommand '%s'\n", subcommand)
		printUsage()
		os.Exit(1)
	}
}

// printUsage prints the usage information for the application
func printUsage() {
	fmt.Println(`Usage:
	go run main.go <subcommand> [options]

Subcommands:
	producer      Start the Kafka producer
	consumer      Start the Kafka consumer
	create_topic  Create a Kafka topic
	delete_topic  Delete a Kafka topic

Use "go run main.go <subcommand> --help" for more information about a subcommand.`)
}
