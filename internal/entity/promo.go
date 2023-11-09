package entity

import (
	"time"
)

type Promo struct {
	ID                   int64      `gorm:"primaryKey;autoIncrement:true" json:"id"`
	PromoTypeID          int32      `json:"promo_type_id"`
	Name                 string     `json:"name"`
	Description          string     `json:"description"`
	Code                 string     `json:"code"`
	DiscountType         string     `json:"discount_type"`
	Discount             int        `json:"discount"`
	MaximumDiscountUsage int        `json:"maximum_discount_usage"`
	MinimumTransaction   int        `json:"minimum_transaction"`
	Image                string     `json:"image"`
	UserType             string     `json:"user_type"`
	StartDate            time.Time  `json:"start_date"`
	EndDate              time.Time  `json:"end_date"`
	IsActive             string     `json:"is_active"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at" gorm:"index"`
}
