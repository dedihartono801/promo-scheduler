package repository

import (
	"github.com/dedihartono801/promo-scheduler/internal/entity"
	"gorm.io/gorm"
)

type PromoRepository interface {
	CreatePromo(tx *gorm.DB, promo *entity.Promo) (int64, error)
	CreateUserPromo(tx *gorm.DB, ti []*entity.UserPromo) error
	GetPromoByCode(code string) (*entity.Promo, error)
}

type promoRepository struct {
	database *gorm.DB
}

func NewPromoRepository(database *gorm.DB) PromoRepository {
	return &promoRepository{database}
}

func (r *promoRepository) CreatePromo(tx *gorm.DB, promo *entity.Promo) (int64, error) {
	result := tx.Table("promo").Create(promo)
	if result.Error != nil {
		return 0, result.Error
	}

	// Fetch the ID of the inserted record from the database
	insertedID := promo.ID // Assuming you have an ID field in your user struct

	return insertedID, nil
}

func (r *promoRepository) CreateUserPromo(tx *gorm.DB, up []*entity.UserPromo) error {
	return tx.Table("user_promo").Create(up).Error
}

func (r *promoRepository) GetPromoByCode(code string) (*entity.Promo, error) {
	var promo entity.Promo
	err := r.database.Table("promo").Where("code", code).First(&promo).Error
	return &promo, err
}
