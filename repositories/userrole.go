package repositories

import (
	"Gin_Event_App_Arc/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type UserRoleRepo struct {
	DB *gorm.DB
}

type UserRoleRepoI interface {
	Update(int, int) error
	GetRolesByUser(int) ([]*models.Role, error)
	GetUsersByRole(int) ([]*models.User, error)
	Insert(*models.UserRole) error
	GetAll() ([]*models.UserRole, error)
}

func InitUserRoleRepo(db *gorm.DB) *UserRoleRepo {
	return &UserRoleRepo{
		DB: db,
	}
}

func (repo *UserRoleRepo) Update(userId, roleId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Where("user_id=? AND role_id=?", userId, roleId).Model(&models.UserRole{}).Updates(models.UserRole{ //Basic Update Query hai ye
		UserId: userId,
		RoleId: roleId,
	}).Error
	return err
}

func (repo *UserRoleRepo) GetRolesByUser(userId int) ([]*models.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var roleListByUser []*models.Role
	err := repo.DB.WithContext(ctx).Table("roles r").Select("r.id,r.name").Joins("JOIN user_roles ur ON ur.role_id = r.id").Where("ur.user_id=?", userId).Scan(&roleListByUser).Error
	return roleListByUser, err
}

func (repo *UserRoleRepo) GetUsersByRole(roleId int) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var userListByRole []*models.User
	err := repo.DB.WithContext(ctx).Table("users u").Select("u.id,u.name").Joins("JOIN user_roles ur ON ur.user_id = u.id").Where("ur.role_id=?", roleId).Scan(&userListByRole).Error
	return userListByRole, err
}

func (repo *UserRoleRepo) Insert(userrole *models.UserRole) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return repo.DB.WithContext(ctx).Create(userrole).Error
}

func (repo *UserRoleRepo) GetAll() ([]*models.UserRole, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var useroleList []*models.UserRole
	err := repo.DB.WithContext(ctx).Find(&useroleList).Error
	return useroleList, err
}
