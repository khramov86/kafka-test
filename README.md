# Kafka Producer Utility

This utility is a Go-based tool designed to produce messages to a Kafka topic. It can be used to create a Kafka topic and send a specified number of messages to it. The utility is statically linked using `musl` for compatibility in environments where dynamic linking is not feasible.

## Features

- **Topic Creation**: Automatically creates a Kafka topic with configurable partitions and replication factor.
- **Message Production**: Sends a specified number of messages to the created Kafka topic.
- **Environment Configuration**: Uses environment variables for Kafka connection and topic configuration.
- **Static Build**: Built with `musl` for static linking, making it portable across different Linux environments.


## Prerequisites

- **Go**: Ensure Go is installed on your system. You can download it from [here](https://golang.org/dl/).
- **musl**: Install `musl-tools` on Ubuntu:
  ```bash
  sudo apt-get update
  sudo apt-get install musl-tools
  ```


## Configuration
The utility is configured using environment variables. Create a .env file in the root directory with the following variables:

```
KAFKA_BOOTSTRAP_SERVERS=<kafka-broker-address>
KAFKA_SECURITY_PROTOCOL=<security-protocol> # e.g., SASL_SSL
KAFKA_SASL_MECHANISM=<sasl-mechanism> # e.g., PLAIN
KAFKA_SASL_USERNAME=<sasl-username>
KAFKA_SASL_PASSWORD=<sasl-password>
KAFKA_TOPIC=<topic-name>
KAFKA_NUM_PARTITIONS=<number-of-partitions>
KAFKA_REPLICATION_FACTOR=<replication-factor>
MESSAGE_COUNT=<number-of-messages-to-produce>
```

## Build and Run
1. Build the Utility
To build the utility with musl for static linking, run the following command:
```
export CC=/usr/bin/musl-gcc
go build --ldflags '-linkmode external -extldflags "-static"' -tags musl -o kafka-producer
```
This will generate a statically linked binary named kafka-producer.

2. Run the Utility
Run the utility using the following command:

```
./kafka-producer
```
The utility will:

Create the specified Kafka topic (if it doesn't already exist).

Produce the specified number of messages to the topic.

Log the delivery status of each message.