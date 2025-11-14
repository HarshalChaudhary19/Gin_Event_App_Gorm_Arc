package repositories

import (
	"Gin_Event_App_Arc/models"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type GraphRepoCourseI interface {
	CreateCourse(string, string) (*models.CourseGraph, error)
	GetAllCourses() ([]*models.CourseGraph, error)
}

type GraphRepoCourse struct {
	DB *gorm.DB
}

func InitGraphRepoCourse(db *gorm.DB) *GraphRepo {
	return &GraphRepo{
		DB: db,
	}
}

func (repo *GraphRepo) CreateCourse(name string, description string) (*models.CourseGraph, error) {
	fmt.Println("HERERERERERER")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	catgraph := models.CourseGraph{
		// ID:          uuid.NewString(),
		Name:        name,
		Description: description,
	}
	fmt.Println("Course Graph send", catgraph)
	err := repo.DB.WithContext(ctx).Create(&catgraph).Error
	return &catgraph, err
}

func (repo *GraphRepo) GetAllCourses() ([]*models.CourseGraph, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var allCategories []*models.CourseGraph
	err := repo.DB.WithContext(ctx).Find(&allCategories).Error
	return allCategories, err
}
