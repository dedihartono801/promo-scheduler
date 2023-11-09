package repository

import (
	"github.com/dedihartono801/promo-scheduler/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByDateBirth(datebirth string) ([]entity.User, error)
}

type userRepository struct {
	database *gorm.DB
}

func NewUserRepository(database *gorm.DB) UserRepository {
	return &userRepository{database}
}

func (r *userRepository) GetUserByDateBirth(datebirth string) ([]entity.User, error) {
	var user []entity.User
	err := r.database.Table("user").Where("DATE_FORMAT(date_birth, '%m-%d') = DATE_FORMAT(?, '%m-%d')", datebirth).Scan(&user).Error
	return user, err
}
