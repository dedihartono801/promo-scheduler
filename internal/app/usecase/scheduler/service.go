package scheduler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/dedihartono801/promo-scheduler/internal/app/queue/kafka"
	"github.com/dedihartono801/promo-scheduler/internal/app/repository"
	"github.com/dedihartono801/promo-scheduler/internal/entity"
	"github.com/dedihartono801/promo-scheduler/pkg/config"
	"github.com/dedihartono801/promo-scheduler/pkg/dto"
	"gorm.io/gorm"
)

type Service interface {
	Scheduler() error
	FetchUserBirthDay(datebirth string) ([]entity.User, error)
	GeneratePromo(tx *gorm.DB, promo *entity.Promo) (int64, error)
	CreateUserPromo(tx *gorm.DB, userPromo []*entity.UserPromo) error
	RetrieveStartDateEndDate() (startDate, endDate time.Time)
}

type service struct {
	promoRepository repository.PromoRepository
	userRepository  repository.UserRepository
	dbTransaction   repository.DbTransactionRepository
	kafkaProducer   kafka.Producer
}

func NewSchedulerService(
	promoRepository repository.PromoRepository,
	userRepository repository.UserRepository,
	dbTransaction repository.DbTransactionRepository,
	kafkaProducer kafka.Producer,
) Service {
	return &service{
		promoRepository: promoRepository,
		userRepository:  userRepository,
		dbTransaction:   dbTransaction,
		kafkaProducer:   kafkaProducer,
	}
}

func (s *service) Scheduler() error {
	// Get the current date and time
	currentTime := time.Now().Format("2006-01-02")

	var userPromo []*entity.UserPromo
	var kafkaMessages []*sarama.ProducerMessage

	fmt.Println("get users birthdays on this day ", currentTime)
	users, err := s.FetchUserBirthDay(currentTime)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("none of the users have birthdays on this day")
		return nil
	}

	tx, err := s.dbTransaction.BeginTransaction()
	if err != nil {
		return err
	}

	_, err = s.promoRepository.GetPromoByCode("BIRTHDAY_PROMO")
	if err == nil {
		fmt.Println("Promo already created on this day")
		return nil
	}

	startDate, endDate := s.RetrieveStartDateEndDate()
	promoID, err := s.GeneratePromo(tx, &entity.Promo{
		PromoTypeID:  1,
		Name:         "Birthday Promo",
		Description:  "Birthday Promo",
		Code:         "BIRTHDAY_PROMO",
		DiscountType: "fix",
		Discount:     5000,
		UserType:     "valid_user",
		StartDate:    startDate,
		EndDate:      endDate,
		IsActive:     "1",
	})

	if err != nil {
		return err
	}

	for _, val := range users {
		jsonMessage, err := json.Marshal(dto.KafkaPromo{
			Name:     val.Name,
			Phone:    val.Phone,
			Birthday: val.DateBirth,
		})
		if err != nil {
			return err
		}
		userPromo = append(userPromo, &entity.UserPromo{
			PromoID:  promoID,
			UserID:   val.ID,
			IsActive: "1",
		})
		kafkaMessages = append(kafkaMessages, &sarama.ProducerMessage{
			Topic: config.GetEnv("PROMO_BIRTHDAY_TOPIC"),
			Value: sarama.StringEncoder(jsonMessage), // Use StringEncoder for JSON data
		})
	}

	err = s.CreateUserPromo(tx, userPromo)
	if err != nil {
		tx.Rollback()
		return err
	}

	fmt.Println("Promo successfully created")

	err = s.kafkaProducer.SendMessage(kafkaMessages)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = s.dbTransaction.CommitTransaction(tx)
	if err != nil {
		return err
	}

	fmt.Println("Event sent to topic")

	return nil
}

func (s *service) FetchUserBirthDay(datebirth string) ([]entity.User, error) {
	user, err := s.userRepository.GetUserByDateBirth(datebirth)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GeneratePromo(tx *gorm.DB, promo *entity.Promo) (int64, error) {
	return s.promoRepository.CreatePromo(tx, promo)
}

func (s *service) CreateUserPromo(tx *gorm.DB, userPromo []*entity.UserPromo) error {
	return s.promoRepository.CreateUserPromo(tx, userPromo)
}

func (s *service) RetrieveStartDateEndDate() (startDate, endDate time.Time) {
	// Get the current date and time
	currentTime := time.Now()

	// Get the start of the day (00:00:00)
	startOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())

	// Get the end of the day (23:59:59)
	endOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 999999999, currentTime.Location())

	return startOfDay, endOfDay
}
