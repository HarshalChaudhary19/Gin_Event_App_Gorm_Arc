package services

import (
	"Gin_Event_App_Arc/models"
	"Gin_Event_App_Arc/repositories"
	"fmt"
	"reflect"
)

type UserRoleService struct {
	UserRoleServ repositories.UserRoleRepoI
	RoleServ     repositories.RoleRepoI
	UserServ     repositories.UserRepoI
}

type UserRoleServiceI interface {
	Update(int, int) error
	GetRolesByUser(int) ([]*models.Role, error)
	GetUsersByRole(int) ([]*models.User, error)
	GetAll() ([]*models.UserRole, error)
}

func NewUserRoleService(userroleRepo repositories.UserRoleRepoI, roleRepo repositories.RoleRepoI, userRepo repositories.UserRepoI) UserRoleServiceI {
	return &UserRoleService{UserRoleServ: userroleRepo, RoleServ: roleRepo, UserServ: userRepo}
}

func (service *UserRoleService) GetRolesByUser(userId int) ([]*models.Role, error) {
	//Check if user Actually exists or not
	user, err1 := service.UserServ.Get(userId)
	if err1 != nil {
		if reflect.DeepEqual(user, models.User{}) {
			fmt.Println("Debug hora hai...")
			return nil, nil
		}
		return nil, err1
	}
	//Now get Roles by user
	rolelList, err := service.UserRoleServ.GetRolesByUser(userId)
	if err != nil {
		return nil, err
	}
	if rolelList == nil {
		fmt.Println("Ye to nil aarha hai")
	}
	return rolelList, nil

}

func (service *UserRoleService) GetUsersByRole(roleId int) ([]*models.User, error) {
	//Check if This Role exists or not
	role, err1 := service.RoleServ.GetRole(roleId)
	if err1 != nil {
		if reflect.DeepEqual(role, models.Role{}) {
			fmt.Println("This is where error is coming")
			return nil, nil
		}
		return nil, err1
	}
	rolelList, err := service.UserRoleServ.GetUsersByRole(roleId)
	if err != nil {
		fmt.Println("This is where error is coming 2")
		return nil, err
	}
	return rolelList, nil
}

func (service *UserRoleService) Update(userId, roleId int) error {

	err := service.UserRoleServ.Update(userId, roleId)
	if err != nil {
		return nil
	}
	return err
}

func (service *UserRoleService) GetAll() ([]*models.UserRole, error) {
	useroleList, err := service.UserRoleServ.GetAll()
	if err != nil {
		return nil, err
	}
	return useroleList, nil
}
