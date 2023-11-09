package entity

import (
	"time"
)

type UserPromo struct {
	ID        int64      `gorm:"primaryKey;autoIncrement:true" json:"id"`
	PromoID   int64      `json:"promo_id"`
	UserID    int64      `json:"user_id"`
	IsActive  string     `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}
