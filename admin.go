// admin.go

package main

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

// CreateTopic creates a Kafka topic with specified configurations
func CreateTopic(broker, topic string, numPartitions, replicationFactor int) error {
	// Establish a connection to the Kafka broker
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka broker: %v", err)
	}
	defer conn.Close()

	// Define the topic configuration
	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}

	// Attempt to create the topic
	err = conn.CreateTopics(topicConfig)
	if err != nil {
		return fmt.Errorf("failed to create topic: %v", err)
	}

	fmt.Printf("Topic %s created successfully\n", topic)
	return nil
}

// DeleteTopic deletes a Kafka topic
func DeleteTopic(broker, topic string) error {
	// Establish a connection to the Kafka broker
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka broker: %v", err)
	}
	defer conn.Close()

	// Attempt to delete the topic
	err = conn.DeleteTopics(topic)
	if err != nil {
		return fmt.Errorf("failed to delete topic: %v", err)
	}

	fmt.Printf("Topic %s deleted successfully\n", topic)
	return nil
}
