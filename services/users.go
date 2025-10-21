package services

import (
	"Gin_Event_App_Arc/kafka"
	"Gin_Event_App_Arc/models"
	"Gin_Event_App_Arc/repositories"
	userpb "Gin_Event_App_Arc/src/proto"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type UserServiceI interface {
	Insertservice(*models.User) error
	Getservice(int) (*models.User, error)
	GetAllUsersservice(int) ([]*models.User, error)
	FindByEmailPassword(string) (*models.User, error)
}

type UserProto struct {
	userpb.UnimplementedUserServiceServer
	Serviceprot *UserService
}

type UserService struct { //Basically Dusre Package ke Interface ko point krna
	UserServ     repositories.UserRepoI
	UserRoleServ repositories.UserRoleRepoI
}

func NewUserService(userRepo repositories.UserRepoI, userolerepo repositories.UserRoleRepoI) UserServiceI {
	return &UserService{UserServ: userRepo, UserRoleServ: userolerepo}
}

func (s *UserProto) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	user, err := s.Serviceprot.UserServ.Get(int(req.Id))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Convert from your model.User to proto response
	return &userpb.GetUserResponse{
		Id:    int32(user.Id),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
func (serve *UserService) Insertservice(user *models.User) error {
	usernew, err2 := serve.UserServ.CreateUser(user)
	userolenew := &models.UserRole{
		UserId: usernew.Id,
		RoleId: 4,
	}
	if err2 != nil {
		return err2
	}
	err := serve.UserRoleServ.Insert(userolenew)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//Kafka User Sent start
	userSent := kafka.UserSent{
		Id:       usernew.Id,
		Email:    usernew.Email,
		Name:     usernew.Name,
		Password: usernew.Password,
	}
	key := strconv.Itoa(userSent.Id)
	if err := kafka.PublishEvent(ctx, key, userSent); err != nil {
		log.Printf("Failed to Sent Data to Kafka", err.Error())
	}
	//Kafka User Sent Complete
	if err != nil {
		return err
	}
	return nil
}

func (serve *UserService) Getservice(id int) (*models.User, error) {

	user, err := serve.UserServ.Get(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, err
}

func (serve *UserService) GetAllUsersservice(page int) ([]*models.User, error) {
	pageSize := 5
	offSet := (page - 1) * pageSize
	allUsers, err := serve.UserServ.GetAll(pageSize, offSet)
	if err != nil {
		return nil, nil
	}
	return allUsers, err
}

func (serve *UserService) FindByEmailPassword(userName string) (*models.User, error) {
	user, err := serve.UserServ.FindByEmail(userName)
	if err != nil {
		return nil, err
	}
	return user, nil
}
