package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type Producer interface {
	SendMessage([]*sarama.ProducerMessage) error
}

type producer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(broker string, topic string) (Producer, error) {
	// Set up configuration for the Kafka producer
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	// Create a new Kafka producer
	prd, err := sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		return nil, err
	}

	// Try to create the topic if it doesn't exist
	err = createTopic(prd, broker, topic)
	if err != nil {
		log.Printf("Error creating topic: %v", err)
	}

	return producer{prd}, nil

}

func (prd producer) SendMessage(message []*sarama.ProducerMessage) error {

	// Send the message to Kafka
	err := prd.producer.SendMessages(message)
	if err != nil {
		return err
	}

	log.Printf("Message sent to topic")
	return nil
}

func createTopic(prd sarama.SyncProducer, broker, topic string) error {
	client, err := sarama.NewClient([]string{broker}, nil)
	if err != nil {
		return err
	}
	defer client.Close()

	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		return err
	}
	defer admin.Close()

	// Check if the topic already exists
	topics, err := admin.ListTopics()
	if err != nil {
		return err
	}

	// If the topic doesn't exist, create it
	if _, exists := topics[topic]; !exists {
		err := admin.CreateTopic(topic, &sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		}, false)
		if err != nil {
			return err
		}

		log.Printf("Topic %s created successfully", topic)
	}

	return nil
}
