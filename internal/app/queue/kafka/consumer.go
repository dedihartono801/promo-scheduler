package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/dedihartono801/promo-scheduler/pkg/config"
	"github.com/dedihartono801/promo-scheduler/pkg/dto"
)

type Consumer interface {
	StartConsumerGroup(ready chan<- bool, topic string) error
	CloseConsumerGroup()
}

type consumer struct {
	consumer sarama.ConsumerGroup
}

func NewKafkaConsumer(kafkaBroker string, groupConsumer string) (Consumer, error) {
	// Set up configuration for the consumer group
	config := sarama.NewConfig()
	//config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	//config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Create a new consumer group
	consumerGroup, err := sarama.NewConsumerGroup([]string{kafkaBroker}, groupConsumer, config)
	if err != nil {
		return nil, err
	}

	return consumer{consumerGroup}, nil
}

// Helper function to start the consumer group and handle messages
func (csmr consumer) StartConsumerGroup(ready chan<- bool, topic string) error {
	for {
		log.Println("start worker kafka")
		err := csmr.consumer.Consume(context.Background(), []string{topic}, &MyConsumerGroupHandler{})
		if err != nil {
			return err
		}

		// Check if the consumer group is still running
		select {
		case <-csmr.consumer.Errors():
			err := <-csmr.consumer.Errors()
			// Handle the error here
			fmt.Println("Error occurred:", err)
			// You can add your custom error handling logic
		default:
			// The consumer group has stopped, signal that it's ready to restart
			ready <- false
		}
	}
}

func (csmr consumer) CloseConsumerGroup() {
	csmr.consumer.Close()
}

// Custom consumer group handler
type MyConsumerGroupHandler struct{}

func (h MyConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	// Called when the consumer group session is set up
	return nil
}

func (h MyConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	// Called when the consumer group session is closed
	return nil
}

func (h MyConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Process each message in a separate goroutine
	for message := range claim.Messages() {
		go func(msg *sarama.ConsumerMessage) {
			// Retry logic
			retryCount := 0
			maxRetries := 3
			retryInterval := 5 * time.Second
			var processingErr error
			err := processMessage(msg)
			if err == nil {
				log.Println("Message processed successfully")
				session.MarkMessage(msg, "")
				session.Commit()
			}
			if err != nil {
				for retryCount < maxRetries {

					log.Printf("Error processing message: %s", err)
					log.Printf("Retrying message processing in %s...", retryInterval)
					time.Sleep(retryInterval)
					err := processMessage(msg)
					if err == nil {
						log.Println("Message processed successfully after retry.")
						session.MarkMessage(msg, "")
						session.Commit()
						return
					}
					processingErr = err
					retryCount++
				}

				log.Printf("Max retries reached. Message processing failed: %s", processingErr)
				// Perform any necessary error handling, such as logging or sending to an error topic
			}
		}(message)
	}

	return nil
}

func processMessage(msg *sarama.ConsumerMessage) error {
	fmt.Printf("Received message: %s\n", string(msg.Value))
	var data dto.KafkaPromo
	err := json.Unmarshal([]byte(msg.Value), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	message := fmt.Sprintf(`
Selamat ulang tahun yang luar biasa, %s!
Hari ini adalah hari istimewa yang penuh dengan kebahagiaan dan cinta. Kami ingin merayakan momen indah ini bersama Anda.
Sebagai ungkapan terima kasih atas dukungan Anda, kami ingin memberikan penawaran eksklusif untuk ulang tahun Anda.

**Promo Ulang Tahun:**
🎁 **Diskon Rp 5000 untuk Semua Produk SayaKaya!**
Gunakan kode promo: BIRTHDAY_PROMO saat checkout.

Selamat ulang tahun sekali lagi!

Dengan cinta,
Tim SayaKaya
`, data.Name)

	url := "https://console.zenziva.net/wareguler/api/sendWA/"

	req := map[string]string{
		"userkey": config.GetEnv("ZENZIVA_API_KEY"),
		"passkey": config.GetEnv("ZENZIVA_SECRET_KEY"),
		"to":      data.Phone,
		"message": message,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return err
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Message sent successfully!")
	} else {
		fmt.Println("Failed to send message. Status code:", resp.StatusCode)
	}
	return nil
}
