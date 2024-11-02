# iot-simulation-kafka

Simulating a meteorology producer-consumer model using Kafka with support for compact payloads.

A Go-based application that simulates IoT sensor data using Kafka. The application supports producing sensor data in both traditional JSON format and a compact 3-byte payload format, consuming and storing data, and managing Kafka topics programmatically.

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

- **Producer**: Generates and sends simulated sensor data to a specified Kafka topic in either JSON or compact format.
- **Consumer**: Listens to a Kafka topic, receives sensor data, writes it to a CSV file, and displays real-time plots for temperature and humidity.
- **Create Topic**: Programmatically creates a Kafka topic with specified configurations.
- **Delete Topic**: Programmatically deletes a specified Kafka topic.
- **Dual Mode Operation**: Supports both traditional JSON payloads and compact 3-byte payloads to accommodate environments with payload size constraints (e.g., LoRa, Sigfox).

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
   go get github.com/gizak/termui/v3
   go get github.com/segmentio/kafka-go
   ```

3. **Build the Executable**

   ```bash
   go build -o weather.exe main.go producer.go consumer.go admin.go sensor.go
   ```

   This command compiles all necessary `.go` files into a single executable named `weather.exe`.

## Usage

The application supports five primary operations: `producer`, `consumer`, `create_topic`, `delete_topic`, and `help`. Each operation has its own set of flags.

### Producer

**Description**: Generates and sends simulated sensor data to a specified Kafka topic in either JSON or compact mode.

**Command**:

```bash
./weather.exe producer --broker=<BROKER_ADDRESS> --topic=<TOPIC_NAME> --mode=<MODE> --min-interval=<MIN_SECONDS> --max-interval=<MAX_SECONDS>
```

**Flags**:

- `--broker` (Optional): Kafka broker address. Default is `localhost:9092`.
- `--topic` (Required): Kafka topic name.
- `--mode` (Optional): Mode of operation. Options:
  - `json` (Default): Sends data as JSON strings.
  - `compact`: Sends data as compact 3-byte payloads.
- `--min-interval` (Optional): Minimum interval between messages in seconds. Default is `15`.
- `--max-interval` (Optional): Maximum interval between messages in seconds. Default is `30`.

**Example**:

```bash
./weather.exe producer --broker=localhost:9092 --topic=weather --mode=compact --min-interval=10 --max-interval=20
```

### Consumer

**Description**: Listens to a Kafka topic, receives sensor data, writes it to a CSV file, and displays real-time plots for temperature and humidity. Supports both JSON and compact payload formats.

**Command**:

```bash
./weather.exe consumer --broker=<BROKER_ADDRESS> --topic=<TOPIC_NAME> --group=<GROUP_ID> --mode=<MODE>
```

**Flags**:

- `--broker` (Optional): Kafka broker address. Default is `localhost:9092`.
- `--topic` (Required): Kafka topic name.
- `--group` (Optional): Kafka consumer group ID. Default is `weather_group`.
- `--mode` (Optional): Mode of operation. Options:
  - `json` (Default): Expects data in JSON format.
  - `compact`: Expects data in compact 3-byte payloads.

**Example**:

```bash
./weather.exe consumer --broker=localhost:9092 --topic=weather --group=weather_group --mode=compact
```

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

**Example**:

```bash
./weather.exe create_topic --broker=localhost:9092 --topic=weather --partitions=3 --replicas=2
```

### Delete Topic

**Description**: Programmatically deletes a specified Kafka topic.

**Command**:

```bash
./weather.exe delete_topic --broker=<BROKER_ADDRESS> --topic=<TOPIC_NAME>
```

**Flags**:

- `--broker` (Optional): Kafka broker address. Default is `localhost:9092`.
- `--topic` (Required): Kafka topic name.

**Example**:

```bash
./weather.exe delete_topic --broker=localhost:9092 --topic=weather
```

## Examples

### 1. **Create a Kafka Topic**

```bash
./weather.exe create_topic --broker=localhost:9092 --topic=weather --partitions=3 --replicas=2
```

**Output**:

```
Topic weather created successfully.
```

### 2. **Run the Producer in JSON Mode**

```bash
./weather.exe producer --broker=localhost:9092 --topic=weather --mode=json --min-interval=10 --max-interval=20
```

**Output**:

```
Kafka Producer started in json mode. Sending data every 10-20 seconds...
Message sent: {"temperature":54.23,"humidity":60,"wind_direction":"NE"}
Message sent: {"temperature":56.78,"humidity":55,"wind_direction":"SW"}
...
```

### 3. **Run the Producer in Compact Mode**

```bash
./weather.exe producer --broker=localhost:9092 --topic=weather --mode=compact --min-interval=10 --max-interval=20
```

**Output**:

```
Kafka Producer started in compact mode. Sending data every 10-20 seconds...
Message sent: [0 125 56]
Message sent: [0 130 12]
...
```

### 4. **Run the Consumer in JSON Mode**

```bash
./weather.exe consumer --broker=localhost:9092 --topic=weather --group=weather_group --mode=json
```

**Output**:

```
[TermUI Interface Displaying Temperature and Humidity Charts]
Log:
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

### 5. **Run the Consumer in Compact Mode**

```bash
./weather.exe consumer --broker=localhost:9092 --topic=weather --group=weather_group --mode=compact
```

**Output**:

```
[TermUI Interface Displaying Temperature and Humidity Charts]
Log:
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

### 6. **Delete a Kafka Topic**

```bash
./weather.exe delete_topic --broker=localhost:9092 --topic=weather
```

**Output**:

```
Topic weather deleted successfully.
```

## Troubleshooting

- **Unknown Topic or Partition Error**: Ensure the Kafka topic exists. Use the `create_topic` subcommand to create it.

  ```bash
  ./weather.exe create_topic --broker=localhost:9092 --topic=weather
  ```

- **Connection Issues**: Verify that the Kafka broker address and port are correct and that the broker is reachable.

  ```bash
  telnet localhost 9092
  ```

- **Invalid Mode Error**: Ensure that the `--mode` flag is set to either `json` or `compact`.

  ```bash
  ./weather.exe producer --mode=compact
  ```

- **Permissions**: Ensure that the user running the application has the necessary permissions to create or delete Kafka topics.

- **CSV File Not Updating**: Check file permissions and ensure no other process is locking `sensor_data.csv`.

- **TermUI Initialization Failure**: Ensure that your terminal supports TermUI and that no other processes are interfering with the UI rendering.

## License

This project is licensed under the [MIT License](LICENSE).
