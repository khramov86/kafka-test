package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	requiredVars := []string{
		"KAFKA_BOOTSTRAP_SERVERS",
		"KAFKA_SECURITY_PROTOCOL",
		"KAFKA_SASL_MECHANISM",
		"KAFKA_SASL_USERNAME",
		"KAFKA_SASL_PASSWORD",
		"KAFKA_TOPIC",
		"KAFKA_NUM_PARTITIONS",
		"KAFKA_REPLICATION_FACTOR",
		"MESSAGE_COUNT",
	}

	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("Required environment variable %s is missing", v)
		}
	}

	// Kafka configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers":   os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"security.protocol":   os.Getenv("KAFKA_SECURITY_PROTOCOL"),
		"sasl.mechanism":      os.Getenv("KAFKA_SASL_MECHANISM"),
		"sasl.username":       os.Getenv("KAFKA_SASL_USERNAME"),
		"sasl.password":       os.Getenv("KAFKA_SASL_PASSWORD"),
		"api.version.request": true,
	}

	// Create Admin Client
	admin, err := kafka.NewAdminClient(config)
	if err != nil {
		log.Fatalf("Failed to create admin client: %v", err)
	}
	defer admin.Close()

	numPartitions, _ := strconv.Atoi(os.Getenv("KAFKA_NUM_PARTITIONS"))
	replicationFactor, _ := strconv.Atoi(os.Getenv("KAFKA_REPLICATION_FACTOR"))

	topic := kafka.TopicSpecification{
		Topic:             os.Getenv("KAFKA_TOPIC"),
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results, err := admin.CreateTopics(ctx, []kafka.TopicSpecification{topic})
	if err != nil {
		log.Fatalf("Failed to create topic: %v", err)
	}

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			log.Fatalf("Topic creation failed: %v", result.Error)
		}
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)
	messageCount, _ := strconv.Atoi(os.Getenv("MESSAGE_COUNT"))
	topicName := os.Getenv("KAFKA_TOPIC")
	for i := 0; i < messageCount; i++ {
		value := fmt.Sprintf("Message %d", i)
		msg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: kafka.PartitionAny},
			Value:          []byte(value),
		}

		err := producer.Produce(msg, deliveryChan)
		if err != nil {
			log.Printf("Failed to produce message %d: %v", i, err)
		}
		e := <-deliveryChan
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			log.Printf("Delivery failed: %v", m.TopicPartition.Error)
		} else {
			log.Printf("Delivered message to topic %s [%d] at offset %v",
				*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		}
	}
	producer.Flush(15 * 1000)
}
