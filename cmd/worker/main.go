package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/dedihartono801/promo-scheduler/internal/app/queue/kafka"
	"github.com/dedihartono801/promo-scheduler/pkg/config"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ready := make(chan bool)

	kafkaConsumer, err := kafka.NewKafkaConsumer(config.GetEnv("KAFKA_ADDRESS"), config.GetEnv("CONSUMER_GROUP"))
	if err != nil {
		log.Fatal(err)
	}

	err = kafkaConsumer.StartConsumerGroup(ready, config.GetEnv("PROMO_BIRTHDAY_TOPIC"))
	if err != nil {
		log.Fatal(err)
	}

	// Wait for the consumer group to be ready
	<-ready
	fmt.Println("Consumer group is ready")

	// Set up a signal handler to gracefully handle termination signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Wait for a termination signal
	<-signals
	fmt.Println("Terminating...")

	// Signal the consumer group to stop gracefully
	kafkaConsumer.CloseConsumerGroup()

	// Wait for the consumer group to finish before exiting
	<-ready
	fmt.Println("Consumer group terminated")
}
