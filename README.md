# iot-simulation-kafka
Simulating a meteorology producer-consumer model using Kafka.

A Go-based application that simulates IoT sensor data using Kafka. The application supports producing sensor data, consuming and storing data, and managing Kafka topics programmatically.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
  - [Producer](#producer)
  - [Consumer](#consumer)
  - [Create Topic](#create-topic)
  - [Delete Topic](#delete-topic)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [License](#license)

## Features

- **Producer**: Generates and sends simulated sensor data to a specified Kafka topic.
- **Consumer**: Listens to a Kafka topic, receives sensor data, and writes it to a CSV file.
- **Create Topic**: Programmatically creates a Kafka topic with specified configurations.
- **Delete Topic**: Programmatically deletes a specified Kafka topic.

## Prerequisites

- **Go**: Ensure Go is installed. [Download Go](https://golang.org/dl/)
- **Kafka Broker**: Access to a running Kafka broker. [Apache Kafka](https://kafka.apache.org/downloads)
- **Kafka-Go Library**: Managed via Go Modules.

## Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/yourusername/iot-simulation-kafka.git
   cd iot-simulation-kafka
   ```

2. **Initialize Go Modules**

   ```bash
   go mod init iot-simulation-kafka
   go get github.com/segmentio/kafka-go
   ```

3. **Build the Executable**

   ```bash
   go build -o weather.exe main.go producer.go consumer.go admin.go sensor.go
   ```

   This command compiles all necessary `.go` files into a single executable named `weather.exe`.

## Usage

The application supports four primary operations: `producer`, `consumer`, `create_topic`, and `delete_topic`. Each operation has its own set of flags.

### Producer

**Description**: Generates and sends simulated sensor data to a specified Kafka topic.

**Command**:

```bash
./weather.exe producer --broker=<BROKER_ADDRESS> --topic=<TOPIC_NAME>
```

**Flags**:

- `--broker` (Optional): Kafka broker address. Default is `localhost:9092`.
- `--topic` (Required): Kafka topic name.

### Consumer

**Description**: Listens to a Kafka topic, receives sensor data, and writes it to a CSV file. And also plot in realtime.

**Command**:

```bash
./weather.exe consumer --broker=<BROKER_ADDRESS> --topic=<TOPIC_NAME> --group=<GROUP_ID>
```

**Flags**:

- `--broker` (Optional): Kafka broker address. Default is `localhost:9092`.
- `--topic` (Required): Kafka topic name.
- `--group` (Optional): Kafka consumer group ID. Default is `weather_group`.

### Create Topic

**Description**: Programmatically creates a Kafka topic with specified configurations.

**Command**:

```bash
./weather.exe create_topic --broker=<BROKER_ADDRESS> --topic=<TOPIC_NAME> --partitions=<NUM_PARTITIONS> --replicas=<REPLICATION_FACTOR>
```

**Flags**:

- `--broker` (Optional): Kafka broker address. Default is `localhost:9092`.
- `--topic` (Required): Kafka topic name.
- `--partitions` (Optional): Number of partitions. Default is `1`.
- `--replicas` (Optional): Replication factor. Default is `1`.

### Delete Topic

**Description**: Programmatically deletes a specified Kafka topic.

**Command**:

```bash
./weather.exe delete_topic --broker=<BROKER_ADDRESS> --topic=<TOPIC_NAME>
```

**Flags**:

- `--broker` (Optional): Kafka broker address. Default is `localhost:9092`.
- `--topic` (Required): Kafka topic name.

## Examples

### 1. **Create a Kafka Topic**

```bash
./weather.exe create_topic --broker=164.92.76.15:9092 --topic=21881 --partitions=3 --replicas=2
```

**Output**:

```
Topic 21881 created successfully
```

### 2. **Run the Producer**

```bash
./weather.exe producer --broker=164.92.76.15:9092 --topic=21881
```

**Output**:

```
Kafka Producer started. Sending data every 15-30 seconds...
Message sent: {"temperature":54.23,"humidity":60,"wind_direction":"NE"}
Message sent: {"temperature":56.78,"humidity":55,"wind_direction":"SW"}
...
```

### 3. **Run the Consumer**

```bash
./weather.exe consumer --broker=164.92.76.15:9092 --topic=21881 --group=weather_group
```

**Output**:

```
Kafka Consumer started. Listening for messages and writing to sensor_data.csv...
Message received - Timestamp: 2024-04-27T12:34:56Z, Temperature: 54.23°C, Humidity: 60%, Wind Direction: NE
Message received - Timestamp: 2024-04-27T12:35:10Z, Temperature: 56.78°C, Humidity: 55%, Wind Direction: SW
...
```

**CSV Output (`sensor_data.csv`):**

```csv
Timestamp,Temperature (°C),Humidity (%),Wind Direction
2024-04-27T12:34:56Z,54.23,60,NE
2024-04-27T12:35:10Z,56.78,55,SW
...
```

### 4. **Delete a Kafka Topic**

```bash
./weather.exe delete_topic --broker=164.92.76.15:9092 --topic=21881
```

**Output**:

```
Topic 21881 deleted successfully
```

## Troubleshooting

- **Unknown Topic or Partition Error**: Ensure the Kafka topic exists. Use the `create_topic` subcommand to create it.

  ```bash
  ./weather.exe create_topic --broker=164.92.76.15:9092 --topic=21881
  ```

- **Connection Issues**: Verify that the Kafka broker address and port are correct and that the broker is reachable.

  ```bash
  telnet 164.92.76.15 9092
  ```

- **Permissions**: Ensure that the user running the application has the necessary permissions to create or delete Kafka topics.

- **CSV File Not Updating**: Check file permissions and ensure no other process is locking `sensor_data.csv`.

## License

This project is licensed under the [MIT License](LICENSE).
