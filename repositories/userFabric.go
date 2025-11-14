package repositories

import (
	"Gin_Event_App_Arc/models"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type UserFabricRepoI interface {
	LoadCerts(string) (*models.FabricUser, error)
	SaveRegisteredUser(*models.FabricUser) (*models.FabricUser, error)
	GetAllUsersRegistered() ([]models.FabricUser, error)
}

type UserFabricRepo struct {
	DB *gorm.DB
}

func InitUserFabricRepo(db *gorm.DB) *UserFabricRepo {
	return &UserFabricRepo{DB: db}
}

func (repo *UserFabricRepo) LoadCerts(username string) (*models.FabricUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var userFabric models.FabricUser
	err := repo.DB.WithContext(ctx).Where("enrollment_id =?", username).First(&userFabric)
	return &userFabric, err.Error
}

func (repo *UserFabricRepo) SaveRegisteredUser(user *models.FabricUser) (*models.FabricUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Create(user).Error
	return user, err
}

func (repo *UserFabricRepo) GetAllUsersRegistered() ([]models.FabricUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var userFabricList []models.FabricUser
	if err := repo.DB.WithContext(ctx).Find(&userFabricList).Error; err != nil {
		// If record not found, return empty slice (Find shouldn't return ErrRecordNotFound,
		// but we handle defensively)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.FabricUser{}, err
		}
		return nil, err
	}
	return userFabricList, nil
}
