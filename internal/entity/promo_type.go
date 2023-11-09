package entity

import (
	"time"
)

type PromoType struct {
	ID        int64      `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}
