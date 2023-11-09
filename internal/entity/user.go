package entity

import (
	"time"
)

type User struct {
	ID        int64      `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	DateBirth time.Time  `json:"date_birth"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}
