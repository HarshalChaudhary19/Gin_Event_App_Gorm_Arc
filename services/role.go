package services

import (
	"Gin_Event_App_Arc/models"
	"Gin_Event_App_Arc/repositories"
)

type RoleService struct {
	RoleServ repositories.RoleRepoI
}

type RoleServiceI interface {
	Insert(*models.Role) error
	Update(*models.Role, int) (*models.Role, error)
	Delete(int) error
	Get(int) (*models.Role, error)
	GetAll() ([]*models.Role, error)
}

func NewRoleService(roleRepo repositories.RoleRepoI) RoleServiceI {
	return &RoleService{RoleServ: roleRepo}
}

func (service *RoleService) Insert(role *models.Role) error {
	err := service.RoleServ.CreateRole(role)
	if err != nil {
		return err
	}
	return nil
}

func (service *RoleService) Update(role *models.Role, id int) (*models.Role, error) {
	newRole, err := service.RoleServ.UpdateRole(role, id)
	if err != nil {
		return nil, err
	}
	return newRole, nil
}

func (service *RoleService) Delete(id int) error {
	err := service.RoleServ.DeleteRole(id)
	if err != nil {
		return err
	}
	return nil
}

func (service *RoleService) Get(id int) (*models.Role, error) {
	newRole, err := service.RoleServ.GetRole(id)
	if err != nil {
		return nil, err
	}
	return newRole, nil
}

func (service *RoleService) GetAll() ([]*models.Role, error) {
	allRoles, err := service.RoleServ.GetAllRoles()
	if err != nil {
		return nil, err
	}
	return allRoles, nil
}
