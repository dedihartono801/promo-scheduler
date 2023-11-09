package dto

import "time"

type KafkaPromo struct {
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Birthday time.Time `json:"birthday"`
}
