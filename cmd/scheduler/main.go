package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dedihartono801/promo-scheduler/database"
	"github.com/dedihartono801/promo-scheduler/internal/app/queue/kafka"
	"github.com/dedihartono801/promo-scheduler/internal/app/repository"
	"github.com/dedihartono801/promo-scheduler/internal/app/usecase/scheduler"
	"github.com/dedihartono801/promo-scheduler/internal/delivery"
	"github.com/dedihartono801/promo-scheduler/pkg/config"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var db *gorm.DB

	// Number of retry attempts
	maxRetries := 5
	retryInterval := 5 * time.Second

	for retries := 1; retries <= maxRetries; retries++ {
		fmt.Printf("Attempt %d to connect to the database...\n", retries)

		// Attempt to connect to the database
		db, err = database.InitMysql()

		if err == nil {
			// Connection successful, break out of the loop
			fmt.Println("Connected to the database!")
			break
		}

		// Connection failed, wait for a short interval before retrying
		fmt.Printf("Error connecting to the database: %v\n", err)
		fmt.Printf("Retrying in %s...\n", retryInterval)
		time.Sleep(retryInterval)
	}

	if err != nil {
		// All retry attempts failed
		log.Fatalf("Failed to connect to the database after %d attempts. Error: %v\n", maxRetries, err)
	}

	kafkaProducer, err := kafka.NewKafkaProducer(config.GetEnv("KAFKA_ADDRESS"), config.GetEnv("PROMO_BIRTHDAY_TOPIC"))
	if err != nil {
		log.Fatalf(err.Error())
	}

	dbTransactionRepository := repository.NewDbTransactionRepository(db)
	promoRepository := repository.NewPromoRepository(db)
	userRepository := repository.NewUserRepository(db)
	schedulerService := scheduler.NewSchedulerService(promoRepository, userRepository, dbTransactionRepository, kafkaProducer)
	schedulerHandler := delivery.NewSchedulerHandler(schedulerService)

	// */30 * * * * * (every 30 second)

	// Membuat scheduler cron baru
	c := cron.New()

	//trigger every 1 minute
	err = c.AddFunc("*/60 * * * * *", func() {
		fmt.Println("Scheduler triggered")
		err := schedulerHandler.Scheduler()
		if err != nil {
			log.Fatal("Error scheduler: ", err)
		}
	})
	if err != nil {
		log.Fatal("Error adding cron job: ", err)
	}

	// Memulai scheduler cron
	c.Start()

	// Menjalankan program secara tak terbatas agar scheduler tetap berjalan
	select {}
}
