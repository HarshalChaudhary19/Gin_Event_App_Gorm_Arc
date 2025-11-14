package repositories

import (
	"Gin_Event_App_Arc/models"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type GraphRepoI interface {
	Create(string, string) (*models.CategoryGraph, error)
	GetAll() ([]*models.CategoryGraph, error)
}

type GraphRepo struct {
	DB *gorm.DB
}

func InitGraphRepo(db *gorm.DB) *GraphRepo {
	return &GraphRepo{
		DB: db,
	}
}

func (repo *GraphRepo) Create(name string, description string) (*models.CategoryGraph, error) {
	fmt.Println("Yahan tk to aarha hai")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	catgraph := models.CategoryGraph{
		// ID:          uuid.NewString(),
		Name:        name,
		Description: description,
	}
	err := repo.DB.WithContext(ctx).Create(&catgraph).Error
	return &catgraph, err
}

func (repo *GraphRepo) GetAll() ([]*models.CategoryGraph, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var allCategories []*models.CategoryGraph
	err := repo.DB.WithContext(ctx).Find(&allCategories).Error
	return allCategories, err
}
