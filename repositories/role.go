package repositories

import (
	"Gin_Event_App_Arc/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type RoleRepo struct {
	DB *gorm.DB
}

type RoleRepoI interface {
	CreateRole(*models.Role) error
	GetRole(int) (*models.Role, error)
	UpdateRole(*models.Role, int) (*models.Role, error)
	DeleteRole(int) error
	GetAllRoles() ([]*models.Role, error)
}

func InitRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{
		DB: db,
	}
}

func (repo *RoleRepo) CreateRole(role *models.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return repo.DB.WithContext(ctx).Create(role).Error
}

func (repo *RoleRepo) GetRole(id int) (*models.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var role models.Role
	//Here first means the first result that comes
	err := repo.DB.WithContext(ctx).First(&role, id)
	return &role, err.Error
}

func (repo *RoleRepo) UpdateRole(role *models.Role, id int) (*models.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Where("id=?", role.Id).Model(&models.Role{}).Updates(models.Role{ //Basic Update Query hai ye
		Name: role.Name,
	}).Error
	return role, err
}

func (repo *RoleRepo) DeleteRole(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return repo.DB.WithContext(ctx).Delete(&models.Role{}, id).Error
}

func (repo *RoleRepo) GetAllRoles() ([]*models.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var allRoles []*models.Role
	err := repo.DB.WithContext(ctx).Find(&allRoles).Error
	return allRoles, err
}
