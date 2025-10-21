package repositories

import (
	"Gin_Event_App_Arc/models"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserRepoI interface {
	CreateUser(*models.User) (*models.User, error)
	Get(int) (*models.User, error)
	GetAll(int, int) ([]*models.User, error)
	FindByEmail(string) (*models.User, error)
}

type UserRepo struct {
	DB *gorm.DB
}

func InitUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (repo *UserRepo) CreateUser(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := repo.DB.WithContext(ctx).Create(user).Error
	return user, err
}

func (repo *UserRepo) Get(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user models.User
	//Here first means the first result that comes
	err := repo.DB.WithContext(ctx).First(&user, id)
	return &user, err.Error
}

func (repo *UserRepo) GetAll(pageLimit, offSet int) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var users []*models.User
	//Here first means the first result that comes
	err := repo.DB.WithContext(ctx).Limit(pageLimit).Offset(offSet).Find(&users)
	return users, err.Error
}

func (repo *UserRepo) FindByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user models.User
	err := repo.DB.WithContext(ctx).Where("email =?", email).First(&user)
	fmt.Println("Error ye rha", err.Error)
	return &user, err.Error
}

// func (repo *UserRepo) FindPassword(email string) (string,error){
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	pass,err:=repo.DB.WithContext(ctx).Where()
// }
